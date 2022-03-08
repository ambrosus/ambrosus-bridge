package receipts_proof

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ProofEvent interface {
	ProofElements() [][]byte
	Log() *types.Log
}

func CheckProofEvent(proof [][]byte, event ProofEvent) common.Hash {
	return CheckProof(proof, event.ProofElements())
}

func CalcProofEvent(receipts []*types.Receipt, event ProofEvent) ([][]byte, error) {
	return CalcProof(receipts, event.Log(), event.ProofElements())
}
