package ethash

import (
	"fmt"
	"math/big"
)

const (
	epochLength uint64 = 30000
	maxEpoch    uint64 = 2048
)

type EpochData struct {
	Epoch                   *big.Int
	FullSizeIn128Resolution *big.Int
	BranchDepth             *big.Int
	MerkleNodes             []*big.Int
}

func GenerateEpochData(epoch uint64) *EpochData {
	fullSize := DatasetSize(epoch * epochLength)
	fullSizeIn128Resolution := fullSize / 128
	branchDepth := len(fmt.Sprintf("%b", fullSizeIn128Resolution-1))

	return &EpochData{
		Epoch:                   big.NewInt(int64(epoch)),
		FullSizeIn128Resolution: big.NewInt(int64(fullSizeIn128Resolution)),
		BranchDepth:             big.NewInt(int64(branchDepth - 10)),
	}
}
