package merkle

import "math/big"

const BranchElementLength = 32

type BranchElement [BranchElementLength]byte

func (b BranchElement) Big() *big.Int { return bytesToBig(b[:]) }

type BranchNode struct {
	Hash             NodeData
	Left             *BranchNode
	Right            *BranchNode
	ElementOnTheLeft bool
}

func (b BranchNode) ToNodeArray() []NodeData {
	if b.Left == nil && b.Right == nil {
		return []NodeData{b.Hash}
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
	HashedData NodeData
	Root       *BranchNode
}

func (t BranchTree) ToNodeArray() []NodeData { return t.Root.ToNodeArray() }

func AcceptRightSibling(branch *BranchNode, data NodeData) *BranchNode {
	return &BranchNode{
		Hash:             nil,
		Right:            &BranchNode{data, nil, nil, false},
		Left:             branch,
		ElementOnTheLeft: true,
	}
}

func AcceptLeftSibling(branch *BranchNode, data NodeData) *BranchNode {
	return &BranchNode{
		Hash:             nil,
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
