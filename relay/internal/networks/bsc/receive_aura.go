package bsc

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *Bridge) SubmitTransferAura(proof *bindings.CheckAuraAuraProof) error {
	// Metric
	defer b.SetRelayBalanceMetric()

	return b.ProcessTx("submitTransferAura", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.Contract.SubmitTransferAura(b.Auth, *proof)
	})
}

func (b *Bridge) GetValidatorSet() ([]common.Address, error) {
	return b.Contract.GetValidatorSet(nil)
}

func (b *Bridge) GetLastProcessedBlockHash() (*common.Hash, error) {
	blockHash, err := b.Contract.LastProcessedBlock(nil)
	if err != nil {
		return nil, err
	}

	return (*common.Hash)(&blockHash), nil
}

func (b *Bridge) GetMinSafetyBlocksValidators() (uint64, error) {
	v, err := b.Contract.MinSafetyBlocksValidators(nil)
	if err != nil {
		return 0, err
	}

	return v.Uint64(), nil
}
