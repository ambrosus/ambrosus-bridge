package amb_faucet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

// AmbFaucet sends money to the users that receive tokens and have zero balance
type AmbFaucet struct {
	networks.Bridge
	prev         service_submit.Submitter
	moneyAccount *bind.TransactOpts
	minBalance   *big.Int
	sendAmount   *big.Int
	logger       *zerolog.Logger
}

func NewAmbFaucet(prev service_submit.Submitter, moneyAccountPK string, minBalance, sendAmount *big.Int) *AmbFaucet {
	logger := prev.GetLogger().With().Str("service", "AmbFaucet").Logger()
	if prev.Receiver().GetName() != "ambrosus" {
		logger.Fatal().Msg("AmbFaucet can be used only with ambrosus receiver")
	}
	moneyAccount, err := createAuth(prev.Receiver().GetClient(), moneyAccountPK)
	if err != nil {
		logger.Fatal().Err(err).Msg("Create auth error")
	}

	return &AmbFaucet{
		Bridge:       prev,
		prev:         prev,
		moneyAccount: moneyAccount,
		minBalance:   minBalance,
		sendAmount:   sendAmount,
		logger:       &logger,
	}
}

func (b *AmbFaucet) Receiver() service_submit.Receiver {
	return b.prev.Receiver()
}

func (b *AmbFaucet) SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error {
	prevRes := b.prev.SendEvent(event, safetyBlocks)

	for _, t := range event.Queue {
		balance, err := b.Receiver().GetClient().BalanceAt(context.Background(), t.ToAddress, nil)
		if err != nil {
			b.logger.Error().Err(err).Str("address", t.ToAddress.String()).Msg("Get balance error")
			continue
		}
		if balance.Cmp(b.minBalance) != -1 {
			b.logger.Debug().Str("address", t.ToAddress.String()).Str("balance", balance.String()).Msg("User have enough balance")
			continue
		}

		b.logger.Info().Str("address", t.ToAddress.String()).Str("balance", balance.String()).Msg("User have not enough balance, sending money")
		tx, err := b.Transfer(t.ToAddress)
		if err != nil {
			b.logger.Error().Err(err).Str("address", t.ToAddress.String()).Msg("Send money error")
			continue
		}
		b.logger.Info().Str("address", t.ToAddress.String()).Str("tx", tx.Hash().String()).Msg("Money sent")
	}

	return prevRes
}

func (b *AmbFaucet) Transfer(addressTo common.Address) (*types.Transaction, error) {
	client := b.Receiver().GetClient()
	nonce, err := client.PendingNonceAt(context.Background(), b.moneyAccount.From)
	if err != nil {
		return nil, err
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &addressTo,
		Value:    b.sendAmount,
		Gas:      21_000,
		GasPrice: big.NewInt(0),
		Data:     []byte(nil),
	})

	signedTx, err := b.moneyAccount.Signer(b.moneyAccount.From, tx)
	if err != nil {
		return nil, err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	return signedTx, err
}

func createAuth(client ethclients.ClientInterface, privateKey string) (*bind.TransactOpts, error) {
	pk, err := helpers.ParsePK(privateKey)
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainId)
	if err != nil {
		return nil, fmt.Errorf("new keyed transactor: %w", err)
	}
	return auth, nil
}
