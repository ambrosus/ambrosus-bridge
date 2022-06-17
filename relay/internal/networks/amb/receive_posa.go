package amb

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *Bridge) SubmitTransferPoSA(proof *bindings.CheckPoSAPoSAProof) error {
	defer b.SetRelayBalanceMetric()

	return b.ProcessTx("submitTransferPoSA", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.Contract.SubmitTransferPoSA(b.Auth, *proof)
	})
}

func (b *Bridge) GetCurrentEpoch() (uint64, error) {
	epoch, err := b.Contract.CurrentEpoch(nil)
	if err != nil {
		return 0, fmt.Errorf("CurrentEpoch: %w", err)
	}
	return epoch.Uint64(), nil
}
