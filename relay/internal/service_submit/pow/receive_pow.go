package pow

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type ReceiverPoW struct {
	service_submit.Receiver
}

func (b *ReceiverPoW) SubmitTransferPoW(proof *bindings.CheckPoWPoWProof) error {
	defer metric.SetRelayBalanceMetric(b)

	return b.ProcessTx("submitTransferPoW", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SubmitTransferPoW(b.GetAuth(), *proof)
	})
}

func (b *ReceiverPoW) SubmitEpochData(epochData *ethash.EpochData) error {
	defer metric.SetRelayBalanceMetric(b)

	return b.ProcessTx("setEpochData", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.GetContract().SetEpochData(b.GetAuth(),
			epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth, epochData.MerkleNodes)
	})
}

func (b *ReceiverPoW) IsEpochSet(epoch uint64) (bool, error) {
	return b.GetContract().IsEpochDataSet(nil, big.NewInt(int64(epoch)))
}
