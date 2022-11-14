package helpers

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/exp/constraints"
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

func NewCache[K comparable, V any](getter func(K) (V, error)) func(arg K) (V, error) {
	cache := map[K]V{}
	return func(arg K) (V, error) {
		if v, ok := cache[arg]; ok {
			return v, nil
		}

		res, err := getter(arg)
		if err == nil {
			cache[arg] = res
		}

		return res, err
	}
}

func Range[T constraints.Integer](start, end T) []T {
	res := make([]T, 0, end-start)
	for i := start; i < end; i++ {
		res = append(res, i)
	}
	return res
}
