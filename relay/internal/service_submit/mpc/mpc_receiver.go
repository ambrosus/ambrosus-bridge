package mpc

import (
	"context"
	"fmt"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type MpcSigner interface {
	Sign(ctx context.Context, msg []byte) ([]byte, error)
}

type MpcReceiver struct {
	service_submit.Receiver
	mpcSigner MpcSigner

	signer types.Signer
}

func NewMpcReceiver(receiver service_submit.Receiver, mpcSigner MpcSigner) (*MpcReceiver, error) {
	chainID, err := receiver.GetClient().ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get chain id: %w", err)
	}
	signer := types.LatestSignerForChainID(chainID)

	return &MpcReceiver{
		Receiver:  receiver,
		mpcSigner: mpcSigner,
		signer:    signer,
	}, nil
}

func (b *MpcReceiver) GetAuth() *bind.TransactOpts {
	originalAuth := *b.Receiver.GetAuth()
	originalAuth.Signer = b.MpcSign
	// todo set originalAuth.Address to b.mpgSigner.GetAddress()
	return &originalAuth
}

func (b *MpcReceiver) MpcSign(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	hash := b.signer.Hash(tx).Bytes()
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	sig, err := b.mpcSigner.Sign(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("mpcSigner sign: %w", err)
	}
	return tx.WithSignature(b.signer, sig)
}
