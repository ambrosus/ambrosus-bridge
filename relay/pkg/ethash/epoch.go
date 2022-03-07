package ethash

import (
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash/merkle"
)

const (
	epochLength = 30000
	maxEpoch    = 2048
)

type EpochData struct {
	Epoch                   *big.Int
	FullSizeIn128Resolution *big.Int
	BranchDepth             *big.Int
	MerkleNodes             []*big.Int
}

func GenerateEpochData(epoch uint64) (*EpochData, error) {
	fullSize := DatasetSize(epoch * epochLength)
	fullSizeIn128Resolution := fullSize / 128
	branchDepth := len(fmt.Sprintf("%b", fullSizeIn128Resolution-1))

	path := PathToDAG(epoch, DefaultDir)
	if err := checkDatasetPath(path, epoch); err != nil {
		return &EpochData{}, nil
	}

	mt := merkle.NewDatasetTree()
	mt.RegisterStoredLevel(uint32(branchDepth), 10)

	if err := ProcessDuringRead(path, mt); err != nil {
		return &EpochData{}, err
	}

	mt.Finalize()

	return &EpochData{
		Epoch:                   big.NewInt(int64(epoch)),
		FullSizeIn128Resolution: big.NewInt(int64(fullSizeIn128Resolution)),
		BranchDepth:             big.NewInt(int64(branchDepth - 10)),
		MerkleNodes:             mt.MerkleNodes(),
	}, nil
}

func checkDatasetPath(path string, epoch uint64) error {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			MakeDAG(epoch*epochLength, DefaultDir)

			return nil
		}

		return err
	}

	return nil
}
