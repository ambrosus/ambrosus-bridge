package amb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/parity"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

const BridgeName = "ambrosus"

type Bridge struct {
	nc.CommonBridge
	ParityClient *parity.Client
	vSContract   *bindings.Vs
	sideBridge   networks.BridgeReceiveAura
}

// New creates a new ambrosus bridge.
func New(cfg *config.AMBConfig, externalLogger logger.Hook) (*Bridge, error) {
	commonBridge, err := nc.New(cfg.Network, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}
	commonBridge.Logger = logger.NewSubLogger(BridgeName, externalLogger)

	// ///////////////////

	parityClient, err := parity.Dial(cfg.HttpURL)
	if err != nil {
		return nil, fmt.Errorf("dial http: %w", err)
	}
	commonBridge.Client = parityClient

	// Creating a new bridge contract instance.
	commonBridge.Contract, err = bindings.NewBridge(common.HexToAddress(cfg.ContractAddr), commonBridge.Client)
	if err != nil {
		return nil, fmt.Errorf("create contract http: %w", err)
	}

	// Create websocket instances if wsUrl provided
	if cfg.WsURL != "" {
		commonBridge.WsClient, err = parity.Dial(cfg.WsURL)
		if err != nil {
			return nil, fmt.Errorf("dial ws: %w", err)
		}

		commonBridge.WsContract, err = bindings.NewBridge(common.HexToAddress(cfg.ContractAddr), commonBridge.WsClient)
		if err != nil {
			return nil, fmt.Errorf("create contract ws: %w", err)
		}
	}

	// Creating a new ambrosus VS contract instance.
	vsContract, err := bindings.NewVs(common.HexToAddress(cfg.VSContractAddr), commonBridge.Client)
	if err != nil {
		return nil, fmt.Errorf("create vs contract: %w", err)
	}

	b := &Bridge{
		CommonBridge: commonBridge,
		ParityClient: parityClient,
		vSContract:   vsContract,
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
	go b.TriggerTransfersLoop()
	b.SubmitTransfersLoop()
}

func (b *Bridge) SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error {
	auraProof, err := b.encodeAuraProof(event, safetyBlocks)
	if err != nil {
		return fmt.Errorf("encodeAuraProof: %w", err)
	}

	proofJson, err := json.MarshalIndent(auraProof, "", "  ")
	if err != nil {
		return err
	}
	ioutil.WriteFile("auraproofdevnet.json", proofJson, 0644)

	b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer Aura...")
	err = b.sideBridge.SubmitTransferAura(auraProof)

	if isRequestTooBig(err) {
		b.Logger.Info().Str("event_id", event.EventId.String()).Msg("The proof is big, so split it")

		changes := splitVsChanges(auraProof)

		for i := 0; i < len(changes)-1; i++ {
			b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Submit validator set change...")
			err = b.sideBridge.SubmitValidatorSetChanges(changes[i])
			if err != nil {
				return fmt.Errorf("SubmitValidatorSetChanges: %w", err)
			}
		}

		err = b.sideBridge.SubmitTransferAura(changes[len(changes)-1])
	}

	if err != nil {
		return fmt.Errorf("SubmitTransferAura: %w", err)
	}
	return nil
}

func isRequestTooBig(err error) bool {
	if err == nil {
		return false
	}

	// TODO: we don't know is there this error in eth and bsc networks
	// if err.Error() == "Transaction is too big, see chain specification for the limit." {
	// 	return true
	// }

	if castedErr, ok := err.(rpc.HTTPError); ok {
		if ok && castedErr.StatusCode == http.StatusRequestEntityTooLarge {
			return true
		}
	}

	return false
}
