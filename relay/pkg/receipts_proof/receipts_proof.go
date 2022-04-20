package receipts_proof

import (
	"bytes"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof/mytrie"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"
)

func CheckProof(proof [][]byte, proofElements [][]byte) common.Hash {
	wrapProofsCount := len(proofElements) + 1
	// wrap proofElements with proof
	els := make([][]byte, 0, len(proofElements)+wrapProofsCount)
	for i, le := range proofElements {
		els = append(els, proof[i], le)
	}
	els = append(els, proof[wrapProofsCount-1])

	el := crypto.Keccak256(helpers.BytesConcat(els...))

	// compute trie root
	for i := wrapProofsCount; i < len(proof); i += 2 {
		el = helpers.BytesConcat(proof[i], el, proof[i+1])
		if len(el) > 32 {
			el = crypto.Keccak256(el)
		}
	}

	return common.BytesToHash(el)
}

func CalcProof(receipts []*types.Receipt, log *types.Log, proofElements [][]byte) ([][]byte, error) {
	rlpLog, err := rlp.EncodeToBytes(log)
	if err != nil {
		return nil, fmt.Errorf("rlp encode log: %w", err)
	}

	logResult, err := helpers.BytesSplit(rlpLog, proofElements)
	if err != nil {
		return nil, fmt.Errorf("split log: %w", err)
	}

	trieResult, err := encodeTrie(receipts, rlpLog)
	if err != nil {
		return nil, fmt.Errorf("encodetrie: %w", err)
	}

	/*
		these bytes are next to each other, so merge
		trieResult[0] + logResult[0]
		and
		logResult[-1] + trieResult[1]
	*/
	logResult[0] = append(trieResult[0], logResult[0]...)
	logResult[len(logResult)-1] = append(logResult[len(logResult)-1], trieResult[1]...)
	trieResult = trieResult[2:]

	return append(logResult, trieResult...), nil
}

type trieProof struct {
	whatSearch []byte
	path       [][]byte
}

func encodeTrie(receipts []*types.Receipt, rlpLog []byte) ([][]byte, error) {
	trie := mytrie.NewStackTrie()
	types.DeriveSha(types.Receipts(receipts), trie)

	p := trieProof{whatSearch: rlpLog, path: [][]byte{}}
	if !p.findTriePath(trie) {
		return nil, fmt.Errorf("rlpReceipt not found in trie")
	}
	return p.makeTrieProof()
}

// findTriePath try to find trie node with p.whatSearch in node.UnhashedVal
// when found - return to the root of the tree, saving node.UnhashedVal to path along the way
// return true if found
func (p *trieProof) findTriePath(node *mytrie.ModifiedStackTrie) bool {
	if node == nil {
		return false
	}
	if bytes.Contains(node.UnhashedVal, p.whatSearch) {
		p.path = append(p.path, node.UnhashedVal)
		return true
	}
	for _, child := range node.Children {
		if p.findTriePath(child) {
			p.path = append(p.path, node.UnhashedVal)
			return true
		}
	}
	return false
}

func (p *trieProof) makeTrieProof() ([][]byte, error) {
	result := make([][]byte, 0, len(p.path)*2)
	whatSearch := p.whatSearch

	for _, unhashedVal := range p.path {
		r := bytes.Split(unhashedVal, whatSearch)
		if len(r) != 2 {
			log.Debug().Msgf("unhashedVal = 0x%x, whatSearch = 0x%x", unhashedVal, whatSearch)
			return nil, fmt.Errorf("split result length (%v) != 2", len(r))
		}
		result = append(result, r[0], r[1])
		whatSearch = crypto.Keccak256(unhashedVal)
	}
	return result, nil
}
