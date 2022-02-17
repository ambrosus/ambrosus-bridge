package merkle

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
