package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Bridge struct {
	Client     *ethclient.Client
	Contract   *contracts.Eth
	sideBridge networks.Bridge
	config     *config.Bridge
}

// Creating a new ethereum bridge.
func New(cfg *config.Bridge) (*Bridge, error) {
	// Creating a new ethereum client.
	client, err := ethclient.Dial(cfg.Url)
	if err != nil {
		return nil, err
	}

	// Creating a new ethereum bridge contract instance.
	contract, err := contracts.NewEth(cfg.ContractAddress, client)
	if err != nil {
		return nil, err
	}

	return &Bridge{Client: client, Contract: contract, config: cfg}, nil
}

func (b *Bridge) SubmitTransfer(proof contracts.TransferProof) error {
	switch proof.(type) {
	case *contracts.CheckAuraAuraProof:
		// todo
	default:
		// todo error

	}
	return nil
}

func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

func (b *Bridge) DisputeBlock(blockHash common.Hash) {
	block, err := b.Client.BlockByHash(context.Background(), blockHash)
	if err != nil {
		log.Error().Err(err).Msg("block not getting")
	}

	blockHeaderWithoutNonce, err := EncodeHeaderWithoutNonceToRLP(block.Header())
	if err != nil {
		log.Error().Err(err).Msg("block header not encode")
	}

	blockHeaderHashWithoutNonce := crypto.Keccak256(blockHeaderWithoutNonce)

	var blockHeaderHashWithoutNonceLength32 [32]byte
	copy(blockHeaderHashWithoutNonceLength32[:], blockHeaderHashWithoutNonce)

	blockMetaData := ethash.NewBlockMetaData(
		block.Header().Number.Uint64(), block.Header().Nonce.Uint64(),
		blockHeaderHashWithoutNonceLength32,
	)

	dataSetLookUp := blockMetaData.DAGElementArray()
	witnessForLookup := blockMetaData.DAGProofArray()

	fmt.Println(dataSetLookUp)
	fmt.Println(witnessForLookup)
}

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.Bridge) {}
