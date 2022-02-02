package eth

import (
	"math/big"
	"relay/config"
	"relay/contracts"
	"relay/networks"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

// не дописано

type Bridge struct {
	client   *ethclient.Client
	contract *contracts.Eth
	config   *config.Bridge
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
		config:   c,
	}
}

func (b *Bridge) SubmitBlockPoA(eventId *big.Int, blocks []contracts.CheckPoABlockPoA, events []contracts.CommonStructsTransfer, proof *contracts.ReceiptsProof) {
	auth, err := b.getAuth()
	if err != nil {
		// todo
	}

	tx, err := b.contract.CheckPoA(auth, blocks, events, *proof)
	if err != nil {
		// todo
	}

}

func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.contract.InputEventId(nil)
}

func (b *Bridge) Run(sideBridge networks.Bridge, submit networks.SubmitPoWF) {
	// todo watch events in eth
	// todo submit block on amb
}

func (b Bridge) getAuth() (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(b.config.PrivateKey, b.config.ChainID)
	if err != nil {
		return nil, err
	}

	nonce, err := b.client.PendingNonceAt(auth.Context, auth.From)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	return auth, nil
}
