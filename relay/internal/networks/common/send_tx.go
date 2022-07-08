package common

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *CommonBridge) ProcessTx(methodName string, txCallback networks.ContractCallFn) error {
	// if the transaction get stuck, then retry it with the higher gas price
	var txOpts = b.Auth
	var receipt *types.Receipt
	var tx *types.Transaction

	err := retry.Do(
		func() (err error) {
			b.ContractCallLock.Lock()
			tx, err = txCallback(txOpts)
			b.ContractCallLock.Unlock()
			if err != nil {
				return err
			}

			b.IncTxCountMetric(methodName)

			b.Logger.Info().
				Str("method", methodName).
				Str("tx_hash", tx.Hash().Hex()).
				Msgf("Wait the tx to be mined...")

			receipt, err = b.waitMined(tx)
			return err
		},

		retry.RetryIf(func(err error) bool {
			return errors.Is(err, context.DeadlineExceeded)
		}),
		retry.OnRetry(func(n uint, err error) {
			b.Logger.Warn().
				Str("method", methodName).
				Str("tx_hash", tx.Hash().Hex()).
				Msgf("Seems the transaction get stuck, making new tx with higher gas price and the same nonce to replace the old one... (%d/%d)", n+1, 2)

			// set gas price higher by 30%
			txOpts.GasPrice, _ = new(big.Float).Mul(
				new(big.Float).SetInt(tx.GasPrice()),
				big.NewFloat(1.30),
			).Int(nil)
			txOpts.Nonce = big.NewInt(int64(tx.Nonce()))
		}),
		retry.Attempts(2),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		return err
	}

	b.SetUsedGasMetric(methodName, receipt.GasUsed, tx.GasPrice())

	if receipt.Status != types.ReceiptStatusSuccessful {
		b.IncFailedTxCountMetric(methodName)
		if err = b.getFailureReason(tx); err != nil {
			return fmt.Errorf("tx %s failed: %w", tx.Hash().Hex(), helpers.ParseError(err))
		}
		b.Logger.Debug().Err(err).Str("tx_hash", tx.Hash().Hex()).Msg("Tx has been mined but failed :(")
	}
	b.Logger.Debug().Str("tx_hash", tx.Hash().Hex()).Msg("Tx has been mined successfully!")

	return nil
}

func (b *CommonBridge) waitMined(tx *types.Transaction) (receipt *types.Receipt, err error) {
	err = retry.Do(
		func() (err error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
			defer cancel()

			receipt, err = bind.WaitMined(ctx, b.Client, tx)
			return err
		},

		retry.RetryIf(func(err error) bool {
			return errors.Is(err, context.DeadlineExceeded)
		}),
		retry.OnRetry(func(n uint, err error) {
			b.Logger.Warn().
				Str("tx_hash", tx.Hash().Hex()).
				Msgf("Timeout waiting for tx to be mined, trying again... (%d/%d)", n+1, 2)
		}),
		retry.Attempts(2),
		retry.LastErrorOnly(true),
	)
	if err != nil {
		return nil, fmt.Errorf("wait mined: %w", err)
	}
	return
}

func (b *CommonBridge) getFailureReason(tx *types.Transaction) error {
	_, err := b.Client.CallContract(context.Background(), ethereum.CallMsg{
		From:     b.Auth.From,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}, nil)

	return err
}
