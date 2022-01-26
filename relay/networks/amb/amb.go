package amb

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"relay/config"
	"relay/contracts"
	"relay/networks"
)

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

func (b *Bridge) SubmitBlockPoW(eventId uint, blocks contracts.CheckPoWBlockPoW, events contracts.CommonStructsTransfer, proof contracts.ReceiptsProof) {
	// todo
}

func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.contract.InputEventId(nil)
}

func (b *Bridge) Run(sideBridge networks.Bridge, submit networks.SubmitPoWF) {
	// todo first start
	needId := sideBridge.GetLastEventId()
	for {
		needId += 1
		_ = needId

	}

	for {
		// todo listen

	}
}

func (b *Bridge) CheckOldEvents(sideBridge networks.Bridge, submit networks.SubmitPoWF) {
