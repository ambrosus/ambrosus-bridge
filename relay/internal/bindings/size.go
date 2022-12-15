package bindings

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var BridgeParsedABI = parseABI()

func parseABI() abi.ABI {
	abi, err := abi.JSON(strings.NewReader(BridgeMetaData.ABI))
	if err != nil {
		panic(fmt.Errorf("failed to parse bridge abi: %w", err))
	}
	return abi
}
