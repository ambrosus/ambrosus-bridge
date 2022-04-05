package merkle

import (
	"math/big"
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

func (d DatasetTree) MerkleNodes() []*big.Int {
	d.assertFinalized()
	return branchesFromHashes(d.exportNodes)
}

func (d DatasetTree) Lookups() (dataSetLookup, witnessForLookup []*big.Int) {
	d.assertFinalized()

	branches := *(d.merkleBuf.Front().Value.(Node).Branches)

	for _, k := range d.orderedIndexes {
		hh := branches[k].ToNodeArray()[1:]
		hashes := hh[:len(hh)-int(d.storedLevel)]
		witnessForLookup = append(witnessForLookup, branchesFromHashes(hashes)...)

		dataSetLookup = append(dataSetLookup, branches[k].RawData.toUint256Array()...)
	}

	return
}

func branchesFromHashes(hashes []SPHash) []*big.Int {
	var result []*big.Int
	for i := 0; i*2 < len(hashes); i++ {
		// todo what does this mean?
		// for anyone who is courious why i*2 + 1 comes before i * 2
		// it's agreement between client side and contract side
		a := SPHash{}
		if i*2+1 < len(hashes) {
			a = hashes[i*2+1]
		}
		b := hashes[i*2]

		result = append(result, branchElement(a, b))
	}
	return result
}

func (d DatasetTree) assertFinalized() {
	if !d.finalized {
		panic("SP Merkle tree needs to be finalized by calling mt.Finalize()")
	}
}

func (w Word) toUint256Array() []*big.Int {
	var result []*big.Int
	for i := 0; i < WordLength/32; i++ {
		b := rev(w[i*32 : (i+1)*32])
		bi := new(big.Int).SetBytes(b)
		result = append(result, bi)
	}
	return result
}

func rev(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func branchElement(a, b SPHash) *big.Int {
	bytes := append(a[:], b[:]...)[:BranchElementLength]
	return new(big.Int).SetBytes(bytes)
}
