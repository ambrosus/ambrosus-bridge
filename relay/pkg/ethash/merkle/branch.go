package merkle

import "math/big"

const BranchElementLength = 32

type BranchElement [BranchElementLength]byte

func (b BranchElement) Big() *big.Int { return bytesToBig(b[:]) }

type BranchNode struct {
	Hash             SPHash
	Left             *BranchNode
	Right            *BranchNode
	ElementOnTheLeft bool
}

func (b BranchNode) ToNodeArray() []SPHash {
	if b.Left == nil && b.Right == nil {
		return []SPHash{b.Hash}
	}
	left := b.Left.ToNodeArray()
	right := b.Right.ToNodeArray()
	if b.ElementOnTheLeft {
		return append(left, right...)
	} else {
		return append(right, left...)
	}
}

type BranchTree struct {
	RawData    ElementData
	HashedData SPHash
	Root       *BranchNode
}

func (t BranchTree) ToNodeArray() []SPHash { return t.Root.ToNodeArray() }

func AcceptRightSibling(branch *BranchNode, data SPHash) *BranchNode {
	return &BranchNode{
		Left:             branch,
		Right:            &BranchNode{data, nil, nil, false},
		ElementOnTheLeft: true,
	}
}

func AcceptLeftSibling(branch *BranchNode, data SPHash) *BranchNode {
	return &BranchNode{
		Left:             &BranchNode{data, nil, nil, false},
		Right:            branch,
		ElementOnTheLeft: false,
	}
}

func BranchElementFromHash(a, b SPHash) BranchElement {
	result := BranchElement{}
	copy(result[:], append(a[:], b[:]...)[:BranchElementLength])
	return result
}

func bytesToBig(data []byte) *big.Int {
	n := new(big.Int)
	n.SetBytes(data)
	return n
}
