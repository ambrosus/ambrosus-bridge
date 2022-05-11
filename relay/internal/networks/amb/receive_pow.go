package amb

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *Bridge) SubmitTransferPoW(proof *contracts.CheckPoWPoWProof) error {
	defer b.SetRelayBalanceMetric()

	return b.ProcessTx(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.Contract.SubmitTransferPoW(b.Auth, *proof)
	}, networks.GetTxErrParams{MethodName: "submitTransferPoW", TxParams: []interface{}{*proof}})
}

func (b *Bridge) SubmitEpochData(epochData *ethash.EpochData) error {
	defer b.SetRelayBalanceMetric()

	return b.ProcessTx(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.Contract.SetEpochData(b.Auth,
			epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth, epochData.MerkleNodes)
	}, networks.GetTxErrParams{MethodName: "setEpochData", TxParams: []interface{}{
		epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth, epochData.MerkleNodes,
	}})
}

func (b *Bridge) IsEpochSet(epoch uint64) (bool, error) {
	return b.Contract.IsEpochDataSet(nil, big.NewInt(int64(epoch)))
}
