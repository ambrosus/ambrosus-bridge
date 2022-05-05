package amb

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
)

func (b *Bridge) SubmitTransferPoSA(proof *contracts.CheckPoSAPoSAProof) error {
	defer b.SetRelayBalanceMetric()

	tx, txErr := b.Contract.SubmitTransferPoSA(b.Auth, *proof)
	return b.ProcessTx(networks.GetTxErrParams{Tx: tx, TxErr: txErr, MethodName: "submitTransferPoSA", TxParams: []interface{}{*proof}})
}
