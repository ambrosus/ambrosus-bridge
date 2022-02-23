package eth

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

func (b *Bridge) SetEpochData(epochData ethash.EpochData) error {
	nodes := []*big.Int{}
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

			auth, err := b.prepareTransaction(
				crypto.PubkeyToAddress(b.config.PrivateKey.PublicKey),
				b.config.PrivateKey, big.NewInt(0),
			)
			if err != nil {
				return err
			}

			// TODO: Change faceFunction() to contranct bind
			tx, err := faceFunction(
				auth, epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth,
				nodes, start, merkelNodesNumber,
			)
			if err != nil {
				return err
			}

			log.Info().Str("hex", tx.Hash().Hex()).Msg("Transaction submitted!")

			receipts, err := b.GetTransactionReceipt(tx.Hash())
			if err != nil {
				return err
			}

			if receipts.Status == 0 {
				// TODO: return transaction failed error and send log.
				return nil
			}

			start.Add(start, merkelNodesNumber)
			nodes = []*big.Int{}
		}
	}

	return nil
}

func faceFunction(*bind.TransactOpts, *big.Int, *big.Int, *big.Int, []*big.Int, *big.Int, *big.Int) (*types.Transaction, error) {
	return &types.Transaction{}, nil
}
