// Copyright 2019 The go-ethereum Authors
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
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

type sliceBuffer []byte

func (b *sliceBuffer) Write(data []byte) (n int, err error) {
	*b = append(*b, data...)
	return len(data), nil
}

func (b *sliceBuffer) Reset() {
	*b = (*b)[:0]
}

// Hasher is a type used for the trie Hash operation. A Hasher has some
// internal preallocated temp space
type Hasher struct {
	sha      crypto.KeccakState
	tmp      sliceBuffer
	parallel bool // Whether to use paralallel threads when hashing
}

func NewHasher() *Hasher {
	return &Hasher{
		tmp: make(sliceBuffer, 0, 550), // cap is as large as a full fullNode.
		sha: sha3.NewLegacyKeccak256().(crypto.KeccakState),
	}
}

func (h *Hasher) HashTmp() []byte {
	buf := make([]byte, 32)
	h.sha.Write(h.tmp)
	h.sha.Read(buf)
	return buf
}

func (h *Hasher) Hash(v []byte) []byte {
	buf := make([]byte, 32)
	h.sha.Write(v)
	h.sha.Read(buf)
	return buf
}

func Hash(v []byte) []byte {
	sha := sha3.NewLegacyKeccak256().(crypto.KeccakState)
	buf := make([]byte, 32)
	sha.Write(v)
	sha.Read(buf)
	return buf
}
