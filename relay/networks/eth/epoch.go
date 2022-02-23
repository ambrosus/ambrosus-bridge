package eth

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
)

func (b *Bridge) SetEpochData(epochData ethash.EpochData) error {
	var nodes []*big.Int
	start := big.NewInt(0)

	for i, node := range epochData.MerkleNodes {
		nodes = append(nodes, node)

		if len(nodes) == 40 || i == len(epochData.MerkleNodes) {
			merkelNodesNumber := big.NewInt(int64(len(nodes)))

			if i < 440 && epochData.Epoch.Uint64() == 128 {
				start.Add(start, merkelNodesNumber)
				nodes = []*big.Int{}

				continue
			}

			// TODO
			err := sideBridgeSetEpochDataMock(
				epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth,
				nodes, start, merkelNodesNumber,
			)
			if err != nil {
				return err
			}

			start.Add(start, merkelNodesNumber)
			nodes = []*big.Int{}
		}
	}

	return nil
}

func sideBridgeSetEpochDataMock(*big.Int, *big.Int, *big.Int, []*big.Int, *big.Int, *big.Int) error {
	return nil
}
