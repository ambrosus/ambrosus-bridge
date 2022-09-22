package untrustless2

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiverUntrustless2 struct {
	service_submit.Receiver
}

func (b *ReceiverUntrustless2) SubmitTransferUntrustless(event *bindings.BridgeTransfer) error {
	defer metric.SetRelayBalanceMetric(b)

	// todo: make the multiplier configurable
	authChangedGasLimit, err := cb.GetMultipliedEstimatedGasLimit(*b.GetAuth(), 3, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SubmitTransferUntrustless(opts, event.EventId, event.Queue)
	})
	if err != nil {
		return fmt.Errorf("get multiplied estimated gas limit: %w", err)
	}

	return b.ProcessTx("submitTransferUntrustless", authChangedGasLimit, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SubmitTransferUntrustless(opts, event.EventId, event.Queue)
	})

}
func (b *ReceiverUntrustless2) IsEventAlreadyConfirmed(event *bindings.BridgeTransfer) (bool, error) {
	return b.GetContract().IsConfirmedByRelay(nil, b.GetAuth().From, event.EventId, event.Queue)
}
