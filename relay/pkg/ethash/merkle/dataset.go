package merkle

import (
	"container/list"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	HashLength = 16
	WordLength = 128
)

type SPHash [HashLength]byte

type Word [WordLength]byte

func (w Word) ToUint256Array() []*big.Int {
	result := []*big.Int{}

	for i := 0; i < WordLength/32; i++ {
		z := big.NewInt(0)
		z.SetBytes(rev(w[i*32 : (i+1)*32]))

		result = append(result, z)
	}

	return result
}

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
			hash:             hash,
			merkleBuf:        merkleBuf,
			elementHash:      elementHash,
			dummyNodeModifie: dummyNodeModifie,
			exportNodeCount:  0,
			storedLevel:      0,
			finalized:        false,
			indexes:          map[uint32]bool{},
			orderedIndexes:   []uint32{},
			exportNodes:      []NodeData{},
		},
	}
}

func (d DatasetTree) MerkleNodes() []*hexutil.Big {
	if !d.finalized {
		panic("Merkle tree needs to be finalized")
	}
	var result []*hexutil.Big

	for i := 0; i*2 < len(d.exportNodes); i++ {
		if i*2+1 >= len(d.exportNodes) {
			result = append(result, (*hexutil.Big)(BranchElementFromHash(
				SPHash(DatasetData{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
				SPHash(d.exportNodes[i*2].(DatasetData)),
			).Big()))
		} else {
			result = append(result, (*hexutil.Big)(BranchElementFromHash(
				SPHash(d.exportNodes[i*2+1].(DatasetData)),
				SPHash(d.exportNodes[i*2].(DatasetData)),
			).Big()))
		}
	}

	return result
}

func (d DatasetTree) AllDAGElements() []Word {
	if d.finalized {
		result := []Word{}
		branches := d.Branches()

		for _, k := range d.Indices() {
			result = append(result, branches[k].RawData.(Word))
		}

		return result
	}

	panic("SP Merkle tree needs to be finalized by calling mt.Finalize()")
}

func (d DatasetTree) AllBranchesArray() []BranchElement {
	if d.finalized {
		result := []BranchElement{}
		branches := d.Branches()

		for _, k := range d.Indices() {
			hh := branches[k].ToNodeArray()[1:]
			hashes := hh[:len(hh)-int(d.StoredLevel())]

			for i := 0; i*2 < len(hashes); i++ {
				// for anyone who is courious why i*2 + 1 comes before i * 2
				// it's agreement between client side and contract side
				if i*2+1 >= len(hashes) {
					result = append(result, BranchElementFromHash(
						SPHash(DatasetData{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}),
						SPHash(hashes[i*2].(DatasetData))))
				} else {
					result = append(result, BranchElementFromHash(
						SPHash(hashes[i*2+1].(DatasetData)),
						SPHash(hashes[i*2].(DatasetData))))
				}
			}
		}

		return result
	}

	panic("SP Merkle tree needs to be finalized by calling mt.Finalize()")
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

func dummyNodeModifie(data NodeData) {}
