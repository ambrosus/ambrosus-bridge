package rolling_finality

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// RollingFinality checker for authority round consensus.
// Stores a chain of unfinalized hashes that can be pushed onto.
type RollingFinality struct {
	blocks       []block
	validatorSet []common.Address
	signCount    map[common.Address]uint
}
type block struct {
	num    uint64
	signer common.Address
}

// NewRollingFinality creates a blank finality checker under the given validator set.
func NewRollingFinality(validatorSet []common.Address) *RollingFinality {
	return &RollingFinality{
		validatorSet: validatorSet,
		blocks:       []block{},
		signCount:    map[common.Address]uint{},
	}
}

// Push a hash onto the rolling finality checker (implying `subchain_head` == head.parent)
//
// Fails if `signer` isn't a member of the active validator set.
// Returns a list of all newly finalized headers.
func (f *RollingFinality) Push(num uint64, signer common.Address) (newlyFinalized []uint64, err error) {
	if !f.isValidator(signer) {
		return nil, fmt.Errorf("unknown validator")
	}

	f.push(num, signer)

	for f.isFinalized() {
		finalizedBlockNum := f.pop()
		newlyFinalized = append(newlyFinalized, finalizedBlockNum)
	}
	return newlyFinalized, nil
}

func (f *RollingFinality) push(blockNum uint64, signer common.Address) {
	// push to blocks
	f.blocks = append(f.blocks, block{blockNum, signer})
	// add signer
	f.signCount[signer] += 1
}

func (f *RollingFinality) pop() uint64 {
	// pop from blocks
	b := f.blocks[0]
	f.blocks = f.blocks[1:]

	// remove signer
	count, ok := f.signCount[b.signer]
	if !ok {
		panic("all hashes in `header` should have entries in `sign_count` for their signers")
	}
	f.signCount[b.signer] = count - 1
	if count <= 1 { // remove signer from map if he has 0 signs
		delete(f.signCount, b.signer)
	}

	return b.num
}

// isFinalized returns whether the first entry in `self.headers` is finalized.
func (f *RollingFinality) isFinalized() bool {
	return len(f.signCount)*2 > len(f.validatorSet)
}

func (f RollingFinality) isValidator(signer common.Address) bool {
	for _, v := range f.validatorSet {
		if v == signer {
			return true
		}
	}
	return false
}
