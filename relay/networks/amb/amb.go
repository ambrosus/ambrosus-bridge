package amb

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/metric"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

const BridgeName = "ambrosus"

type Bridge struct {
	Client         *ethclient.Client
	Contract       *contracts.Amb
	ContractRaw    *contracts.AmbRaw
	VSContract     *contracts.Vs
	HttpUrl        string // TODO: delete this field
	sideBridge     networks.BridgeReceiveAura
	auth           *bind.TransactOpts
	ExternalLogger external_logger.ExternalLogger
}

func (b *Bridge) SubmitEpochData(
	epoch *big.Int,
	fullSizeIn128Resultion *big.Int,
	branchDepth *big.Int,
	nodes []*big.Int,
	start *big.Int,
	merkelNodesNumber *big.Int,
) error {
	tx, txErr := b.Contract.SetEpochData(b.auth, epoch, fullSizeIn128Resultion, branchDepth, nodes, start, merkelNodesNumber)
	if txErr != nil {
		return b.getFailureReasonViaCall(
			txErr,
			"setEpochData",
			epoch,
			fullSizeIn128Resultion,
			branchDepth,
			nodes,
			start,
			merkelNodesNumber,
		)
	}

	// Metric
	if err := metric.SetContractBalance(BridgeName, b.Client, b.auth.From); err != nil {
		return err
	}

	return b.waitForTxMined(tx)
}

// New creates a new ambrosus bridge.
func New(cfg *config.AMBConfig, externalLogger external_logger.ExternalLogger) (*Bridge, error) {
	// Creating a new ethereum client.
	client, err := ethclient.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	// Creating a new ambrosus bridge contract instance.
	contract, err := contracts.NewAmb(common.HexToAddress(cfg.ContractAddr), client)
	if err != nil {
		return nil, err
	}

	// Creating a new ambrosus VS contract instance.
	vsContract, err := contracts.NewVs(common.HexToAddress(cfg.VSContractAddr), client)
	if err != nil {
		return nil, err
	}

	var auth *bind.TransactOpts

	if cfg.PrivateKey != nil {
		chainId, err := client.ChainID(context.Background())
		if err != nil {
			return nil, err
		}

		auth, err = bind.NewKeyedTransactorWithChainID(cfg.PrivateKey, chainId)
		if err != nil {
			return nil, err
		}

		// Metric
		if err := metric.SetContractBalance(BridgeName, client, auth.From); err != nil {
			return nil, fmt.Errorf("failed to init metric for ambrosus bridge: %w", err)
		}
	}

	return &Bridge{
		Client:         client,
		Contract:       contract,
		ContractRaw:    &contracts.AmbRaw{Contract: contract},
		VSContract:     vsContract,
		HttpUrl:        "https://network.ambrosus.io",
		auth:           auth,
		ExternalLogger: externalLogger,
	}, nil
}

func (b *Bridge) SubmitTransferPoW(proof *contracts.CheckPoWPoWProof) error {
	tx, txErr := b.Contract.SubmitTransfer(b.auth, *proof)

	if txErr != nil {
		// we've got here probably due to error at eth_estimateGas (e.g. revert(), require())
		// openethereum doesn't give us a full error message
		// so, make low-level call method to get the full error message
		return b.getFailureReasonViaCall(txErr, "submitTransfer", *proof)
	}

	// Metric
	if err := metric.SetContractBalance(BridgeName, b.Client, b.auth.From); err != nil {
		return err
	}

	return b.waitForTxMined(tx)
}

// GetLastEventId gets last contract event id.
func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

// GetEventById gets contract event by id.
func (b *Bridge) GetEventById(eventId *big.Int) (*contracts.TransferEvent, error) {
	opts := &bind.FilterOpts{Context: context.Background()}

	logs, err := b.Contract.FilterTransfer(opts, []*big.Int{eventId})
	if err != nil {
		return nil, err
	}

	if logs.Next() {
		return &logs.Event.TransferEvent, nil
	}

	return nil, networks.ErrEventNotFound
}

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.BridgeReceiveAura) {
	b.sideBridge = sideBridge

	for {
		if err := b.listen(); err != nil {
			log.Error().Err(err).Msg("listen ambrosus error")

			err = b.ExternalLogger.LogError(err)
			if err != nil {
				log.Error().Err(err).Msg("external logger error")
			}
		}
	}
}

func (b *Bridge) checkOldEvents() error {
	lastEventId, err := b.sideBridge.GetLastEventId()
	if err != nil {
		return err
	}

	i := big.NewInt(1)
	for {
		nextEventId := big.NewInt(0).Add(lastEventId, i)
		nextEvent, err := b.GetEventById(nextEventId)
		if err != nil {
			if errors.Is(err, networks.ErrEventNotFound) {
				// no more old events
				return nil
			}
			return err
		}

		err = b.sendEvent(nextEvent)
		if err != nil {
			return err
		}

		i = big.NewInt(0).Add(i, big.NewInt(1))
	}
}

func (b *Bridge) listen() error {
	err := b.checkOldEvents()
	if err != nil {
		return err
	}

	// Subscribe to events
	watchOpts := &bind.WatchOpts{Context: context.Background()}
	eventChannel := make(chan *contracts.AmbTransfer) // <-- тут я хз как сделать общий(common) тип для канала
	eventSub, err := b.Contract.WatchTransfer(watchOpts, eventChannel, nil)
	if err != nil {
		return err
	}

	defer eventSub.Unsubscribe()

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			return err
		case event := <-eventChannel:
			if err := b.sendEvent(&event.TransferEvent); err != nil {
				return err
			}
		}
	}
}

func (b *Bridge) sendEvent(event *contracts.TransferEvent) error {
	// Wait for safety blocks.
	safetyBlocks, err := b.sideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return err
	}

	if err := ethereum.WaitForBlock(b.Client, event.Raw.BlockNumber+safetyBlocks); err != nil {
		return err
	}

	// Check if the event has been removed.
	if err := b.isEventRemoved(event); err != nil {
		return err
	}

	ambTransfer, err := b.getBlocksAndEvents(event)
	if err != nil {
		return err
	}

	return b.sideBridge.SubmitTransferAura(ambTransfer)
}

func (b *Bridge) isEventRemoved(event *contracts.TransferEvent) error {
	block, err := b.HeaderByNumber(big.NewInt(int64(event.Raw.BlockNumber)))
	if err != nil {
		return err
	}

	if block.Hash(true) != event.Raw.BlockHash {
		return fmt.Errorf("block hash != event's block hash")
	}
	return nil
}

func (b *Bridge) GetMinSafetyBlocksNum() (uint64, error) {
	safetyBlocks, err := b.Contract.MinSafetyBlocks(nil)
	if err != nil {
		return 0, err
	}
	return safetyBlocks.Uint64(), nil
}
