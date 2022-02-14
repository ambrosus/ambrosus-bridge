package eth

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/networks"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

type Bridge struct {
	Client     *ethclient.Client
	Contract   *contracts.Eth
	sideBridge networks.Bridge
	config     *config.Bridge
}

// Creating a new ethereum bridge.
func New(cfg *config.Bridge) *Bridge {
	// Creating a new ethereum client.
	client, err := ethclient.Dial(cfg.Url)
	if err != nil {
		log.Fatal().Err(err).Msg("ethereum client not created")
	}

	// Creating a new ethereum bridge contract instance.
	contract, err := contracts.NewEth(cfg.ContractAddress, client)
	if err != nil {
		log.Fatal().Err(err).Msg("ethereum bridge contract instance not created")
	}

	return &Bridge{Client: client, Contract: contract, config: cfg}
}

func (b *Bridge) SubmitTransfer(proof contracts.TransferProof) error {
	switch proof.(type) {
	case *contracts.CheckAuraAuraProof:
		// todo
	default:
		// todo error

	}
	return nil
}

func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.Bridge) {}
