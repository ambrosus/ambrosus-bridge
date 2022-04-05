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

func (d SPHash) Copy() NodeData {
	result := SPHash{}
	copy(result[:], d[:])
	return result
}

func NewDatasetTree() *DatasetTree {
	return &DatasetTree{
		merkleBuf:      list.New(),
		indexes:        map[uint32]bool{},
		orderedIndexes: []uint32{},
		exportNodes:    []NodeData{},
	}
}

func (d DatasetTree) MerkleNodes() []*big.Int {
	d.assertFinalized()
	result := []*big.Int{}

	for i := 0; i*2 < len(d.exportNodes); i++ {
		if i*2+1 >= len(d.exportNodes) {
			result = append(result, BranchElementFromHash(SPHash{}, d.exportNodes[i*2].(SPHash)).Big())
		} else {
			result = append(result, BranchElementFromHash(d.exportNodes[i*2+1].(SPHash), d.exportNodes[i*2].(SPHash)).Big())
		}
	}

	return result
}

func (d DatasetTree) DatasetLookups() []*big.Int {
	d.assertFinalized()
	var result []*big.Int
	branches := d.Branches()

	for _, k := range d.orderedIndexes {
		result = append(result, branches[k].RawData.(Word).ToUint256Array()...)
	}

	return result
}

func (d DatasetTree) WitnessForLookups() []*big.Int {
	d.assertFinalized()

	var result []*big.Int
	branches := d.Branches()

	for _, k := range d.orderedIndexes {
		hh := branches[k].ToNodeArray()[1:]
		hashes := hh[:len(hh)-int(d.storedLevel)]

		for i := 0; i*2 < len(hashes); i++ {
			// for anyone who is courious why i*2 + 1 comes before i * 2
			// it's agreement between client side and contract side
			if i*2+1 >= len(hashes) {
				result = append(result, BranchElementFromHash(SPHash{}, hashes[i*2].(SPHash)).Big())
			} else {
				result = append(result, BranchElementFromHash(hashes[i*2+1].(SPHash), hashes[i*2].(SPHash)).Big())
			}
		}
	}

	return result
}

func (d *DatasetTree) hash(a, b NodeData) NodeData {
	var keccak []byte

	left := a.(SPHash)
	right := b.(SPHash)

	keccak = crypto.Keccak256(
		append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, left[:]...),
		append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, right[:]...),
	)

	result := SPHash{}
	copy(result[:HashLength], keccak[HashLength:])
	return result
}

func (d *DatasetTree) elementHash(data ElementData) NodeData {
	first, second := conventionalWord(data.(Word))
	keccak := crypto.Keccak256(first, second)

	result := SPHash{}
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
