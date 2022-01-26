package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"relay/config"
	"relay/contracts"
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
	ethBridge, err := contracts.NewEth(c.BridgeContractAddress, client)
	if err != nil {
		panic(err)
	}
	return &Bridge{
		client:   client,
		contract: ethBridge,
	}
}

func (b *Bridge) SubmitBlockPoW(eventId uint, blocks contracts.CheckPoWBlockPoW, events contracts.CommonStructsTransfer, proof contracts.ReceiptsProof) {
	tx, err := contracts.Submit(nil, withdraw.Blocks, withdraw.Events, withdraw.ReceiptsProof)
}

func (b *Bridge) GetLastEventId() {
	// todo
}
