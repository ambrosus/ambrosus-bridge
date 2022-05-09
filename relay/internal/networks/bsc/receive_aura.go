package bsc

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/common"
)

func (b *Bridge) SubmitTransferAura(proof *contracts.CheckAuraAuraProof) error {
	// Metric
	defer b.SetRelayBalanceMetric()

	tx, txErr := b.Contract.SubmitTransferAura(b.Auth, *proof)
	return b.ProcessTx(networks.GetTxErrParams{Tx: tx, TxErr: txErr, MethodName: "submitTransferAura", TxParams: []interface{}{*proof}})
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
