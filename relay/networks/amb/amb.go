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

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.Bridge, submit networks.SubmitPoWF) {
	// todo save args to struct?
	b.CheckOldEvents()
	b.Listen()
}

func (b *Bridge) CheckOldEvents() {
	for {
		needId := sideBridge.GetLastEventId() + 1
		// todo get event by id `needId`

		if !event {
			return
		}

		b.sendEvent()
	}
}

func (b *Bridge) Listen() {
	// todo listen
	b.sendEvent(event)
}

func (b *Bridge) sendEvent(event) {
	// todo wait for safety blocks
	// todo encode blocks
	// todo estimate gas
	// todo send
	// todo wait status ok
}
