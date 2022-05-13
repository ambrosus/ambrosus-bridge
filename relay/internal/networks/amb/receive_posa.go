package amb

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

func (b *Bridge) SubmitTransferPoSA(proof *contracts.CheckPoSAPoSAProof) error {
	defer b.SetRelayBalanceMetric()

	return b.ProcessTx(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.Contract.SubmitTransferPoSA(b.Auth, *proof)
	}, networks.GetTxErrParams{MethodName: "submitTransferPoSA", TxParams: []interface{}{*proof}})
}
