package merkle

import (
	"container/list"

	"github.com/ethereum/go-ethereum/crypto"
)

type DatasetTree struct {
	merkleBuf       *list.List
	exportNodeCount uint32
	storedLevel     uint32
	finalized       bool
	indexes         map[uint32]bool
	orderedIndexes  []uint32
	exportNodes     []SPHash
}

func NewDatasetTree() *DatasetTree {
	return &DatasetTree{
		merkleBuf:      list.New(),
		indexes:        map[uint32]bool{},
		orderedIndexes: []uint32{},
		exportNodes:    []SPHash{},
	}
}

type Node struct {
	Data      SPHash
	NodeCount uint32
	Branches  *map[uint32]BranchTree
}

func (n Node) Copy() Node {
	return Node{n.Data.Copy(), n.NodeCount, &map[uint32]BranchTree{}}
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
			Root:       &BranchNode{Hash: node.Data},
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
		if prev == nil {
			break
		}

		cNode = element.Value.(Node)
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

	for d.merkleBuf.Len() > 1 {
		dupNode := d.merkleBuf.Back().Value.(Node).Copy()
		d.insertNode(dupNode)
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

func (d *DatasetTree) hash(left, right SPHash) SPHash {
	return d.wtfHash(
		append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, left[:]...),
		append([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, right[:]...),
	)
}

func (d *DatasetTree) elementHash(data Word) SPHash {
	rev(data[0:32])
	rev(data[32:64])
	rev(data[64:96])
	rev(data[96:128])
	first, second := data[:64], data[64:]
	return d.wtfHash(first, second)
}

func (d *DatasetTree) wtfHash(first, second []byte) SPHash {
	keccak := crypto.Keccak256(first, second)
	result := SPHash{}
	copy(result[:HashLength], keccak[HashLength:])
	return result
}
