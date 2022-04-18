package eth

import (
	"context"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/metric"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func (b *Bridge) SubmitTransferAura(proof *contracts.CheckAuraAuraProof) error {
	// Metric
	defer metric.SetContractBalance(BridgeName, b.Client, b.auth.From)

	tx, txErr := b.Contract.SubmitTransferAura(b.auth, *proof)
	if txErr != nil {
		if txErr.Error() == "execution reverted" {
			dataErr := txErr.(rpc.DataError)
			return fmt.Errorf("contract runtime error: %s", dataErr.ErrorData())
		}
		return txErr
	}

	receipt, err := bind.WaitMined(context.Background(), b.Client, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		// we've got here probably due to low gas limit,
		// and revert() that hasn't been caught at eth_estimateGas
		err = ethereum.GetFailureReason(b.Client, b.auth, tx)
		if err != nil {
			return fmt.Errorf("GetFailureReason: %w", err)
		}
	}

	return nil
}

func (b *Bridge) GetValidatorSet() ([]common.Address, error) {
	return b.Contract.GetValidatorSet(nil)
}

func (b *Bridge) GetLastProcessedBlockHash() (*common.Hash, error) {
	blockHash, err := b.Contract.LastProcessedBlock(nil)
	if err != nil {
		return nil, err
	}

	return (*common.Hash)(&blockHash), nil
}
