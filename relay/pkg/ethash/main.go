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

	// cache (in default meaning)
	caches   map[uint64][]byte
	dags     map[uint64][]byte
	dagKLock *Kmutex // dagKLock need to make sure that for each epoch only one dag is generating at a time
	// no lock for cache because is doesn't take that long to generate it
}

func New(dir string) *Ethash {
	logger := log.New("epoch", epoch) // todo

	return &Ethash{
		dir:    dir,
		logger: logger,

		caches:   map[uint64][]byte{},
		dags:     map[uint64][]byte{},
		dagKLock: NewKmutex(),
	}
}

type EpochData struct {
	Epoch                   *big.Int
	FullSizeIn128Resolution *big.Int
	BranchDepth             *big.Int
	MerkleNodes             []*big.Int
}

func (e *Ethash) GetEpochData(epoch uint64) (*EpochData, error) {
	mt := merkle.NewDatasetTree()
	fullSize, branchDepth, err := e.populateMerkle(epoch, mt)
	if err != nil {
		return nil, err
	}

	return &EpochData{
		Epoch:                   big.NewInt(int64(epoch)),
		FullSizeIn128Resolution: big.NewInt(int64(fullSize / 128)),
		BranchDepth:             big.NewInt(int64(branchDepth - 10)),
		MerkleNodes:             mt.MerkleNodes(),
	}, nil
}

func (e *Ethash) GetBlockLookups(blockNumber uint64, nonce uint64, hashNoNonce [32]byte) (dataSetLookup, witnessForLookup []*big.Int, err error) {
	indices, err := e.getVerificationIndices(blockNumber, hashNoNonce, nonce)
	if err != nil {
		return
	}

	mt := merkle.NewDatasetTree()
	mt.RegisterIndex(indices...)
	_, _, err = e.populateMerkle(epoch(blockNumber), mt)
	if err != nil {
		return nil, nil, err
	}

	for _, w := range mt.AllDAGElements() {
		dataSetLookup = append(dataSetLookup, w.ToUint256Array()...)
	}

	for _, be := range mt.AllBranchesArray() {
		witnessForLookup = append(witnessForLookup, be.Big())
	}

	return
}

func (e *Ethash) deleteOldData(epoch uint64) {
	// create new if need

	// Iterate over all previous instances and delete old ones
	for ep := epoch; ep >= 0; ep-- {
		_ = os.Remove(e.pathToDag(ep))
		_ = os.Remove(e.pathToCache(ep))
		delete(e.dags, epoch)
		delete(e.caches, epoch)
	}
}

func (e *Ethash) pathToCache(epoch uint64) string {
	return filepath.Join(e.dir, fmt.Sprintf("cache-%d", epoch))
}
func (e *Ethash) pathToDag(epoch uint64) string {
	return filepath.Join(e.dir, fmt.Sprintf("dag-%d", epoch))
}
