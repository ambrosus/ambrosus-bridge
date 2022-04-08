package amb

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/metric"
)

func (b *Bridge) SubmitTransferPoW(proof *contracts.CheckPoWPoWProof) error {
	// Metric
	defer metric.SetContractBalance(BridgeName, b.Client, b.auth.From)

	tx, txErr := b.Contract.SubmitTransfer(b.auth, *proof)

	if txErr != nil {
		// we've got here probably due to error at eth_estimateGas (e.g. revert(), require())
		// openethereum doesn't give us a full error message
		// so, make low-level call method to get the full error message
		return b.getFailureReasonViaCall(txErr, "submitTransfer", *proof)
	}

	return b.waitForTxMined(tx)
}

func (b *Bridge) SubmitEpochData(epochData *ethash.EpochData) error {
	// Metric
	defer metric.SetContractBalance(BridgeName, b.Client, b.auth.From)

	tx, txErr := b.Contract.SetEpochData(b.auth,
		epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth, epochData.MerkleNodes)
	if txErr != nil {
		return b.getFailureReasonViaCall(
			txErr,
			"setEpochData",
			epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth, epochData.MerkleNodes,
		)
	}

	return b.waitForTxMined(tx)
}

func (b *Bridge) IsEpochSet(epoch uint64) (bool, error) {
	return b.Contract.IsEpochDataSet(nil, big.NewInt(int64(epoch)))
}
