package eth

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/networks"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Bridge struct {
	Client     *ethclient.Client
	Contract   *contracts.Eth
	sideBridge networks.Bridge
	config     *config.Bridge
}

func New(c *config.Bridge) *Bridge {
	client, err := ethclient.Dial(c.Url)
	if err != nil {
		panic(err)
	}
	ethBridge, err := contracts.NewEth(c.ContractAddress, client)
	if err != nil {
		panic(err)
	}
	return &Bridge{
		Client:   client,
		Contract: ethBridge,
		config:   c,
	}
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

func (b *Bridge) Run(sideBridge networks.Bridge) {
}
