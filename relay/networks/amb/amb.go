package amb

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"relay/config"
	"relay/contracts"
	"relay/networks/common"
)

type Bridge struct {
	client   *ethclient.Client
	contract *contracts.Eth
}

func New(c config.Network) *Bridge {
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

func (b *Bridge) SendDeposit(deposit *common.Deposit) {
	//tx, err := b.contract.TestAll(nil, deposit.Blocks, deposit.Events, deposit.ReceiptsProof)
}

func (b *Bridge) GetLastEventId() {
	// todo
}
