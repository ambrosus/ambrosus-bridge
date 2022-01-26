package contracts

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// CommonStructsTransfer is an auto generated low-level Go binding around an user-defined struct.
type CommonStructsTransfer struct {
	TokenAddress common.Address
	ToAddress    common.Address
	Amount       *big.Int
}

// CheckPoWBlockPoW is an auto generated low-level Go binding around an user-defined struct.
type CheckPoWBlockPoW struct {
	P1                    []byte
	PrevHashOrReceiptRoot [32]byte
	P2                    []byte
	Difficulty            []byte
	P3                    []byte
}

// CheckPoABlockPoA is an auto generated low-level Go binding around an user-defined struct.
type CheckPoABlockPoA struct {
	P0Seal                []byte
	P0Bare                []byte
	P1                    []byte
	PrevHashOrReceiptRoot [32]byte
	P2                    []byte
	Timestamp             []byte
	P3                    []byte
	S1                    []byte
	Signature             []byte
	S2                    []byte
}

type ReceiptsProof [][]byte
