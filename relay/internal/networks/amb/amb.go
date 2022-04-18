package amb

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
	"github.com/ethereum/go-ethereum/ethclient"
)

const BridgeName = "ambrosus"

type Bridge struct {
	networks.CommonBridge
	VSContract *contracts.Vs
	Config     *config.AMBConfig
	sideBridge networks.BridgeReceiveAura
}

// New creates a new ambrosus bridge.
func New(cfg *config.AMBConfig, externalLogger external_logger.ExternalLogger) (*Bridge, error) {
	logger := logger.NewSubLogger(BridgeName, externalLogger)

	logger.Debug().Msg("Creating ambrosus bridge...")

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

	// Creating a new ambrosus VS contract instance.
	vsContract, err := contracts.NewVs(common.HexToAddress(cfg.VSContractAddr), client)
	if err != nil {
		return nil, fmt.Errorf("create vs contract: %w", err)
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
		VSContract: vsContract,
		Config:     cfg,
	}
	b.CommonBridge.Bridge = b
	return b, nil
}

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.BridgeReceiveAura) {
	b.sideBridge = sideBridge
	b.CommonBridge.SideBridge = sideBridge

	go b.UnlockOldestTransfersLoop()
	go b.WatchValidityLockedTransfersLoop()

	b.Logger.Info().Msg("Ambrosus bridge runned!")

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

	b.Logger.Debug().Str("event_id", event.EventId.String()).Msg("Submit transfer Aura...")

	err = b.sideBridge.SubmitTransferAura(ambTransfer)
	if err != nil {
		return fmt.Errorf("SubmitTransferAura: %w", err)
	}
	return nil
}

func (b *Bridge) GetTransactionError(params networks.GetTransactionErrorParams, txParams ...interface{}) error {
	if params.TxErr != nil {
		// we've got here probably due to error at eth_estimateGas (e.g. revert(), require())
		// openethereum doesn't give us a full error message
		// so, make low-level call method to get the full error message
		err := b.getFailureReasonViaCall(params.MethodName, txParams...)
		if err != nil {
			return fmt.Errorf("getFailureReasonViaCall: %w", err)
		}
		return params.TxErr
	}

	err := b.waitForTxMined(params.Tx)
	if err != nil {
		return fmt.Errorf("waitForTxMined: %w", err)
	}
	return nil
}

func (b *Bridge) isEventRemoved(event *contracts.BridgeTransfer) error {
	block, err := b.HeaderByNumber(big.NewInt(int64(event.Raw.BlockNumber)))
	if err != nil {
		return fmt.Errorf("HeaderByNumber: %w", err)
	}

	if block.Hash(true) != event.Raw.BlockHash {
		return fmt.Errorf("block hash != event's block hash")
	}
	return nil
}
