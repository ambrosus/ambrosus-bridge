package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

func main() {
	// todo как то выбрать шо именно нагенерить

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
	if err := dataForReceiptProof(ambCfg); err != nil {
		log.Error().Err(err).Msgf("error generating '%s' test data", ambCfg.Name)
	}

	// Generating ethereum testing data.
	if err := dataForReceiptProof(ethCfg); err != nil {
		log.Error().Err(err).Msgf("error generating '%s' test data", ethCfg.Name)
	}

	encodePoWBlock(ethCfg)
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
func dataForReceiptProof(cfg networkConfig) error {
	// Creating a new ambrosus bridge client.
	bridge, err := amb.New(&config.Bridge{Url: cfg.Url})
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

// Generating dataset testing data.
func encodePoWBlock(cfg networkConfig) {
	// Creating a new ethereum bridge client.
	bridge, err := eth.New(&config.Bridge{Url: "https://mainnet.infura.io/v3/ab050ca98686478e9e9b06dfc3b2f069"})
	if err != nil {
		log.Fatal().Err(err).Msg("block not getting")
	}
	defer bridge.Client.Close()

	// Getting block by hash.
	block, err := bridge.Client.BlockByHash(context.Background(), cfg.Block)
	if err != nil {
		log.Fatal().Err(err).Msg("block not getting")
	}

	data, err := eth.EncodeBlock(block.Header(), true)
	if err != nil {
		log.Fatal().Err(err).Msg("encode pow block err")
	}

	err = writeToJSONFile(data, fmt.Sprintf("./assets/testdata/BlockPoW.json"))
	if err != nil {
		log.Fatal().Err(err).Msg("encode pow block err")
	}

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
