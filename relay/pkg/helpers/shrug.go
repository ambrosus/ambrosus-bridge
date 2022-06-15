package helpers

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func ParsePK(pk string) (*ecdsa.PrivateKey, error) {
	if pk == "" {
		return nil, fmt.Errorf("parsePk: empty pk string")
	}
	b, err := hex.DecodeString(pk)
	if err != nil {
		return nil, fmt.Errorf("parsePk: decode Hex: %w", err)
	}
	p, err := crypto.ToECDSA(b)
	if err != nil {
		return nil, fmt.Errorf("parsePk: ToECDSA: %w", err)
	}
	return p, nil
}
