package amb_faucet

import (
	"context"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

// AmbFaucet sends money to the users that receive tokens and have zero balance
type AmbFaucet struct {
	networks.Bridge // used for sending transactions, must be amb bridge
	prev            service_submit.Submitter
	faucetAddress   common.Address
	faucetContract  *bindings.Faucet
	minBalance      *big.Int
	sendAmount      *big.Int
	logger          *zerolog.Logger
}

func NewAmbFaucet(bridge networks.Bridge, prev service_submit.Submitter, faucetAddress common.Address, minBalance, sendAmount *big.Int) *AmbFaucet {
	logger := prev.GetLogger().With().Str("service", "AmbFaucet").Logger()
	if prev.Receiver().GetName() != "ambrosus" {
		logger.Fatal().Msg("AmbFaucet can be used only with ambrosus receiver")
	}
	if bridge.GetName() != "ambrosus" {
		logger.Fatal().Msg("AmbFaucet can be used only with ambrosus bridge")
	}

	faucetContract, err := bindings.NewFaucet(faucetAddress, bridge.GetClient())
	if err != nil {
		logger.Fatal().Err(err).Msg("Create faucet contract error")
	}

	return &AmbFaucet{
		Bridge:         bridge,
		prev:           prev,
		faucetAddress:  faucetAddress,
		faucetContract: faucetContract,
		minBalance:     minBalance,
		sendAmount:     sendAmount,
		logger:         &logger,
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
		if balance.Cmp(b.minBalance) == 1 { // do nothing if balance > minBalance
			b.logger.Debug().Str("address", t.ToAddress.String()).Str("balance", balance.String()).Msg("User have enough balance")
			continue
		}

		b.logger.Info().Str("address", t.ToAddress.String()).Str("balance", balance.String()).Msg("User have not enough balance, sending money")
		tx, err := b.Transfer(t.ToAddress, event.EventId)
		if err != nil {
			b.logger.Error().Err(err).Str("address", t.ToAddress.String()).Msg("Send money error")
			continue
		}
		b.logger.Info().Str("address", t.ToAddress.String()).Str("tx", tx.Hash().String()).Msg("Money sent")
	}

	return prevRes
}

func (b *AmbFaucet) Transfer(addressTo common.Address, eventId *big.Int) (*types.Transaction, error) {
	defer metric.SetAmbFaucetBalanceMetric(b.Bridge, b.faucetAddress)
	return b.faucetContract.Faucet(b.Bridge.GetAuth(), addressTo, eventId, b.sendAmount)
}
