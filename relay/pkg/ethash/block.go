package ethash

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash/merkle"

	"github.com/ethereum/go-ethereum/common"
)

type BlockMetaData struct {
	blockNumber uint64
	nonce       uint64
	hashNoNonce common.Hash
	DatasetTree *merkle.DatasetTree
}

func NewBlockMetaData(blockNumber uint64, nonce uint64, rlpHeaderHashWithoutNonce [32]byte) *BlockMetaData {
	return &BlockMetaData{
		blockNumber: blockNumber,
		nonce:       nonce,
		hashNoNonce: rlpHeaderHashWithoutNonce,
		DatasetTree: nil,
	}
}

func (s *BlockMetaData) DAGElementArray() []*big.Int {
	if s.DatasetTree == nil {
		s.buildDagTree()
	}

	result := []*big.Int{}
	for _, w := range s.DatasetTree.AllDAGElements() {
		result = append(result, w.ToUint256Array()...)
	}

	return result
}

func (s *BlockMetaData) buildDagTree() {
	indices := Instance.GetVerificationIndices(
		s.blockNumber,
		s.hashNoNonce,
		s.nonce,
	)

	s.DatasetTree = merkle.NewDatasetTree()
	s.DatasetTree.RegisterIndex(indices...)

	MakeDAG(s.blockNumber, DefaultDir)

	fullSize := DatasetSize(s.blockNumber)
	fullSizeIn128Resolution := fullSize / 128
	branchDepth := len(fmt.Sprintf("%b", fullSizeIn128Resolution-1))

	s.DatasetTree.RegisterStoredLevel(uint32(branchDepth), uint32(10))

	path := PathToDAG(s.blockNumber/30000, DefaultDir)

	ProcessDuringRead(path, s.DatasetTree)

	s.DatasetTree.Finalize()
}

func (s *BlockMetaData) DAGProofArray() []*big.Int {
	if s.DatasetTree == nil {
		s.buildDagTree()
	}

	result := []*big.Int{}

	for _, be := range s.DatasetTree.AllBranchesArray() {
		result = append(result, be.Big())
	}

	return result
}

func ProcessDuringRead(datasetPath string, mt *merkle.DatasetTree) error {
	var (
		file *os.File
		err  error
	)

	file, err = os.Open(datasetPath)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	buf := [128]byte{}

	_, err = io.ReadFull(reader, buf[:8])
	if err != nil {
		return err
	}

	var i uint32 = 0

	for {
		n, err := io.ReadFull(reader, buf[:128])
		if n == 0 {
			if err == nil {
				continue
			} else if err == io.EOF {
				break
			}

			return err
		} else if n != 128 {
			return errors.New("error malformed dataset")
		}

		mt.Insert(merkle.Word(buf), i)

		if err != nil && err != io.EOF {
			return err
		}

		i++
	}

	return nil
}
