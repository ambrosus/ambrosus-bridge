package receipts_proof

import (
	"bytes"
	"fmt"
	"relay/helpers"
	"relay/receipts_proof/mytrie"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func CheckProof(proof [][]byte, log *types.Log) common.Hash {
	el := helpers.BytesConcat(
		proof[0],
		log.Address.Bytes(),
		proof[1],
		log.Topics[1].Bytes(), // event id
		proof[2],
		log.Data, //transfers
		proof[3],
	)

	for i := 4; i < len(proof); i += 2 {
		el = helpers.BytesConcat(proof[i], el, proof[i+1])
		if len(el) > 32 {
			el = mytrie.Hash(el)
		}
	}

	return common.BytesToHash(el)
}

func CalcProof(receipts []*types.Receipt, log *types.Log) ([][]byte, error) {
	rlpLog, logResult, err := encodeLog(log)
	if err != nil {
		return nil, err
	}

	trieResult, err := encodeTrie(receipts, rlpLog)
	if err != nil {
		return nil, err
	}

	return append(logResult, trieResult...), nil
}

func encodeLog(log *types.Log) (rlpLog []byte, result [][]byte, err error) {
	if len(log.Topics) != 2 { // topic[0] always here, topic[1] must be event_id
		return nil, nil, fmt.Errorf("log.Topics length != 2 (only event_id must be indexed)")
	}

	rlpLog, err = rlp.EncodeToBytes(log)
	if err != nil {
		return
	}

	splitEls := [][]byte{
		log.Address.Bytes(),
		log.Topics[1].Bytes(),
		log.Data,
	}
	result, err = helpers.BytesSplit(rlpLog, splitEls)
	return
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
			return nil, fmt.Errorf("split result length (%v) != 2", len(r)) // todo pass split args here
		}
		result = append(result, r[0], r[1])
		whatSearch = mytrie.Hash(unhashedVal)
	}
	return result, nil
}
