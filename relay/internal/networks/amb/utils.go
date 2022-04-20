package amb

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *Bridge) waitForTxMined(tx *types.Transaction) error {
	receipt, err := bind.WaitMined(context.Background(), b.Client, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		if err = b.GetFailureReason(tx); err != nil {
			return fmt.Errorf("GetFailureReason: %w", parseError(err))
		}
	}

	return nil
}

func (b *Bridge) getFailureReasonViaCall(funcName string, params ...interface{}) error {
	err := b.Contract.Raw().Call(&bind.CallOpts{
		From: b.Auth.From,
	}, nil, funcName, params...)

	if err != nil {
		return parseError(err)
	}
	return nil
}

func parseError(err error) error {
	if dataError, ok := err.(rpc.DataError); ok {
		errStr := dataError.ErrorData().(string)

		if errStr[:9] == "Reverted " {
			decodedMsg, err := decodeRevertMessage(errStr[9:]) // remove 'Reverted ' from string
			if err != nil {
				return errors.New(errStr)
			}
			return errors.New(decodedMsg)
		}

		return fmt.Errorf("%s %s", dataError, errStr) // e.g. "VM execution error. Out of gas"
	}
	return err
}

func decodeRevertMessage(errStr string) (string, error) {
	if len(errStr) < 138 {
		return fmt.Sprintf("probably runtime error: %s", errStr), nil
	}

	res := errStr[138:] // https://github.com/authereum/eth-revert-reason/blob/e33f4df82426a177dbd69c0f97ff53153592809b/index.js#L93
	res = strings.TrimRight(res, "0")
	// If the res is an odd number of characters, add a trailing 0
	if len(res)%2 == 1 {
		res = "0" + res
	}

	resBytes, err := hex.DecodeString(res)
	if err != nil {
		return "", err
	}
	return string(resBytes), nil
}
