package bsc

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sort"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof"
)

const AddressLength = 20
const ExtraVanityLength = 32
const ExtraSealLength = 65
const EpochLength = 200

func (b *Bridge) encodePoSAProof(transferEvent *c.BridgeTransfer) (*c.CheckPoSAPoSAProof, error) {
	var blocks []c.CheckPoSABlockPoSA
	var blocksMap = make(map[uint64]*c.CheckPoSABlockPoSA)
	var currEventBlockNum = transferEvent.Raw.BlockNumber
	var prevEventBlockNum uint64

	// get epochs changes blocks
	prevEventBlockNum, err := b.getPrevEventBlockNum()
	if err != nil {
		return nil, fmt.Errorf("get prev event block num: %w", err)
	}
	firstEpochChange := findFirstEpochChange(prevEventBlockNum)
	epochChanges := findEpochChanges(firstEpochChange, currEventBlockNum)

	// save blocks into blocksMap
	for _, epochChange := range epochChanges {
		// save epoch change block and get VS length
		encodedEpochChange, err := b.saveBlock(epochChange, blocksMap)
		if err != nil {
			return nil, fmt.Errorf("save epoch change block: %w", err)
		}
		vsLength := getVSLength(encodedEpochChange)

		// start from +1 cuz the epoch change block is already saved
		for i := uint64(1); i < uint64(vsLength)/2; i++ {
			if _, err := b.saveBlock(i+epochChange, blocksMap); err != nil {
				return nil, fmt.Errorf("save block %d: %w", i+epochChange, err)
			}
		}
	}

	// fill up blocks and get transfer event index
	indexToBlockNum := sortedKeys(blocksMap)
	var transferEventIndex uint64
	for i, blockNum := range indexToBlockNum {
		if blockNum == currEventBlockNum {
			transferEventIndex = uint64(i) // set transferEventIndex to index in blocks array
		}
		blocks = append(blocks, *blocksMap[blockNum])
	}

	// encode transferProof and save event block to blocksMap
	transfer, err := b.encodeTransferEvent(blocksMap, transferEvent)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferEvent: %w", err)
	}

	return &c.CheckPoSAPoSAProof{
		Blocks:             blocks,
		Transfer:           transfer,
		TransferEventBlock: transferEventIndex,
	}, nil
}

func (b *Bridge) getPrevEventBlockNum() (uint64, error) {
	lastEventId, err := b.GetLastEventId() // 70% шо тут треба не інпут, а аутпут
	if err != nil {
		return 0, fmt.Errorf("GetLastEventId: %w", err)
	}

	lastEvent, err := b.GetEventById(lastEventId)
	if err != nil {
		return 0, fmt.Errorf("GetEventById: %w", err) // TODO: обработать помилку коли останнього івенту нема
	}

	return lastEvent.Raw.BlockNumber, nil
}

func findFirstEpochChange(start uint64) uint64 {
	return uint64(math.Ceil(float64(start)/float64(EpochLength))) * EpochLength
}

func findEpochChanges(start, end uint64) []uint64 {
	var epochChanges []uint64
	for blockNum := start; blockNum < end; blockNum += 200 {
		epochChanges = append(epochChanges, blockNum)
	}
	return epochChanges
}

func getVSLength(encodedEpochChange *c.CheckPoSABlockPoSA) int {
	return (len(encodedEpochChange.ExtraData) - ExtraSealLength - ExtraVanityLength) / AddressLength
}

func (b *Bridge) saveBlock(blockNumber uint64, blocksMap map[uint64]*c.CheckPoSABlockPoSA) (*c.CheckPoSABlockPoSA, error) {
	if encodedBlock, ok := blocksMap[blockNumber]; ok {
		return encodedBlock, nil
	}

	block, err := b.Client.HeaderByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, fmt.Errorf("HeaderByNumber: %w", err)
	}
	encodedBlock, err := b.EncodeBlock(block)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	blocksMap[blockNumber] = encodedBlock
	return encodedBlock, nil
}

// TODO: винести в коммон
func (b *Bridge) getProof(event receipts_proof.ProofEvent) ([][]byte, error) {
	receipts, err := b.GetReceipts(event.Log().BlockHash)
	if err != nil {
		return nil, fmt.Errorf("GetReceipts: %w", err)
	}
	return receipts_proof.CalcProofEvent(receipts, event)
}

// used for 'ordered' map
// TODO: шось з цим теж зробити, мб заюзати дженеріки
func sortedKeys(m map[uint64]*c.CheckPoSABlockPoSA) []uint64 {
	keys := make([]uint64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func (b *Bridge) encodeTransferEvent(blocks map[uint64]*c.CheckPoSABlockPoSA, event *c.BridgeTransfer) (c.CommonStructsTransferProof, error) {
	proof, err := b.getProof(event)
	if err != nil {
		return c.CommonStructsTransferProof{}, err
	}

	if _, err := b.saveBlock(event.Raw.BlockNumber, blocks); err != nil {
		return c.CommonStructsTransferProof{}, err
	}

	return c.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}
