package merkle

import (
	"container/list"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	HashLength = 16
	WordLength = 128
)

type SPHash [HashLength]byte

type Word [WordLength]byte

type DatasetData SPHash

func (d DatasetData) Copy() NodeData {
	result := DatasetData{}
	copy(result[:], d[:])

	return result
}

type DatasetTree struct{ MerkleTree }

func NewDatasetTree() *DatasetTree {
	merkleBuf := list.New()

	return &DatasetTree{
		MerkleTree{
			hash:            hash,
			merkleBuf:       merkleBuf,
			elementHash:     elementHash,
			exportNodeCount: 0,
			storedLevel:     0,
			finalized:       false,
			indexes:         map[uint32]bool{},
			exportNodes:     []NodeData{},
		},
	}
}

func hash(a, b NodeData) NodeData {
	var keccak []byte

	left := a.(DatasetData)
	right := b.(DatasetData)

	keccak = crypto.Keccak256(
		append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, left[:]...),
		append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, right[:]...),
	)

	result := DatasetData{}

	copy(result[:HashLength], keccak[HashLength:])

	return result
}

func elementHash(data ElementData) NodeData {
	first, second := conventionalWord(data.(Word))
	keccak := crypto.Keccak256(first, second)

	result := DatasetData{}

	copy(result[:HashLength], keccak[HashLength:])

	return result
}

func conventionalWord(data Word) ([]byte, []byte) {
	first := rev(data[:32])
	first = append(first, rev(data[32:64])...)

	second := rev(data[64:96])
	second = append(second, rev(data[96:128])...)

	return first, second
}

func rev(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return b
}
