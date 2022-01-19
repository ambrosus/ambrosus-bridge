package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"relay/config"
	"relay/contracts"
	common2 "relay/networks/common"
)

type Bridge struct {
	client   *ethclient.Client
	contract *contracts.Eth
}

func NewEthBridge(c config.Network) *Bridge {
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

func (b *Bridge) SendWithdraw(withdraw *common2.Withdraw) {
	tx, err := ethBridge.TestAll(nil, withdraw.Blocks, withdraw.Events, withdraw.ReceiptsProof)
}

func (b *Bridge) GetLastEventId() {
	// todo
}
