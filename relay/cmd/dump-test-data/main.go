package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

func main() {
	// Ambrosus configgurations for generating testing data.
	ambCfg := networkConfig{
		Name: "amb",
		Url:  "https://network.ambrosus-dev.io",
		Tx:   common.HexToHash("0xfe802d29486f3648fd082ecad8c8455aa751b18f36be7cff3cfd65245253233a"),
	}
	// Ethereum configgurations for generating testing data.
	ethCfg := networkConfig{
		Name:  "eth",
		Url:   "https://rinkeby.infura.io/v3/01117e8ede8e4f36801a6a838b24f36c",
		Tx:    common.HexToHash("0x3ccead6d9ce5dba1add833c4b2631ad2066aad6d2ef7bfe278eb259a7355d9e2"),
		Block: common.HexToHash("0xc4ca0efd5d528d67691abd9e10e9d4ca570f16235779e1f314b036caa5b455a1"),
	}

	// Generating ambrosus testing data.
	if err := generateNetworkTestData(ambCfg); err != nil {
		log.Error().Err(err).Msgf("error generating '%s' test data", ambCfg.Name)
	}

	// Generating ethereum testing data.
	if err := generateNetworkTestData(ethCfg); err != nil {
		log.Error().Err(err).Msgf("error generating '%s' test data", ethCfg.Name)
	}

	// Generating dataset testing data.
	if err := generateDatasetTestData(ethCfg); err != nil {
		log.Error().Err(err).Msg("error generating dataset test data")
	}
}

// Network testing data.
type networkData struct {
	Log      *types.Log
	Header   *types.Header
	Receipts []*types.Receipt
}

// Configurations for generating network testing data.
type networkConfig struct {
	Name  string      // Network name.
	Url   string      // Network URL.
	Tx    common.Hash // Network transaction hash.
	Block common.Hash // Network block hash.
}

// Generating network testing data.
func generateNetworkTestData(cfg networkConfig) error {
	// Creating a new ambrosus bridge client.
	bridge, err := amb.New(&config.AMBConfig{Network: config.Network{URL: cfg.Url}})
	if err != nil {
		return err
	}
	defer bridge.Client.Close()

	// Getting log receipts.
	logReceipt, err := bridge.Client.TransactionReceipt(context.Background(), cfg.Tx)
	if err != nil {
		return err
	}

	logs := logReceipt.Logs[3]

	// Getting block by hash.
	block, err := bridge.Client.BlockByHash(context.Background(), logs.BlockHash)
	if err != nil {
		return err
	}

	// Getting receipts from block.
	receipts, err := bridge.GetReceipts(logs.BlockHash)
	if err != nil {
		return err
	}

	data := networkData{Log: logs, Header: block.Header(), Receipts: receipts}

	return writeToJSONFile(data, fmt.Sprintf("./receipts_proof/fixtures/%s-data.json", cfg.Name))
}

// Dataset testing data.
type datasetData struct {
	BlockHash        common.Hash
	EpochData        *ethash.EpochData
	DataSetLookUp    []*big.Int
	WitnessForLookup []*big.Int
}

// Generating dataset testing data.
func generateDatasetTestData(cfg networkConfig) error {
	// Creating a new ethereum bridge client.
	bridge, err := eth.New(&config.ETHConfig{
		Network: config.Network{URL: "https://mainnet.infura.io/v3/ab050ca98686478e9e9b06dfc3b2f069"},
	})
	if err != nil {
		return err
	}
	defer bridge.Client.Close()

	// Getting block by hash.
	block, err := bridge.Client.BlockByHash(context.Background(), cfg.Block)
	if err != nil {
		log.Error().Err(err).Msg("block not getting")
	}

	// Encode block header without nonce.
	blockHeaderWithoutNonce, err := eth.EncodeHeaderWithoutNonceToRLP(block.Header())
	if err != nil {
		log.Error().Err(err).Msg("block header not encode")
	}

	blockHeaderHashWithoutNonce := crypto.Keccak256(blockHeaderWithoutNonce)

	var blockHeaderHashWithoutNonceLength32 [32]byte
	copy(blockHeaderHashWithoutNonceLength32[:], blockHeaderHashWithoutNonce)

	// Creating a new meta data block.
	blockMetaData := ethash.NewBlockMetaData(
		block.Header().Number.Uint64(), block.Header().Nonce.Uint64(),
		blockHeaderHashWithoutNonceLength32,
	)

	// Creating dataset lookup.
	dataSetLookUp := blockMetaData.DAGElementArray()
	// Creating witness for lookup
	witnessForLookup := blockMetaData.DAGProofArray()

	epoch := block.Header().Number.Uint64() / 30000
	// Generating epoch data.
	epochData, err := ethash.GenerateEpochData(epoch)
	if err != nil {
		return err
	}

	data := datasetData{
		BlockHash:        cfg.Block,
		EpochData:        epochData,
		DataSetLookUp:    dataSetLookUp,
		WitnessForLookup: witnessForLookup,
	}

	return writeToJSONFile(data, fmt.Sprintf("./assets/testdata/epoch-%d.json", epoch))
}

// Wrire data to json file.
func writeToJSONFile(data interface{}, path string) error {
	// Marshal data.
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatal().Err(err)
	}

	// Write data to file.
	return ioutil.WriteFile(path, file, 0644)
}
