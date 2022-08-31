package untrustless

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiverUntrustless struct {
	service_submit.Receiver
}

func (b *ReceiverUntrustless) SubmitTransferUntrustless(event *bindings.BridgeTransfer) error {
	defer metric.SetRelayBalanceMetric(b)

	return b.ProcessTx("submitTransferUntrustless", b.GetAuth(), func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SubmitTransferUntrustless(opts, event.EventId, event.Queue)
	})

}
func (b *ReceiverUntrustless) IsEventAlreadyConfirmed(event *bindings.BridgeTransfer) (bool, error) {
	return b.GetContract().IsConfirmedByRelay(nil, b.GetAuth().From, event.EventId, event.Queue)
}
