package amb

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *Bridge) waitForTxMined(tx *types.Transaction) error {
	receipt, err := bind.WaitMined(context.Background(), b.Client, tx)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		if err = b.getFailureReason(tx); err != nil {
			return fmt.Errorf("%s", parseError(err))
		}
	}

	return nil
}

func (b *Bridge) getFailureReason(tx *types.Transaction) error {
	_, err := b.Client.CallContract(context.Background(), ethereum.CallMsg{
		From:     b.auth.From,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}, nil)

	return err
}

type JsonError interface {
	ErrorData() interface{}
}

func parseError(err error) string {
	var jsonErr = err.(JsonError)
	errStr := jsonErr.ErrorData().(string)

	if strings.HasPrefix(errStr, "Reverted") {
		errBytes, err := hex.DecodeString(errStr[11:])
		if err == nil {
			return string(errBytes)
		}
	}

	return errStr
}
