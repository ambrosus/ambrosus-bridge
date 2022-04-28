package eth

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ethereum/go-ethereum/core/types"
)

// todo
func (b *Bridge) EncodeBlock(header *types.Header) (*contracts.CheckPoWBlockPoW, error) {
	return &contracts.CheckPoWBlockPoW{}, nil
}
