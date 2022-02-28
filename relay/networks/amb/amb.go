package amb

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/networks"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

type Bridge struct {
	Client      *ethclient.Client
	Contract    *contracts.Amb
	VSContract  *contracts.Vs
	ContractRaw *contracts.AmbRaw

	HttpUrl    string // TODO: delete this field
	sideBridge networks.BridgeReceiveAura
	auth       *bind.TransactOpts
}

func (b *Bridge) SubmitEpochData() error {
	//TODO implement me
	panic("implement me")
}

// Creating a new ambrosus bridge.
func New(cfg *config.AMBConfig) (*Bridge, error) {
	// Creating a new ethereum client.
	client, err := ethclient.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	// Creating a new ambrosus bridge contract instance.
	contract, err := contracts.NewAmb(common.HexToAddress(cfg.ContractAddr), client)
	if err != nil {
		return nil, err
	}

	// Creating a new ambrosus VS contract instance.
	vsContract, err := contracts.NewVs(common.HexToAddress(cfg.VSContractAddr), client)
	if err != nil {
		return nil, err
	}

	bridge := &Bridge{
		Client:      client,
		Contract:    contract,
		ContractRaw: &contracts.AmbRaw{Contract: contract},
		VSContract:  vsContract,
		HttpUrl:     "https://network.ambrosus.io",
	}

	err = bridge.setAuth(cfg.PrivateKey)

	return bridge, err
}

func (b *Bridge) SubmitTransferPoW(proof *contracts.CheckPoWPoWProof) error {
	tx, txErr := b.Contract.SubmitTransfer(b.auth, *proof)

	// todo find way to make this part common for different contract methods
	if txErr != nil {
		// we've got here probably due to error at eth_estimateGas (e.g. revert(), require())
		// openethereum doesn't give us a full error message
		// so, make low-level call method to get the full error message

		err := b.ContractRaw.Call(&bind.CallOpts{
			From: b.auth.From,
		}, nil, "submitTransfer", *proof)

		if err != nil {
			return fmt.Errorf("%s", parseError(err))
		}
		return fmt.Errorf("%s", parseError(txErr))
	}

	return b.waitForTxMined(tx)
}

// Getting last contract event id.
func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

// Getting contract event by id.
func (b *Bridge) GetEventById(eventId *big.Int) (*contracts.TransferEvent, error) {
	opts := &bind.FilterOpts{Context: context.Background()}

	logs, err := b.Contract.FilterTransfer(opts, []*big.Int{eventId})
	if err != nil {
		return nil, err
	}

	if logs.Next() {
		return &logs.Event.TransferEvent, nil
	}

	return nil, networks.ErrEventNotFound
}

func (b *Bridge) Run(sideBridge networks.BridgeReceiveAura) {
	b.sideBridge = sideBridge

	for {
		if err := b.listen(); err != nil {
			log.Error().Err(err).Msg("listen ambrosus error")
		}
	}
}
