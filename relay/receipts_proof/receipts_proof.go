package receipts_proof

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"relay/helpers"
	"relay/receipts_proof/mytrie"
)

func CheckProof(proof [][]byte, log *types.Log) common.Hash {
	el := helpers.BytesConcat(
		proof[0],
		log.Address.Bytes(),
		proof[1],
		log.Topics[0].Bytes(), // event id
		proof[2],
		log.Data, //transfers
		proof[3],
	)

	for i := 4; i < len(proof); i += 2 {
		el = helpers.BytesConcat(proof[i], el, proof[i+1])
		//fmt.Printf("%x\n", el)
		if len(el) > 32 {
			el = mytrie.Hash(el)
		}
		//fmt.Printf("%x\n", el)
	}

	return common.BytesToHash(el)
}

func CalcProof(receipts []*types.Receipt, log *types.Log) ([][]byte, error) {
	// rlp encoded receipt with given log in it
	rlpReceipt, err := rlpBytes(receipts[log.TxIndex])
	if err != nil {
		panic(err)
	}

	logPart, err := encodeLog(rlpReceipt, log)
	if err != nil {
		return nil, err
	}
	triePart, err := encodeTrie(rlpReceipt, receipts)
	if err != nil {
		return nil, err
	}

	return append(logPart, triePart...), nil

}

func encodeLog(rlpReceipt []byte, log *types.Log) ([][]byte, error) {
	// encode log first (log is part of receipt) for accuracy
	rlpLog, err := rlpBytes(log)
	if err != nil {
		return nil, err
	}
	if len(log.Topics) != 1 {
		return nil, fmt.Errorf("log.Topics length != 1 (only event_id must be indexed)")
	}

	splitEls := [][]byte{
		log.Address.Bytes(),
		log.Topics[0].Bytes(),
		log.Data,
	}
	splittedLog, err := helpers.BytesSplit(rlpLog, splitEls)
	if err != nil {
		return nil, err
	}

	// wrap log with receipt:
	// (receipt1, log1), address, log2, event_id, log3, data, (log4, receipt2)

	splittedReceipt := bytes.Split(rlpReceipt, rlpLog)
	if len(splittedReceipt) != 2 {
		panic("split not 2")
	}

	splittedLog[0] = append(splittedReceipt[0], splittedLog[0]...)
	splittedLog[2] = append(splittedLog[2], splittedReceipt[1]...)

	return splittedLog, nil
}

type trieProof struct {
	whatSearch []byte
	result     [][]byte
	err        error
}

func encodeTrie(rlpReceipt []byte, receipts []*types.Receipt) ([][]byte, error) {
	trie := mytrie.NewStackTrie()
	types.DeriveSha(types.Receipts(receipts), trie)

	p := trieProof{whatSearch: rlpReceipt, result: [][]byte{}}
	r := p.trieProof(trie)
	if p.err != nil {
		return nil, p.err
	}
	if !r {
		return nil, fmt.Errorf("rlpReceipt not found in trie")
	}
	return p.result, nil
}

// trieProof try to find p.whatSearch in receipt tree
// when found - return to the root of the tree, writing the result along the way:
// push to result node.UnhashedVal splitted by child.Val => push some_bytes_before and some_bytes_after
// ex.: node.UnhashedVal is {some_bytes_before + child.Val + some_bytes_after}
func (p *trieProof) trieProof(node *mytrie.ModifiedStackTrie) bool {
	if node == nil {
		return false
	}
	if bytes.Equal(node.UnhashedVal, p.whatSearch) {
		return true
	}

	for _, child := range node.Children {
		if p.trieProof(child) {

			r := bytes.Split(child.Val, node.UnhashedVal)
			if len(r) != 2 {
				p.err = fmt.Errorf("split not 2")
				return false
			}
			p.result = append(p.result, r[0], r[1])

			return true
		}
	}

	return false
}

func rlpBytes(item rlp.Encoder) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := item.EncodeRLP(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
