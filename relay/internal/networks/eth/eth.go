package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/metric"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
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
	logger     zerolog.Logger
}

// New creates a new ethereum bridge.
func New(cfg *config.ETHConfig, externalLogger external_logger.ExternalLogger) (*Bridge, error) {
	logger := logger.NewSubLogger(BridgeName, externalLogger)

	logger.Debug().Msg("Creating ethereum bridge...")

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
		logger:     logger,
	}, nil
}

func (b *Bridge) Name() string {
	return BridgeName
}

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

func (b *Bridge) Run(sideBridge networks.BridgeReceiveEthash) {
	b.logger.Debug().Msg("Running ethereum bridge...")

	b.sideBridge = sideBridge

	// Getting last ethereum block number.
	blockNumber, err := b.Client.BlockNumber(context.Background())
	if err != nil {
		b.logger.Error().Msgf("error getting last block number: %s", err.Error())
	}

	// Checking epoch data dir.
	if err = b.checkEpochDataDir(blockNumber/30000, b.cfg.EpochLength); err != nil {
		b.logger.Error().Msgf("error checking epoch data dir: %s", err.Error())
	}

	b.logger.Info().Msg("Ethereum bridge runned!")

	for {
		if err := b.listen(); err != nil {
			b.logger.Error().Err(err)
		}
	}
}

func (b *Bridge) checkOldEvents() error {
	b.logger.Info().Msg("Checking old events...")

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

		b.logger.Info().Str("event_id", nextEventId.String()).Msg("Send old event...")

		if err := b.sendEvent(nextEvent); err != nil {
			return err
		}

		i = big.NewInt(0).Add(i, big.NewInt(1))
	}
}

func (b *Bridge) listen() error {
	if err := b.checkOldEvents(); err != nil {
		return err
	}

	b.logger.Info().Msg("Listening new events...")

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
			b.logger.Info().Str("event_id", event.EventId.String()).Msg("Send event...")

			if err := b.sendEvent(&event.TransferEvent); err != nil {
				return err
			}
		}
	}
}

func (b *Bridge) sendEvent(event *contracts.TransferEvent) error {
	b.logger.Debug().Str("event_id", event.EventId.String()).Msg("Waiting for safety blocks...")

	// Wait for safety blocks.
	safetyBlocks, err := b.sideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return err
	}

	if err := ethereum.WaitForBlock(b.WsClient, event.Raw.BlockNumber+safetyBlocks); err != nil {
		return err
	}

	b.logger.Debug().Str("event_id", event.EventId.String()).Msg("Checking if the event has been removed...")

	// Check if the event has been removed.
	if err := b.isEventRemoved(event); err != nil {
		return err
	}

	ambTransfer, err := b.getBlocksAndEvents(event)
	if err != nil {
		return err
	}

	for _, blockNum := range []uint64{event.Raw.BlockNumber, event.Raw.BlockNumber + safetyBlocks} {
		if err := b.checkEpochData(blockNum, event.EventId); err != nil {
			return err
		}
	}

	b.logger.Debug().Str("event_id", event.EventId.String()).Msg("Submit transfer PoW...")
	return b.sideBridge.SubmitTransferPoW(ambTransfer)
}

func (b *Bridge) checkEpochData(blockNumber uint64, eventId *big.Int) error {
	epoch := blockNumber / 30000
	isEpochSet, err := b.sideBridge.IsEpochSet(epoch)
	if err != nil {
		return err
	}
	if isEpochSet {
		return nil
	}

	b.logger.Info().Str("event_id", eventId.String()).Msg("Submit epoch data...")
	epochData, err := b.loadEpochDataFile(epoch)
	if err != nil {
		return err
	}
	return b.sideBridge.SubmitEpochData(epochData)
	// todo delete old epochs, generate new (if need)
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
