package ethash

import (
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash/merkle"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	epochLength = 30000
	maxEpoch    = 2048
)

type EpochData struct {
	Epoch                   *hexutil.Big
	FullSizeIn128Resolution *hexutil.Big
	BranchDepth             *hexutil.Big
	MerkleNodes             []*hexutil.Big
}

func GenerateEpochData(epoch uint64) (*EpochData, error) {
	fullSize := DatasetSize(epoch * epochLength)
	fullSizeIn128Resolution := fullSize / 128
	branchDepth := len(fmt.Sprintf("%b", fullSizeIn128Resolution-1))

	path, err := CheckDatasetFile(epoch, DefaultDir)
	if err != nil {
		return &EpochData{}, nil
	}

	mt := merkle.NewDatasetTree()
	mt.RegisterStoredLevel(uint32(branchDepth), 10)

	if err := ProcessDuringRead(path, mt); err != nil {
		return &EpochData{}, err
	}

	mt.Finalize()

	return &EpochData{
		Epoch:                   (*hexutil.Big)(big.NewInt(int64(epoch))),
		FullSizeIn128Resolution: (*hexutil.Big)(big.NewInt(int64(fullSizeIn128Resolution))),
		BranchDepth:             (*hexutil.Big)(big.NewInt(int64(branchDepth - 10))),
		MerkleNodes:             mt.MerkleNodes(),
	}, nil
}
