package merkle

const BranchElementLength = 32

type BranchNode struct {
	Hash             SPHash
	Left             *BranchNode
	Right            *BranchNode
	ElementOnTheLeft bool
}

type BranchTree struct {
	RawData    Word
	HashedData SPHash
	Root       *BranchNode
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

func (t BranchTree) ToNodeArray() []SPHash {
	return t.Root.ToNodeArray()
}

func AcceptRightSibling(branch *BranchNode, data SPHash) *BranchNode {
	return &BranchNode{
		Left:             branch,
		Right:            &BranchNode{Hash: data},
		ElementOnTheLeft: true,
	}
}

func AcceptLeftSibling(branch *BranchNode, data SPHash) *BranchNode {
	return &BranchNode{
		Left:             &BranchNode{Hash: data},
		Right:            branch,
		ElementOnTheLeft: false,
	}
}
