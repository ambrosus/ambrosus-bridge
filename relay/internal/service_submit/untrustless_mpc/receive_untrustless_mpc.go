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
	Sign(ctx context.Context, party []string, msg []byte) ([]byte, error)
	SetFullMsg(fullMsg []byte)
	GetFullMsg() ([]byte, error)
	GetTssAddress() (common.Address, error)
}

type ReceiverUntrustlessMpc struct {
	service_submit.Receiver
	mpcSigner MpcSigner
	signer    types.Signer

	fromAddress common.Address
	signers     []string
	auth        *bind.TransactOpts
}

func NewReceiverUntrustlessMpc(receiver service_submit.Receiver, mpcSigner MpcSigner, signersIDs []string) (*ReceiverUntrustlessMpc, error) {
	chainID, err := receiver.GetClient().ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}
	signer := types.LatestSignerForChainID(chainID)

	auth := *receiver.GetAuth()
	auth.From, err = mpcSigner.GetTssAddress()
	if err != nil {
		return nil, fmt.Errorf("get tss address: %w", err)
	}
	auth.Signer = func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		hash := signer.Hash(tx).Bytes()
		txBytes, err := tx.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("tx to bytes: %w", err)
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Minute)
		mpcSigner.SetFullMsg(txBytes)
		sig, err := mpcSigner.Sign(ctx, signersIDs, hash)
		if err != nil {
			return nil, fmt.Errorf("mpcSigner sign: %w", err)
		}
		return tx.WithSignature(signer, sig)
	}

	return &ReceiverUntrustlessMpc{
		Receiver: receiver,
		auth:     &auth,
	}, nil
}

func (b *ReceiverUntrustlessMpc) GetAuth() *bind.TransactOpts {
	return b.auth
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

		retry.Delay(3*time.Second),
		retry.Attempts(20),
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
