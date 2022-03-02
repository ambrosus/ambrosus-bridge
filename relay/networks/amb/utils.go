package amb

import (
	"context"
	"encoding/hex"
	"fmt"
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
			return fmt.Errorf("%s", parseError(err))
		}
	}

	return nil
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
