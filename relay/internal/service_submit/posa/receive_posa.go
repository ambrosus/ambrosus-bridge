package posa

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiverPoSA struct {
	service_submit.Receiver
}

func (b *ReceiverPoSA) SubmitTransferPoSA(proof *bindings.CheckPoSAPoSAProof) error {
	defer metric.SetRelayBalanceMetric(b)

	return b.ProcessTx("submitTransferPoSA", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SubmitTransferPoSA(b.GetAuth(), *proof)
	})
}

func (b *ReceiverPoSA) SubmitValidatorSetChangesPoSA(proof *bindings.CheckPoSAPoSAProof) error {
	defer metric.SetRelayBalanceMetric(b)

	return b.ProcessTx("SubmitValidatorSetChangesPoSA", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SubmitValidatorSetChangesPoSA(b.GetAuth(), *proof)
	})
}

func (b *ReceiverPoSA) GetCurrentEpoch() (uint64, error) {
	epoch, err := b.GetContract().CurrentEpoch(nil)
	if err != nil {
		return 0, fmt.Errorf("CurrentEpoch: %w", err)
	}
	return epoch.Uint64(), nil
}
