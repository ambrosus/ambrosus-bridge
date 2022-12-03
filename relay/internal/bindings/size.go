package bindings

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var BridgeParsedABI = parseABI()

var ErrProofTooBig = errors.New("proof is too big")

type SizedProofs interface {
	*CheckAuraAuraProof | *CheckPoWPoWProof | *CheckPoSAPoSAProof
	Size() (uint64, error)
}

func IsProofTooBig[T SizedProofs](proof T, maxAllowedSize uint64) error {
	size, err := proof.Size()
	if err != nil {
		return err
	}

	if size >= maxAllowedSize {
		return ErrProofTooBig
	}
	return nil
}

func (p *CheckAuraAuraProof) Size() (uint64, error) {
	return getSize("submitTransferAura", *p)
}

func (p *CheckPoWPoWProof) Size() (uint64, error) {
	return getSize("submitTransferPoW", *p)
}

func (p *CheckPoSAPoSAProof) Size() (uint64, error) {
	return getSize("submitTransferPoSA", *p)
}

func getSize(methodName string, args ...interface{}) (uint64, error) {
	bytes, err := BridgeParsedABI.Pack(methodName, args...)
	if err != nil {
		return 0, err
	}

	return uint64(len(bytes)), nil
}

func parseABI() abi.ABI {
	abi, err := abi.JSON(strings.NewReader(BridgeMetaData.ABI))
	if err != nil {
		panic(fmt.Errorf("failed to parse bridge abi: %w", err))
	}
	return abi
}
