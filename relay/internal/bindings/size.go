package bindings

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func (p *CheckAuraAuraProof) Size() (int, error) {
	return getSize("submitTransferAura", *p)
}

func (p *CheckPoWPoWProof) Size() (int, error) {
	return getSize("submitTransferPoW", *p)
}

func (p *CheckPoSAPoSAProof) Size() (int, error) {
	return getSize("submitTransferPoSA", *p)
}

func getSize(methodName string, args ...interface{}) (int, error) {
	parsedAbi, err := parseABI()
	if err != nil {
		return 0, err
	}

	bytes, err := parsedAbi.Pack(methodName, args...)
	if err != nil {
		return 0, err
	}

	return len(bytes), nil
}

func parseABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(BridgeMetaData.ABI))
}
