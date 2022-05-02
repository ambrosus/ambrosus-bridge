package bsc

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ethereum/go-ethereum/core/types"
)

// todo
func (b *Bridge) EncodeBlock(header *types.Header) (*contracts.CheckPoSABlockPoSA, error) {
	return &contracts.CheckPoSABlockPoSA{}, nil
}
