package untrustless_mpc

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type MpcSigner interface {
	Sign(ctx context.Context, msg []byte) ([]byte, error)
	SetFullMsg(fullMsg []byte)
	GetFullMsg() ([]byte, error)
	GetTssAddress() (common.Address, error)
}

type ReceiverUntrustlessMpc struct {
	service_submit.Receiver
	mpcSigner MpcSigner
	signer    types.Signer

	fromAddress common.Address
}

func NewReceiverUntrustlessMpc(receiver service_submit.Receiver, mpcSigner MpcSigner) (*ReceiverUntrustlessMpc, error) {
	chainID, err := receiver.GetClient().ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}
	signer := types.LatestSignerForChainID(chainID)

	fromAddress, err := mpcSigner.GetTssAddress()
	if err != nil {
		return nil, fmt.Errorf("get tss address: %w", err)
	}

	return &ReceiverUntrustlessMpc{
		Receiver:    receiver,
		mpcSigner:   mpcSigner,
		signer:      signer,
		fromAddress: fromAddress,
	}, nil
}

func (b *ReceiverUntrustlessMpc) GetAuth() *bind.TransactOpts {
	originalAuth := *b.Receiver.GetAuth()
	originalAuth.Signer = b.MpcSign
	originalAuth.From = b.fromAddress
	return &originalAuth
}

func (b *ReceiverUntrustlessMpc) MpcSign(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	hash := b.signer.Hash(tx).Bytes()
	txBytes, err := tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("tx to bytes: %w", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	b.mpcSigner.SetFullMsg(txBytes)
	sig, err := b.mpcSigner.Sign(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("mpcSigner sign: %w", err)
	}
	return tx.WithSignature(b.signer, sig)
}

func (b *ReceiverUntrustlessMpc) SubmitTransferUntrustlessMpcServer(event *bindings.BridgeTransfer) error {
	defer metric.SetRelayBalanceMetric(b)
	defer b.mpcSigner.SetFullMsg(nil)

	return b.ProcessTx("submitTransferUntrustless", b.GetAuth(), func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SubmitTransferUntrustless(opts, event.EventId, event.Queue)
	})
}

func (b *ReceiverUntrustlessMpc) SubmitTransferUntrustlessMpcClient(event *bindings.BridgeTransfer) error {
	serverTx, err := b.getServerTxRetryable()
	if err != nil {
		return fmt.Errorf("get server tx: %w", err)
	}

	auth := *b.GetAuth()
	auth.NoSend = true
	auth.Nonce = big.NewInt(int64(serverTx.Nonce()))
	auth.GasLimit = serverTx.Gas()

	if serverTx.Type() == types.LegacyTxType {
		auth.GasPrice = serverTx.GasPrice()
	} else {
		auth.GasFeeCap = serverTx.GasFeeCap()
		auth.GasTipCap = serverTx.GasTipCap()
	}

	//
	_, err = b.GetContract().SubmitTransferUntrustless(&auth, event.EventId, event.Queue)
	return err
}

func (b *ReceiverUntrustlessMpc) getServerTxRetryable() (*types.Transaction, error) {
	var tx *types.Transaction
	err := retry.Do(
		func() (err error) {
			tx, err = b.getServerTx()
			return err
		},

		retry.Delay(2*time.Second),
		retry.DelayType(retry.FixedDelay),
	)
	return tx, err
}

func (b *ReceiverUntrustlessMpc) getServerTx() (*types.Transaction, error) {
	txBytes, err := b.mpcSigner.GetFullMsg()
	if err != nil {
		return nil, err
	}
	var tx types.Transaction
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		return nil, fmt.Errorf("unmarshal binary tx: %w", err)
	}

	return &tx, nil
}
