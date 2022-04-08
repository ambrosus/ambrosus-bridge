package amb

import (
	"fmt"
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
		err := b.getFailureReasonViaCall("submitTransfer", *proof)
		if err != nil {
			return fmt.Errorf("getFailureReasonViaCall: %w", err)
		}
		return txErr
	}

	err := b.waitForTxMined(tx)
	if err != nil {
		return fmt.Errorf("waitForTxMined: %w", err)
	}
	return nil
}

func (b *Bridge) SubmitEpochData(epochData *ethash.EpochData) error {
	// Metric
	defer metric.SetContractBalance(BridgeName, b.Client, b.auth.From)

	tx, txErr := b.Contract.SetEpochData(b.auth,
		epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth, epochData.MerkleNodes)
	if txErr != nil {
		err := b.getFailureReasonViaCall(
			"setEpochData",
			epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth, epochData.MerkleNodes,
		)
		if err != nil {
			return fmt.Errorf("getFailureReasonViaCall: %w", err)
		}
		return txErr
	}

	err := b.waitForTxMined(tx)
	if err != nil {
		return fmt.Errorf("waitForTxMined: %w", err)
	}
	return nil
}

func (b *Bridge) IsEpochSet(epoch uint64) (bool, error) {
	return b.Contract.IsEpochDataSet(nil, big.NewInt(int64(epoch)))
}
