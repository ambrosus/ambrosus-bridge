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

type BranchTree struct {
	RawData    ElementData
	HashedData NodeData
	Root       *BranchNode
}

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
