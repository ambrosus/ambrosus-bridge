package eth

import (
	"context"
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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const BridgeName = "ethereum"

type Bridge struct {
	networks.CommonBridge
	Config     *config.ETHConfig
	sideBridge networks.BridgeReceiveEthash
}

// New creates a new ethereum bridge.
func New(cfg *config.ETHConfig, externalLogger external_logger.ExternalLogger) (*Bridge, error) {
	logger := logger.NewSubLogger(BridgeName, externalLogger)

	logger.Debug().Msg("Creating ethereum bridge...")

	// Creating a new ethereum client (HTTP & WS).
	client, err := ethclient.Dial(cfg.HttpURL)
	if err != nil {
		return nil, fmt.Errorf("dial http: %w", err)
	}
	// Compatibility with tests.
	var wsClient *ethclient.Client
	if cfg.WsURL != "" {
		wsClient, err = ethclient.Dial(cfg.WsURL)
		if err != nil {
			return nil, fmt.Errorf("dial ws: %w", err)
		}
	}

	// Creating a new ambrosus bridge contract instance (HTTP & WS).
	contract, err := contracts.NewBridge(common.HexToAddress(cfg.ContractAddr), client)
	if err != nil {
		return nil, fmt.Errorf("create contract http: %w", err)
	}
	// Compatibility with tests.
	var wsContract *contracts.Bridge
	if wsClient != nil {
		wsContract, err = contracts.NewBridge(common.HexToAddress(cfg.ContractAddr), wsClient)
		if err != nil {
			return nil, fmt.Errorf("create contract ws: %w", err)
		}
	}

	var auth *bind.TransactOpts

	if cfg.PrivateKey != nil {
		chainId, err := client.ChainID(context.Background())
		if err != nil {
			return nil, fmt.Errorf("chain id: %w", err)
		}

		auth, err = bind.NewKeyedTransactorWithChainID(cfg.PrivateKey, chainId)
		if err != nil {
			return nil, fmt.Errorf("new keyed transactor: %w", err)
		}

		// Metric
		metric.SetContractBalance(BridgeName, client, auth.From)
	}

	b := &Bridge{
		CommonBridge: networks.CommonBridge{
			Client:      client,
			WsClient:    wsClient,
			Contract:    contract,
			WsContract:  wsContract,
			Auth:        auth,
			Logger:      logger,
		},
		Config: cfg,
	}
	b.CommonBridge.Bridge = b
	return b, nil
}

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.BridgeReceiveEthash) {
	b.Logger.Debug().Msg("Running ethereum bridge...")

	b.sideBridge = sideBridge
	b.CommonBridge.SideBridge = sideBridge

	// Getting last ethereum block number.
	blockNumber, err := b.Client.BlockNumber(context.Background())
	if err != nil {
		b.Logger.Error().Msgf("error getting last block number: %s", err.Error())
	}

	// Checking epoch data dir.
	if err = b.checkEpochDataDir(blockNumber/30000, b.Config.EpochLength); err != nil {
		b.Logger.Error().Msgf("error checking epoch data dir: %s", err.Error())
	}

	go b.UnlockOldestTransfersLoop()
	go b.WatchValidityLockedTransfersLoop()

	b.Logger.Info().Msg("Ethereum bridge runned!")

	for {
		if err := b.Listen(); err != nil {
			b.Logger.Error().Err(err).Msg("Listen error")
		}
	}
}

func (b *Bridge) SendEvent(event *contracts.BridgeTransfer) error {
	b.Logger.Debug().Str("event_id", event.EventId.String()).Msg("Waiting for safety blocks...")

	// Wait for safety blocks.
	safetyBlocks, err := b.sideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return fmt.Errorf("GetMinSafetyBlocksNum: %w", err)
	}

	if err := ethereum.WaitForBlock(b.WsClient, event.Raw.BlockNumber+safetyBlocks); err != nil {
		return fmt.Errorf("WaitForBlock: %w", err)
	}

	b.Logger.Debug().Str("event_id", event.EventId.String()).Msg("Checking if the event has been removed...")

	// Check if the event has been removed.
	if err := b.isEventRemoved(event); err != nil {
		return fmt.Errorf("isEventRemoved: %w", err)
	}

	ambTransfer, err := b.getBlocksAndEvents(event)
	if err != nil {
		return fmt.Errorf("getBlocksAndEvents: %w", err)
	}

	for _, blockNum := range []uint64{event.Raw.BlockNumber, event.Raw.BlockNumber + safetyBlocks} {
		if err := b.checkEpochData(blockNum, event.EventId); err != nil {
			return fmt.Errorf("checkEpochData on block %v: %w", blockNum, err)
		}
	}

	b.Logger.Debug().Str("event_id", event.EventId.String()).Msg("Submit transfer PoW...")
	err = b.sideBridge.SubmitTransferPoW(ambTransfer)
	if err != nil {
		return fmt.Errorf("SubmitTransferPoW: %w", err)
	}
	return nil
}

func (b *Bridge) GetTransactionError(params networks.GetTransactionErrorParams, txParams ...interface{}) error {
	if params.TxErr != nil {
		if params.TxErr.Error() == "execution reverted" {
			dataErr := params.TxErr.(rpc.DataError)
			return fmt.Errorf("contract runtime error: %s", dataErr.ErrorData())
		}
		return params.TxErr
	}

	receipt, err := bind.WaitMined(context.Background(), b.Client, params.Tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		// we've got here probably due to low gas limit,
		// and revert() that hasn't been caught at eth_estimateGas
		err = ethereum.GetFailureReason(b.Client, b.Auth, params.Tx)
		if err != nil {
			return fmt.Errorf("GetFailureReason: %w", err)
		}
	}
	return nil
}

func (b *Bridge) checkEpochData(blockNumber uint64, eventId *big.Int) error {
	epoch := blockNumber / 30000
	isEpochSet, err := b.sideBridge.IsEpochSet(epoch)
	if err != nil {
		return fmt.Errorf("IsEpochSet: %w", err)
	}
	if isEpochSet {
		return nil
	}

	b.Logger.Info().Str("event_id", eventId.String()).Msg("Submit epoch data...")
	epochData, err := b.loadEpochDataFile(epoch)
	if err != nil {
		return fmt.Errorf("loadEpochDataFile: %w", err)
	}

	err = b.sideBridge.SubmitEpochData(epochData)
	if err != nil {
		return fmt.Errorf("SubmitEpochData: %w", err)
	}
	return nil
	// todo delete old epochs, generate new (if need)
}

func (b *Bridge) isEventRemoved(event *contracts.BridgeTransfer) error {
	block, err := b.Client.BlockByNumber(context.Background(), big.NewInt(int64(event.Raw.BlockNumber)))
	if err != nil {
		return fmt.Errorf("HeaderByNumber: %w", err)
	}

	if block.Hash() != event.Raw.BlockHash {
		return fmt.Errorf("block hash != event's block hash")
	}
	return nil
}
