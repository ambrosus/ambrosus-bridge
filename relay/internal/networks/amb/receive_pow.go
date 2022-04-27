package amb

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
)

func (b *Bridge) SubmitTransferPoW(proof *contracts.CheckPoWPoWProof) error {
	defer b.SetRelayBalanceMetric()

	tx, txErr := b.Contract.SubmitTransferPoW(b.Auth, *proof)
	txErr = b.ProcessTx(networks.GetTxErrParams{Tx: tx, TxErr: txErr, MethodName: "submitTransferPoW", TxParams: []interface{}{*proof}})
	if txErr == nil {
		b.AddWithdrawalsCountMetric(len(proof.Transfer.Transfers))
	}
	return txErr
}

func (b *Bridge) SubmitEpochData(epochData *ethash.EpochData) error {
	defer b.SetRelayBalanceMetric()

	tx, txErr := b.Contract.SetEpochData(b.Auth,
		epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth, epochData.MerkleNodes)

	return b.ProcessTx(networks.GetTxErrParams{Tx: tx, TxErr: txErr, MethodName: "setEpochData", TxParams: []interface{}{
		epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth, epochData.MerkleNodes,
	}})
}

func (b *Bridge) IsEpochSet(epoch uint64) (bool, error) {
	return b.Contract.IsEpochDataSet(nil, big.NewInt(int64(epoch)))
}
