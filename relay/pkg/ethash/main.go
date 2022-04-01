package ethash

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash/merkle"
	"github.com/ethereum/go-ethereum/log"
)

type Ethash struct {
	dir    string
	logger log.Logger
	// todo caches
	// todo logger

}

func New(dir string) *Ethash {
	logger := log.New("epoch", epoch) // todo

	return &Ethash{
		dir:    dir,
		logger: logger,
	}
}

type EpochData struct {
	Epoch                   *big.Int
	FullSizeIn128Resolution *big.Int
	BranchDepth             *big.Int
	MerkleNodes             []*big.Int
}

func (e *Ethash) GenerateEpochData(epoch uint64) (*EpochData, error) {
	// todo use cache
	dag, err := e.getOrGenerateDag(epoch)
	if err != nil {
		return nil, err
	}

	fullSize := len(dag)
	fullSizeIn128Resolution := fullSize / 128
	branchDepth := len(fmt.Sprintf("%b", fullSizeIn128Resolution-1))

	mt := merkle.NewDatasetTree()
	mt.RegisterStoredLevel(uint32(branchDepth), uint32(10))
	dagToMerkle(dag, mt)
	mt.Finalize()

	return &EpochData{
		Epoch:                   big.NewInt(int64(epoch)),
		FullSizeIn128Resolution: big.NewInt(int64(fullSizeIn128Resolution)),
		BranchDepth:             big.NewInt(int64(branchDepth - 10)),
		MerkleNodes:             mt.MerkleNodes(),
	}, nil
}

func (e *Ethash) GetBlockLookups(blockNumber uint64, nonce uint64, hashNoNonce [32]byte) (dataSetLookup, witnessForLookup []*big.Int, err error) {
	// todo use cache
	dag, err := e.getOrGenerateDag(epoch(blockNumber))
	if err != nil {
		return
	}

	fullSize := len(dag)
	fullSizeIn128Resolution := fullSize / 128
	branchDepth := len(fmt.Sprintf("%b", fullSizeIn128Resolution-1))

	mt := merkle.NewDatasetTree()

	indices, err := e.getVerificationIndices(blockNumber, hashNoNonce, nonce)
	if err != nil {
		return
	}
	mt.RegisterIndex(indices...) // diff with GenerateEpochData :(

	mt.RegisterStoredLevel(uint32(branchDepth), 10)
	dagToMerkle(dag, mt)
	mt.Finalize()

	for _, w := range mt.AllDAGElements() {
		dataSetLookup = append(dataSetLookup, w.ToUint256Array()...)
	}

	for _, be := range mt.AllBranchesArray() {
		witnessForLookup = append(witnessForLookup, be.Big())
	}

	return
}

func (e *Ethash) deleteOldData(epoch uint64) {
	// Iterate over all previous instances and delete old ones
	for ep := epoch; ep >= 0; ep-- {
		_ = os.Remove(e.pathToDag(ep))
		_ = os.Remove(e.pathToCache(ep))
	}
}

func (e *Ethash) pathToCache(epoch uint64) string {
	return filepath.Join(e.dir, fmt.Sprintf("cache-%d", epoch))
}
func (e *Ethash) pathToDag(epoch uint64) string {
	return filepath.Join(e.dir, fmt.Sprintf("dag-%d", epoch))
}
