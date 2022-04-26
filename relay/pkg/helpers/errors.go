package helpers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
)

func ParseError(txErr error) error {
	dataError, ok := txErr.(rpc.DataError)
	if !ok {
		return txErr
	}

	if txErr.Error() == "execution reverted" {
		return fmt.Errorf("contract runtime error: %s", dataError.ErrorData())
	}

	errorData := decodeErrorData(dataError.ErrorData())
	return fmt.Errorf("%s %s", dataError, errorData) // e.g. "VM execution error. Out of gas"
}

func decodeErrorData(errData interface{}) error {
	if errStr, ok := errData.(string); ok {

		// parity case, need to decode revert message
		if errStr[:9] == "Reverted " {
			decodedMsg, err := decodeRevertMessage(errStr[9:]) // remove 'Reverted ' from string
			if err != nil {
				return errors.New(errStr)
			}
			return errors.New(decodedMsg)
		}

		return errors.New(errStr)
	}
	return fmt.Errorf("")
}

func decodeRevertMessage(errStr string) (string, error) {
	// https://github.com/authereum/eth-revert-reason/blob/e33f4df82426a177dbd69c0f97ff53153592809b/index.js#L93
	// "0x08c379a0" is `Error(string)` method signature, it's called by revert/require
	if len(errStr) < 138 || errStr[:10] != "0x08c379a0" {
		return fmt.Sprintf("probably runtime error: %s", errStr), nil
	}

	res := strings.TrimRight(errStr[138:], "0")
	if len(res)%2 == 1 {
		res += "0" // If the res is an odd number of characters, add a trailing 0
	}

	resBytes, err := hex.DecodeString(res)
	if err != nil {
		return "", err
	}
	return string(resBytes), nil
}
