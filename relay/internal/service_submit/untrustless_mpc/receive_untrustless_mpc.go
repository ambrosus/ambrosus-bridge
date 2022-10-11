package untrustless_mpc

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiverUntrustlessMpc struct {
	service_submit.Receiver
}

func (b *ReceiverUntrustlessMpc) SubmitTransferUntrustlessMpcServer(event *bindings.BridgeTransfer) error {
	defer metric.SetRelayBalanceMetric(b)

	return b.ProcessTx("submitTransferUntrustless", b.GetAuth(), func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SubmitTransferUntrustless(opts, event.EventId, event.Queue)
	})
}

func (b *ReceiverUntrustlessMpc) SubmitTransferUntrustlessMpcClient(event *bindings.BridgeTransfer) error {
	defer metric.SetRelayBalanceMetric(b)

	auth := *b.GetAuth()
	auth.NoSend = true
	_, err := b.GetContract().SubmitTransferUntrustless(&auth, event.EventId, event.Queue)
	return err
}
