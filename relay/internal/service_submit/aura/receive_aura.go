package aura

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiverAura struct {
	service_submit.Receiver
}

func (b *ReceiverAura) SubmitTransferAura(proof *bindings.CheckAuraAuraProof) error {
	defer metric.SetRelayBalanceMetric(b)

	return b.ProcessTx("submitTransferAura", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SubmitTransferAura(b.GetAuth(), *proof)
	})
}

func (b *ReceiverAura) SubmitValidatorSetChangesAura(proof *bindings.CheckAuraAuraProof) error {
	defer metric.SetRelayBalanceMetric(b)

	return b.ProcessTx("SubmitValidatorSetChangesAura", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SubmitValidatorSetChangesAura(b.GetAuth(), *proof)
	})
}

func (b *ReceiverAura) GetValidatorSet() ([]common.Address, error) {
	return b.GetContract().GetValidatorSet(nil)
}

func (b *ReceiverAura) GetLastProcessedBlockHash() (*common.Hash, error) {
	blockHash, err := b.GetContract().LastProcessedBlock(nil)
	if err != nil {
		return nil, err
	}

	return (*common.Hash)(&blockHash), nil
}

func (b *ReceiverAura) GetMinSafetyBlocksValidators() (uint64, error) {
	v, err := b.GetContract().MinSafetyBlocksValidators(nil)
	if err != nil {
		return 0, err
	}

	return v.Uint64(), nil
}
