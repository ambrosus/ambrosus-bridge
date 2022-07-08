package aura_proof

import (
	"context"
	"fmt"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/aura/aura_proof/rolling_finality"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type vsChangeInBlock struct {
	eventBlock     uint64
	finalizedBlock uint64
	changes        []c.CheckAuraValidatorSetChange
	lastEvent      *c.VsInitiateChange
}

func (e *AuraEncoder) getVsChanges(toBlock uint64) (map[uint64]*vsChangeInBlock, error) {
	start, err := e.getLastProcessedBlockNum()
	if err != nil {
		return nil, fmt.Errorf("getLastProcessedBlockNum: %w", err)
	}

	initialValidatorSet, err := e.auraReceiver.GetValidatorSet()
	if err != nil {
		return nil, fmt.Errorf("GetValidatorSet: %w", err)
	}

	vsChangeEvents, err := e.fetchVSChangeEvents(start, toBlock)
	if err != nil {
		return nil, err
	}

	blockToEvents, err := e.preprocessVSChangeEvents(initialValidatorSet, vsChangeEvents)
	if err != nil {
		return nil, err
	}

	err = e.findWhenFinalize(start, initialValidatorSet, blockToEvents)
	if err != nil {
		return nil, err
	}

	return blockToEvents, nil
}

func (e *AuraEncoder) findWhenFinalize(start uint64, initialValidatorSet []common.Address, blockToEvents map[uint64]*vsChangeInBlock) error {
	rf := rolling_finality.NewRollingFinality(initialValidatorSet)
	eventsToFinalize := len(blockToEvents)

	for blockNum := start; eventsToFinalize > 0; blockNum++ {
		block, err := e.fetchBlockCache(blockNum)
		if err != nil {
			return err
		}

		finalizedBlocks, err := rf.Push(block.Number.ToInt().Uint64(), *block.Coinbase)
		if err != nil {
			return err
		}

		for _, finalizedBlock := range finalizedBlocks {
			finalizedEvent, ok := blockToEvents[finalizedBlock]
			if ok { // finalized block has vs change event
				rf = rolling_finality.NewRollingFinality(finalizedEvent.lastEvent.NewSet) // OpenEthereum do the same
				finalizedEvent.finalizedBlock = blockNum
				eventsToFinalize--
			}
		}

	}

	return nil
}

func (e *AuraEncoder) preprocessVSChangeEvents(initialValidatorSet []common.Address, events []*c.VsInitiateChange) (map[uint64]*vsChangeInBlock, error) {
	blocksToEvents := map[uint64]*vsChangeInBlock{}

	prevSet := initialValidatorSet

	for _, event := range events {
		address, index, err := deltaVS(prevSet, event.NewSet)
		if err != nil {
			return nil, fmt.Errorf("deltaVS: %w", err)
		}
		vsChange := c.CheckAuraValidatorSetChange{DeltaAddress: address, DeltaIndex: index}
		prevSet = event.NewSet

		if _, ok := blocksToEvents[event.Raw.BlockNumber]; !ok {
			blocksToEvents[event.Raw.BlockNumber] = &vsChangeInBlock{eventBlock: event.Raw.BlockNumber}
		}
		blocksToEvents[event.Raw.BlockNumber].changes = append(blocksToEvents[event.Raw.BlockNumber].changes, vsChange)
		blocksToEvents[event.Raw.BlockNumber].lastEvent = event
	}

	return blocksToEvents, nil
}

func (e *AuraEncoder) fetchVSChangeEvents(start, end uint64) ([]*c.VsInitiateChange, error) {
	opts := &bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}

	logs, err := e.vsContract.FilterInitiateChange(opts, nil)
	if err != nil {
		return nil, fmt.Errorf("filter initiate changes: %w", err)
	}

	var res []*c.VsInitiateChange
	for logs.Next() {
		res = append(res, logs.Event)
	}

	return res, nil
}

func (e *AuraEncoder) getLastProcessedBlockNum() (uint64, error) {
	blockHash, err := e.auraReceiver.GetLastProcessedBlockHash()
	if err != nil {
		return 0, fmt.Errorf("GetLastProcessedBlockHash: %w", err)
	}

	block, err := e.bridge.GetClient().BlockByHash(context.Background(), *blockHash)
	if err != nil {
		return 0, fmt.Errorf("get rfBlock by hash: %w", err)
	}

	return block.Number().Uint64(), nil
}

func deltaVS(prev, curr []common.Address) (common.Address, int64, error) {
	d := len(curr) - len(prev)
	if d != 1 && d != -1 {
		return common.Address{}, 0, fmt.Errorf("delta has more (or less) than 1 change")
	}

	for i, prevEl := range prev {
		if i >= len(curr) { // deleted at the end
			return prev[i], int64(-i - 1), nil
		}

		if curr[i] != prevEl {
			if d == 1 { // added
				return curr[i], int64(i), nil
			} else { // deleted
				return prev[i], int64(-i - 1), nil
			}
		}
	}

	// add at the end
	i := len(curr) - 1
	return curr[i], int64(i), nil

	// return common.Address{}, 0, fmt.Errorf("this error shouln't exist")
}
