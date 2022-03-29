package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/metric"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

const BridgeName = "ethereum"

type Bridge struct {
	Client     *ethclient.Client
	WsClient   *ethclient.Client
	Contract   *contracts.Eth
	WsContract *contracts.Eth
	sideBridge networks.BridgeReceiveEthash
	config     *config.ETHConfig
	auth       *bind.TransactOpts
	cfg        *config.ETHConfig
	logger     external_logger.ExternalLogger
}

// New creates a new ethereum bridge.
func New(cfg *config.ETHConfig, externalLogger external_logger.ExternalLogger) (*Bridge, error) {
	log.Debug().Str("bridge", BridgeName).Msg("Creating ethereum bridge...")

	// Creating a new ethereum client (HTTP & WS).
	client, err := ethclient.Dial(cfg.HttpURL)
	if err != nil {
		return nil, err
	}
	// Compatibility with tests.
	var wsClient *ethclient.Client
	if cfg.WsURL != "" {
		wsClient, err = ethclient.Dial(cfg.WsURL)
		if err != nil {
			return nil, err
		}
	}

	// Creating a new ambrosus bridge contract instance (HTTP & WS).
	contract, err := contracts.NewEth(common.HexToAddress(cfg.ContractAddr), client)
	if err != nil {
		return nil, err
	}
	// Compatibility with tests.
	var wsContract *contracts.Eth
	if wsClient != nil {
		wsContract, err = contracts.NewEth(common.HexToAddress(cfg.ContractAddr), wsClient)
		if err != nil {
			return nil, err
		}
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
		metric.SetContractBalance(BridgeName, client, auth.From)
	}

	return &Bridge{
		Client:     client,
		WsClient:   wsClient,
		Contract:   contract,
		WsContract: wsContract,
		auth:       auth,
		cfg:        cfg,
		logger:     externalLogger,
	}, nil
}

func (b *Bridge) SubmitTransferAura(proof *contracts.CheckAuraAuraProof) error {
	// Metric
	defer metric.SetContractBalance(BridgeName, b.Client, b.auth.From)

	tx, err := b.Contract.SubmitTransfer(b.auth, *proof)
	if err != nil {
		return err
	}

	receipt, err := bind.WaitMined(context.Background(), b.Client, tx)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		// we've got here probably due to low gas limit,
		// and revert() that hasn't been caught at eth_estimateGas
		return ethereum.GetFailureReason(b.Client, b.auth, tx)
	}

	return nil
}

func (b *Bridge) GetValidatorSet() ([]common.Address, error) {
	return b.Contract.GetValidatorSet(nil)
}

func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

func (b *Bridge) GetLastProcessedBlockNum() (*big.Int, error) {
	blockHash, err := b.Contract.LastProcessedBlock(nil)
	if err != nil {
		return nil, err
	}

	header, err := b.Client.HeaderByHash(context.Background(), blockHash)
	return header.Number, err
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

func (b *Bridge) Run(sideBridge networks.BridgeReceiveEthash) {
	log.Debug().Str("bridge", BridgeName).Msg("Running ethereum bridge...")

	b.sideBridge = sideBridge

	// Getting last ethereum block number.
	blockNumber, err := b.Client.BlockNumber(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("error getting last block number")
	}

	// Checking epoch data dir.
	if err = b.checkEpochDataDir(blockNumber/30000, b.cfg.EpochLength); err != nil {
		log.Error().Err(err).Msg("error checking epoch data dir")
	}

	log.Info().Str("bridge", BridgeName).Msg("Ethereum bridge runned!")

	for {
		if err := b.listen(); err != nil {
			log.Error().Err(err).Msg("listen ethereum error")

			if err := b.logger.LogError(err); err != nil {
				log.Error().Err(err).Msg("external logger error")
			}
		}
	}
}

func (b *Bridge) checkOldEvents() error {
	log.Info().Str("bridge", BridgeName).Msg("Checking old events...")

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

		log.Info().
			Str("bridge", BridgeName).
			Str("event_id", nextEventId.String()).
			Msg("Send old event...")

		if err := b.sendEvent(nextEvent); err != nil {
			return err
		}

		i = big.NewInt(0).Add(i, big.NewInt(1))
	}
}

func (b *Bridge) listen() error {
	log.Debug().Str("bridge", BridgeName).Msg("Listening ethereum events...")

	err := b.checkOldEvents()
	if err != nil {
		return err
	}

	// Subscribe to events
	watchOpts := &bind.WatchOpts{Context: context.Background()}
	eventChannel := make(chan *contracts.EthTransfer) // <-- тут я хз как сделать общий(common) тип для канала
	eventSub, err := b.WsContract.WatchTransfer(watchOpts, eventChannel, nil)
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
			log.Info().
				Str("bridge", BridgeName).
				Str("event_id", event.EventId.String()).
				Msg("Send event...")

			if err := b.sendEvent(&event.TransferEvent); err != nil {
				return err
			}
		}
	}
}

func (b *Bridge) sendEvent(event *contracts.TransferEvent) error {
	log.Debug().
		Str("bridge", BridgeName).
		Str("event_id", event.EventId.String()).
		Msg("Waiting for safety blocks...")

	// Wait for safety blocks.
	safetyBlocks, err := b.sideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return err
	}

	if err := ethereum.WaitForBlock(b.WsClient, event.Raw.BlockNumber+safetyBlocks); err != nil {
		return err
	}

	log.Debug().
		Str("bridge", BridgeName).
		Str("event_id", event.EventId.String()).
		Msg("Checking if the event has been removed...")

	// Check if the event has been removed.
	if err := b.isEventRemoved(event); err != nil {
		return err
	}

	ambTransfer, err := b.getBlocksAndEvents(event)
	if err != nil {
		return err
	}

	log.Debug().
		Str("bridge", BridgeName).
		Str("event_id", event.EventId.String()).
		Msg("Submit transfer PoW...")

	if err := b.sideBridge.SubmitTransferPoW(ambTransfer); err != nil {
		if errors.Is(err, networks.ErrEpochData) {
			epoch := event.Raw.BlockNumber / 30000

			epochData, err := b.loadEpochDataFile(epoch)
			if err != nil {
				return err
			}

			if err := b.SetEpochData(epochData); err != nil {
				return err
			}

			if err := b.sideBridge.SubmitTransferPoW(ambTransfer); err != nil {
				return err
			}

			return b.checkEpochDataDir(epoch, b.cfg.EpochLength)
		}

		return err
	}

	return nil
}

func (b *Bridge) isEventRemoved(event *contracts.TransferEvent) error {
	block, err := b.Client.BlockByNumber(context.Background(), big.NewInt(int64(event.Raw.BlockNumber)))
	if err != nil {
		return err
	}

	if block.Hash() != event.Raw.BlockHash {
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
