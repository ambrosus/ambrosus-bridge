package merkle

import "container/list"

type Node struct {
	Data      SPHash
	NodeCount uint32
	Branches  *map[uint32]BranchTree
}

func (n Node) Copy() Node {
	return Node{n.Data.Copy(), n.NodeCount, &map[uint32]BranchTree{}}
}

//type NodeData interface{ Copy() NodeData }

type DatasetTree struct {
	merkleBuf       *list.List
	exportNodeCount uint32
	storedLevel     uint32
	finalized       bool
	indexes         map[uint32]bool
	orderedIndexes  []uint32
	exportNodes     []SPHash
}

func (d *DatasetTree) Insert(data Word, index uint32) {
	node := Node{
		Data:      d.elementHash(data),
		NodeCount: 1,
		Branches:  &map[uint32]BranchTree{},
	}

	if d.indexes[index] {
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

	d.insertNode(node)
}

func (d *DatasetTree) insertNode(node Node) {
	var (
		element, prev   *list.Element
		cNode, prevNode Node
	)

	element = d.merkleBuf.PushBack(node)

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

		prevNode.Data = d.hash(prevNode.Data, cNode.Data)

		prevNode.NodeCount = cNode.NodeCount*2 + 1

		if prevNode.NodeCount == d.exportNodeCount {
			d.exportNodes = append(d.exportNodes, prevNode.Data)
		}

		d.merkleBuf.Remove(element)
		d.merkleBuf.Remove(prev)

		element = d.merkleBuf.PushBack(prevNode)
	}
}

func (d *DatasetTree) Finalize() {
	if d.finalized {
		return
	}

	if d.merkleBuf.Len() > 1 {
		for d.merkleBuf.Len() != 1 {
			dupNode := d.merkleBuf.Back().Value.(Node).Copy()
			d.insertNode(dupNode)
		}
	}

	d.finalized = true
}

func (d *DatasetTree) RegisterStoredLevel(depth, level uint32) {
	d.storedLevel = level
	d.exportNodeCount = 1<<(depth-level+1) - 1
}

func (d *DatasetTree) RegisterIndex(indexes ...uint32) {
	for _, i := range indexes {
		d.indexes[i] = true
		d.orderedIndexes = append(d.orderedIndexes, i)
	}
}

func (d DatasetTree) Branches() map[uint32]BranchTree {
	d.assertFinalized()
	return *(d.merkleBuf.Front().Value.(Node).Branches)
}

func (d DatasetTree) assertFinalized() {
	if !d.finalized {
		panic("SP Merkle tree needs to be finalized by calling mt.Finalize()")
	}
}
