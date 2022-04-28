package eth

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
)

// todo
func (b *Bridge) encodeCliqueProof(transferEvent *contracts.BridgeTransfer) (*contracts.CheckPoWPoWProof, error) {
	return &contracts.CheckPoWPoWProof{}, nil
}
