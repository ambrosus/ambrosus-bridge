package merkle

import "container/list"

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

type Node struct {
	Data      NodeData
	NodeCount uint32
	Branches  *map[uint32]BranchTree
}

type ElementData interface{}

type NodeData interface{ Copy() NodeData }

type hashFunc func(NodeData, NodeData) NodeData

type elementHashFunc func(ElementData) NodeData

type MerkleTree struct {
	hash            hashFunc
	merkleBuf       *list.List
	elementHash     elementHashFunc
	exportNodeCount uint32
	indexes         map[uint32]bool
	exportNodes     []NodeData
}

func (m *MerkleTree) Insert(data, index uint32) {
	node := Node{
		Data:      m.elementHash(data),
		NodeCount: 1,
		Branches:  &map[uint32]BranchTree{},
	}

	if m.indexes[index] {
		(*node.Branches)[index] = BranchTree{
			RawData:    data,
			HashedData: node.Data,
			Root: &BranchNode{
				Hash:  node.Data,
				Left:  nil,
				Right: nil,
			},
		}
	}

	m.insertNode(node)
}

func (m *MerkleTree) insertNode(node Node) {
	var (
		element, prev   *list.Element
		cNode, prevNode Node
	)

	element = m.merkleBuf.PushBack(node)

	for {
		prev = element.Prev()
		cNode = element.Value.(Node)
		if prev == nil {
			break
		}

		prevNode = prev.Value.(Node)
		if cNode.NodeCount != prevNode.NodeCount {
			break
		}

		if prevNode.Branches != nil {
			for k, v := range *prevNode.Branches {
				v.Root = AcceptRightSibling(v.Root, cNode.Data)
				(*prevNode.Branches)[k] = v
			}
		}

		if cNode.Branches != nil {
			for k, v := range *cNode.Branches {
				v.Root = AcceptLeftSibling(v.Root, prevNode.Data)
				(*prevNode.Branches)[k] = v
			}
		}

		prevNode.Data = m.hash(prevNode.Data, cNode.Data)

		prevNode.NodeCount = cNode.NodeCount*2 + 1

		if prevNode.NodeCount == m.exportNodeCount {
			m.exportNodes = append(m.exportNodes, prevNode.Data)
		}

		m.merkleBuf.Remove(element)
		m.merkleBuf.Remove(prev)

		element = m.merkleBuf.PushBack(prevNode)
	}
}
