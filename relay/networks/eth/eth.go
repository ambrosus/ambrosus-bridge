package eth

import (
	"crypto/ecdsa"
	"math/big"
	"relay/config"
	"relay/contracts"
	"relay/networks"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
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
	// todo вынести создание публичного ключа в отдельную ф-ию
	publicKey := b.config.PrivateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	opts := &bind.TransactOpts{
		From: fromAddress,
	}
	tx, err := b.contract.CheckPoA(opts, blocks, events, *proof)
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
