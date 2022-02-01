package eth

import (
	"math/big"
	"relay/config"
	"relay/contracts"
	"relay/networks"

	"github.com/ethereum/go-ethereum/ethclient"
)

// не дописано

type Bridge struct {
	client   *ethclient.Client
	contract *contracts.Eth
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
		client:   client,
		contract: ethBridge,
	}
}

func (b *Bridge) SubmitBlockPoA(eventId *big.Int, blocks []*contracts.CheckPoABlockPoA, events *[]contracts.CommonStructsTransfer, proof *contracts.ReceiptsProof) {
	//tx, err := contracts.Submit(nil, withdraw.Blocks, withdraw.Events, withdraw.ReceiptsProof)
}

func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.contract.InputEventId(nil)
}

func (b *Bridge) Run(sideBridge networks.Bridge, submit networks.SubmitPoWF) {
	// todo watch events in eth
	// todo submit block on amb
}
