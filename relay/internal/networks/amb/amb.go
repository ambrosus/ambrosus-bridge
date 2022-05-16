package amb

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/parity"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/price"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const BridgeName = "ambrosus"

type Bridge struct {
	nc.CommonBridge
	Client     *parity.Client
	WsClient   *parity.Client
	VSContract *contracts.Vs
	sideBridge networks.BridgeReceiveAura
}

// New creates a new ambrosus bridge.
func New(cfg *config.AMBConfig, externalLogger logger.Hook) (*Bridge, error) {
	commonBridge, err := nc.New(cfg.Network, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}
	commonBridge.Logger = logger.NewSubLogger(BridgeName, externalLogger)

	// ///////////////////

	client, err := parity.Dial(cfg.HttpURL)
	if err != nil {
		return nil, fmt.Errorf("dial http: %w", err)
	}

	// Creating a new ambrosus VS contract instance.
	vsContract, err := contracts.NewVs(common.HexToAddress(cfg.VSContractAddr), client)
	if err != nil {
		return nil, fmt.Errorf("create vs contract: %w", err)
	}

	// Creating a new bridge contract instance.
	commonBridge.Contract, err = contracts.NewBridge(common.HexToAddress(cfg.ContractAddr), client)
	if err != nil {
		return nil, fmt.Errorf("create contract http: %w", err)
	}

	// Create websocket instances if wsUrl provided
	var wsClient *parity.Client
	if cfg.WsURL != "" {
		wsClient, err = parity.Dial(cfg.WsURL)
		if err != nil {
			return nil, fmt.Errorf("dial ws: %w", err)
		}

		commonBridge.WsContract, err = contracts.NewBridge(common.HexToAddress(cfg.ContractAddr), wsClient)
		if err != nil {
			return nil, fmt.Errorf("create contract ws: %w", err)
		}
	}

	b := &Bridge{
		CommonBridge: commonBridge,
		Client:       client,
		WsClient:     wsClient,
		VSContract:   vsContract,
	}
	b.CommonBridge.Bridge = b
	return b, nil
}

func (b *Bridge) SetSideBridge(sideBridge networks.BridgeReceiveAura) {
	b.sideBridge = sideBridge
	b.CommonBridge.SideBridge = sideBridge
}

func (b *Bridge) Run() {
	b.Logger.Debug().Msg("Running ambrosus bridge...")

	go b.UnlockTransfersLoop()
	b.SubmitTransfersLoop()
}

func (b *Bridge) SendEvent(event *contracts.BridgeTransfer, safetyBlocks uint64) error {
	auraProof, err := b.encodeAuraProof(event, safetyBlocks)
	if err != nil {
		return fmt.Errorf("encodeAuraProof: %w", err)
	}

	b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer Aura...")
	err = b.sideBridge.SubmitTransferAura(auraProof)
	if err != nil {
		return fmt.Errorf("SubmitTransferAura: %w", err)
	}
	return nil
}

func (b *Bridge) GetTxErr(params networks.GetTxErrParams) error {
	if params.TxErr != nil {
		// we've got here probably due to error at eth_estimateGas (e.g. revert(), require())
		// openethereum doesn't give us a full error message at eth_estimateGas, so
		// do eth_call method to get the full error message
		err := b.Contract.Raw().Call(&bind.CallOpts{From: b.Auth.From}, nil, params.MethodName, params.TxParams...)
		if err != nil {
			return fmt.Errorf("getFailureReasonViaCall: %w", helpers.ParseError(err))
		}
		return params.TxErr
	}
	return nil
}

func (b *Bridge) CoinPrice() (float64, error) {
	return price.CoinToUsdt(price.Amb)
}
