package merkle

import (
	"container/list"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	HashLength = 16
	WordLength = 128
)

type SPHash [HashLength]byte
type Word [WordLength]byte

func (d SPHash) Copy() SPHash {
	result := SPHash{}
	copy(result[:], d[:])
	return result
}

func NewDatasetTree() *DatasetTree {
	return &DatasetTree{
		merkleBuf:      list.New(),
		indexes:        map[uint32]bool{},
		orderedIndexes: []uint32{},
		exportNodes:    []SPHash{},
	}
}

func (d DatasetTree) MerkleNodes() []*big.Int {
	d.assertFinalized()
	var result []*big.Int

	for i := 0; i*2 < len(d.exportNodes); i++ {
		f := SPHash{}
		if i*2+1 < len(d.exportNodes) {
			f = d.exportNodes[i*2+1]
		}

		result = append(result, BranchElementFromHash(f, d.exportNodes[i*2]).Big())
	}

	return result
}

func (d DatasetTree) Lookups() (dataSetLookup, witnessForLookup []*big.Int) {
	d.assertFinalized()

	branches := d.Branches()

	for _, k := range d.orderedIndexes {

		hh := branches[k].ToNodeArray()[1:]
		hashes := hh[:len(hh)-int(d.storedLevel)]

		for i := 0; i*2 < len(hashes); i++ {
			// for anyone who is courious why i*2 + 1 comes before i * 2
			// it's agreement between client side and contract side
			f := SPHash{}
			if i*2+1 < len(hashes) {
				f = hashes[i*2+1]
			}

			witnessForLookup = append(witnessForLookup, BranchElementFromHash(f, hashes[i*2]).Big())
		}

		dataSetLookup = append(dataSetLookup, branches[k].RawData.ToUint256Array()...)
	}

	return
}

func (d *DatasetTree) hash(left, right SPHash) SPHash {
	keccak := crypto.Keccak256(
		append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, left[:]...),
		append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, right[:]...),
	)

	result := SPHash{}
	copy(result[:HashLength], keccak[HashLength:])
	return result
}

func (d *DatasetTree) elementHash(data Word) SPHash {
	first, second := conventionalWord(data)
	keccak := crypto.Keccak256(first, second)

	result := SPHash{}
	copy(result[:HashLength], keccak[HashLength:])
	return result
}

func conventionalWord(data Word) (first, second []byte) {
	rev(data[0:32])
	rev(data[32:64])
	rev(data[64:96])
	rev(data[96:128])
	return data[:64], data[64:]
}

func (w Word) ToUint256Array() []*big.Int {
	result := []*big.Int{}

	for i := 0; i < WordLength/32; i++ {
		z := big.NewInt(0)
		z.SetBytes(rev(w[i*32 : (i+1)*32]))

		result = append(result, z)
	}

	return result
}

func rev(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return b
}
