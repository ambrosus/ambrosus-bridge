// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package mytrie

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	// emptyRoot is the known root hash of an empty trie.
	emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

// ModifiedStackTrie is a trie implementation that expects keys to be inserted
// in order. Once it determines that a subtree will no longer be inserted
// into, it will hash it and free up the memory it uses.
type ModifiedStackTrie struct {
	NodeType    uint8  // node type (as in branch, ext, leaf)
	Val         []byte // value contained by this node if it's a leaf
	UnhashedVal []byte
	Key         []byte                 // Key chunk covered by this (full|ext) node
	keyOffset   int                    // offset of the Key chunk inside a full Key
	Children    [16]*ModifiedStackTrie // list of Children (for fullnodes and exts)
}

// NewStackTrie allocates and initializes an empty trie.
func NewStackTrie() *ModifiedStackTrie {
	return &ModifiedStackTrie{
		NodeType: emptyNode,
	}
}

// MarshalBinary implements encoding.BinaryMarshaler
func (st *ModifiedStackTrie) MarshalBinary() (data []byte, err error) {
	var (
		b bytes.Buffer
		w = bufio.NewWriter(&b)
	)
	if err := gob.NewEncoder(w).Encode(struct {
		Nodetype  uint8
		Val       []byte
		Key       []byte
		KeyOffset uint8
	}{
		st.NodeType,
		st.Val,
		st.Key,
		uint8(st.keyOffset),
	}); err != nil {
		return nil, err
	}
	for _, child := range st.Children {
		if child == nil {
			w.WriteByte(0)
			continue
		}
		w.WriteByte(1)
		if childData, err := child.MarshalBinary(); err != nil {
			return nil, err
		} else {
			w.Write(childData)
		}
	}
	w.Flush()
	return b.Bytes(), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (st *ModifiedStackTrie) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	return st.unmarshalBinary(r)
}

func (st *ModifiedStackTrie) unmarshalBinary(r io.Reader) error {
	var dec struct {
		Nodetype  uint8
		Val       []byte
		Key       []byte
		KeyOffset uint8
	}
	gob.NewDecoder(r).Decode(&dec)
	st.NodeType = dec.Nodetype
	st.Val = dec.Val
	st.Key = dec.Key
	st.keyOffset = int(dec.KeyOffset)

	var hasChild = make([]byte, 1)
	for i := range st.Children {
		if _, err := r.Read(hasChild); err != nil {
			return err
		} else if hasChild[0] == 0 {
			continue
		}
		var child ModifiedStackTrie
		child.unmarshalBinary(r)
		st.Children[i] = &child
	}
	return nil
}

func newLeaf(ko int, key, val []byte) *ModifiedStackTrie {
	st := NewStackTrie()
	st.NodeType = leafNode
	st.keyOffset = ko
	st.Key = append(st.Key, key[ko:]...)
	st.Val = val
	return st
}

func newExt(ko int, key []byte, child *ModifiedStackTrie) *ModifiedStackTrie {
	st := NewStackTrie()
	st.NodeType = extNode
	st.keyOffset = ko
	st.Key = append(st.Key, key[ko:]...)
	st.Children[0] = child
	return st
}

// List all values that ModifiedStackTrie#NodeType can hold
const (
	emptyNode = iota
	branchNode
	extNode
	leafNode
	hashedNode
)

// TryUpdate inserts a (Key, value) pair into the stack trie
func (st *ModifiedStackTrie) TryUpdate(key, value []byte) error {
	k := keybytesToHex(key)
	if len(value) == 0 {
		panic("deletion not supported")
	}
	st.insert(k[:len(k)-1], value)
	return nil
}

func (st *ModifiedStackTrie) Update(key, value []byte) {
	if err := st.TryUpdate(key, value); err != nil {
		log.Error(fmt.Sprintf("Unhandled trie error: %v", err))
	}
}

func (st *ModifiedStackTrie) Reset() {
	st.Key = st.Key[:0]
	st.Val = nil
	for i := range st.Children {
		st.Children[i] = nil
	}
	st.NodeType = emptyNode
	st.keyOffset = 0
}

// Helper function that, given a full Key, determines the index
// at which the chunk pointed by st.keyOffset is different from
// the same chunk in the full Key.
func (st *ModifiedStackTrie) getDiffIndex(key []byte) int {
	diffindex := 0
	for ; diffindex < len(st.Key) && st.Key[diffindex] == key[st.keyOffset+diffindex]; diffindex++ {
	}
	return diffindex
}

// Helper function to that inserts a (Key, value) pair into
// the trie.
func (st *ModifiedStackTrie) insert(key, value []byte) {
	switch st.NodeType {
	case branchNode: /* Branch */
		idx := int(key[st.keyOffset])
		// Unresolve elder siblings
		for i := idx - 1; i >= 0; i-- {
			if st.Children[i] != nil {
				if st.Children[i].NodeType != hashedNode {
					st.Children[i].hash()
				}
				break
			}
		}
		// Add new child
		if st.Children[idx] == nil {
			st.Children[idx] = NewStackTrie()
			st.Children[idx].keyOffset = st.keyOffset + 1
		}
		st.Children[idx].insert(key, value)
	case extNode: /* Side */
		// Compare both Key chunks and see where they differ
		diffidx := st.getDiffIndex(key)

		// Check if chunks are identical. If so, recurse into
		// the child node. Otherwise, the Key has to be split
		// into 1) an optional common prefix, 2) the fullnode
		// representing the two differing path, and 3) a leaf
		// for each of the differentiated subtrees.
		if diffidx == len(st.Key) {
			// Side Key and Key segment are identical, recurse into
			// the child node.
			st.Children[0].insert(key, value)
			return
		}
		// Save the original part. Depending if the break is
		// at the extension's last byte or not, create an
		// intermediate extension or use the extension's child
		// node directly.
		var n *ModifiedStackTrie
		if diffidx < len(st.Key)-1 {
			n = newExt(diffidx+1, st.Key, st.Children[0])
		} else {
			// Break on the last byte, no need to insert
			// an extension node: reuse the current node
			n = st.Children[0]
		}
		// Convert to hash
		n.hash()
		var p *ModifiedStackTrie
		if diffidx == 0 {
			// the break is on the first byte, so
			// the current node is converted into
			// a branch node.
			st.Children[0] = nil
			p = st
			st.NodeType = branchNode
		} else {
			// the common prefix is at least one byte
			// long, insert a new intermediate branch
			// node.
			st.Children[0] = NewStackTrie()
			st.Children[0].NodeType = branchNode
			st.Children[0].keyOffset = st.keyOffset + diffidx
			p = st.Children[0]
		}
		// Create a leaf for the inserted part
		o := newLeaf(st.keyOffset+diffidx+1, key, value)

		// Insert both child leaves where they belong:
		origIdx := st.Key[diffidx]
		newIdx := key[diffidx+st.keyOffset]
		p.Children[origIdx] = n
		p.Children[newIdx] = o
		st.Key = st.Key[:diffidx]

	case leafNode: /* Leaf */
		// Compare both Key chunks and see where they differ
		diffidx := st.getDiffIndex(key)

		// Overwriting a Key isn't supported, which means that
		// the current leaf is expected to be split into 1) an
		// optional extension for the common prefix of these 2
		// keys, 2) a fullnode selecting the path on which the
		// keys differ, and 3) one leaf for the differentiated
		// component of each Key.
		if diffidx >= len(st.Key) {
			panic("Trying to insert into existing Key")
		}

		// Check if the split occurs at the first nibble of the
		// chunk. In that case, no prefix extnode is necessary.
		// Otherwise, create that
		var p *ModifiedStackTrie
		if diffidx == 0 {
			// Convert current leaf into a branch
			st.NodeType = branchNode
			p = st
			st.Children[0] = nil
		} else {
			// Convert current node into an ext,
			// and insert a child branch node.
			st.NodeType = extNode
			st.Children[0] = NewStackTrie()
			st.Children[0].NodeType = branchNode
			st.Children[0].keyOffset = st.keyOffset + diffidx
			p = st.Children[0]
		}

		// Create the two child leaves: the one containing the
		// original value and the one containing the new value
		// The child leave will be hashed directly in order to
		// free up some memory.
		origIdx := st.Key[diffidx]
		p.Children[origIdx] = newLeaf(diffidx+1, st.Key, st.Val)
		p.Children[origIdx].hash()

		newIdx := key[diffidx+st.keyOffset]
		p.Children[newIdx] = newLeaf(p.keyOffset+1, key, value)

		// Finally, cut off the Key part that has been passed
		// over to the Children.
		st.Key = st.Key[:diffidx]
		st.Val = nil
	case emptyNode: /* Empty */
		st.NodeType = leafNode
		st.Key = key[st.keyOffset:]
		st.Val = value
	case hashedNode:
		panic("trying to insert into hash")
	default:
		panic("invalid type")
	}
}

type extNodeStruct struct {
	Key []byte
	Val []byte
}

// hash() hashes the node 'st' and converts it into 'hashedNode', if possible.
// Possible outcomes:
// 1. The rlp-encoded value was >= 32 bytes:
//  - Then the 32-byte `hash` will be accessible in `st.Val`.
//  - And the 'st.type' will be 'hashedNode'
// 2. The rlp-encoded value was < 32 bytes
//  - Then the <32 byte rlp-encoded value will be accessible in 'st.Val'.
//  - And the 'st.type' will be 'hashedNode' AGAIN
//
// This method will also:
// set 'st.type' to hashedNode
// clear 'st.Key'
func (st *ModifiedStackTrie) hash() {
	/* Shortcut if node is already hashed */
	if st.NodeType == hashedNode {
		return
	}

	buf := new(bytes.Buffer)

	switch st.NodeType {
	case emptyNode:
		st.Val = emptyRoot.Bytes()
		st.Key = st.Key[:0]
		st.NodeType = hashedNode
		return

	case branchNode:
		var nodes [17][]byte
		for i, child := range st.Children {
			if child != nil {
				child.hash()
				nodes[i] = child.Val
			}
		}
		nodes[16] = nil
		if err := rlp.Encode(buf, nodes); err != nil {
			panic(err)
		}

	case extNode:
		st.Children[0].hash()
		n := extNodeStruct{hexToCompact(st.Key), st.Children[0].Val}

		if err := rlp.Encode(buf, n); err != nil {
			panic(err)
		}

	case leafNode:

		st.Key = append(st.Key, byte(16))
		sz := hexToCompactInPlace(st.Key)
		n := [][]byte{st.Key[:sz], st.Val}

		if err := rlp.Encode(buf, n); err != nil {
			panic(err)
		}

	default:
		panic("Invalid node type")
	}

	st.Key = st.Key[:0]
	st.NodeType = hashedNode

	// new
	st.UnhashedVal = buf.Bytes()

	if buf.Len() > 32 {
		st.Val = crypto.Keccak256(buf.Bytes())
		return
	}
	st.Val = buf.Bytes()

}

// Hash returns the hash of the current node
func (st *ModifiedStackTrie) Hash() (h common.Hash) {
	st.hash()
	if len(st.Val) != 32 {
		// If the node's RLP isn't 32 bytes long, the node will not
		// be hashed, and instead contain the  rlp-encoding of the
		// node. For the top level node, we need to force the hashing.
		return common.BytesToHash(crypto.Keccak256(st.Val))
	}
	return common.BytesToHash(st.Val)
}
