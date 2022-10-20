package common

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *CommonBridge) ProcessTx(methodName string, txOpts *bind.TransactOpts, txCallback networks.ContractCallFn) error {
	// if the transaction get stuck, then retry it with the higher gas price
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

			metric.IncTxCountMetric(b, methodName)

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
			var txHash string
			if tx != nil {
				txHash = tx.Hash().Hex()
			} else {
				txHash = "unknown"
			}

			b.Logger.Warn().
				Str("method", methodName).
				Str("tx_hash", txHash).
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

		// add useful info about insufficient funds error
		if isInsufficientFundsErr(err) {
			// bsc returns this info anyway, but I rather add this comment than another if statement
			balance, err := b.GetClient().BalanceAt(context.Background(), b.GetAuth().From, nil)
			if err != nil {
				return fmt.Errorf("get balance error: %w", err)
			}
			if tx == nil {
				return fmt.Errorf("%w. (Balance: %v; Require unknown)", err, balance)
			} else {
				return fmt.Errorf("%w. (Balance: %v; Require %v)", err, balance, tx.Cost())
			}
		}

		return err
	}

	txGasPrice, err := GetTxGasPrice(b.Client, tx)
	if err != nil {
		return fmt.Errorf("get tx gas price: %w", err)
	}
	metric.SetUsedGasMetric(b, methodName, receipt.GasUsed, txGasPrice)

	if receipt.Status != types.ReceiptStatusSuccessful {
		metric.IncFailedTxCountMetric(b, methodName)
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

func isInsufficientFundsErr(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "insufficient funds")
}
