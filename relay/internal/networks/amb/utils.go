package amb

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *Bridge) waitForTxMined(tx *types.Transaction) error {
	receipt, err := bind.WaitMined(context.Background(), b.Client, tx)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		if err = ethereum.GetFailureReason(b.Client, b.auth, tx); err != nil {
			return parseError(err)
		}
	}

	return nil
}

func (b *Bridge) getFailureReasonViaCall(txErr error, funcName string, params ...interface{}) error {
	err := b.ContractRaw.Call(&bind.CallOpts{
		From: b.auth.From,
	}, nil, funcName, params...)

	if err != nil {
		return parseError(err)
	}
	return txErr
}

func parseError(err error) error {
	if dataError, ok := err.(rpc.DataError); ok {
		errStr := dataError.ErrorData().(string)
		decodedMsg, err := decodeRevertMessage(errStr[9:]) // remove 'Reverted ' from string
		if err != nil {
			return errors.New(errStr)
		}
		return errors.New(decodedMsg)
	}
	return err
}

func decodeRevertMessage(errStr string) (string, error) {
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
