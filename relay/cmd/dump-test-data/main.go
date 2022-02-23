package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

var (
	// Ambrosus configgurations for generating testing data.
	ambCfg = networkConfig{
		Name:  "amb",
		Url:   "https://network.ambrosus-dev.io",
		Tx:    common.HexToHash("0xfe802d29486f3648fd082ecad8c8455aa751b18f36be7cff3cfd65245253233a"),
		Block: common.HexToHash("0x860249608cfb868c93cd960718d2bc658cc53ec1e9482b4a051f5c2a5c5f6922"),
	}
	// Ethereum configgurations for generating testing data.
	ethCfg = networkConfig{
		Name:  "eth",
		Url:   "https://rinkeby.infura.io/v3/01117e8ede8e4f36801a6a838b24f36c",
		Tx:    common.HexToHash("0x3ccead6d9ce5dba1add833c4b2631ad2066aad6d2ef7bfe278eb259a7355d9e2"),
		Block: common.HexToHash("0xc4ca0efd5d528d67691abd9e10e9d4ca570f16235779e1f314b036caa5b455a1"),
	}
)

func main() {
	switch os.Args[1] {
	case "receipts-proof":
		receiptsProofInput(os.Args[2])
	case "pow-block":
		// Encoding PoW block.
		if err := encodePoWBlock(ethCfg); err != nil {
			log.Fatal().Err(err).Msg("error encoding pow block")
		}

		return
	case "poa-block":
		// Encoding PoA block.
		if err := encodePoABlock(ambCfg); err != nil {
			log.Fatal().Err(err).Msg("error encoding poa block")
		}

		return
	default:
		log.Warn().Msg("command not found")
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

// Receipt proof input.
func receiptsProofInput(args string) {
	switch args {
	case "amb":
		// Generating ambrosus testing data for receipts proof.
		if err := dataForReceiptProof(ambCfg); err != nil {
			log.Fatal().Err(err).Msgf("error generating '%s' test data for receipts proof", ambCfg.Name)
		}

		return
	case "eth":
		// Generating ethereum testing data for receitps proof.
		if err := dataForReceiptProof(ethCfg); err != nil {
			log.Fatal().Err(err).Msgf("error generating '%s' test data for receipts proof", ethCfg.Name)
		}

		return
	default:
		log.Warn().Msg("Specify the network tag!")
	}
}

// Generating network testing data for receipts proof.
func dataForReceiptProof(cfg networkConfig) error {
	log.Info().Msgf("Generating data for %s receipts proof..", cfg.Name)

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

// Encoding PoW block.
func encodePoWBlock(cfg networkConfig) error {
	log.Info().Msg("Encoding pow block...")

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
		return err
	}

	data, err := eth.EncodeBlock(block.Header(), true)
	if err != nil {
		return err
	}

	return writeToJSONFile(data, fmt.Sprintf("./assets/testdata/BlockPoW-%d.json", block.Header().Number.Uint64()))
}

// Encoding PoA block.
func encodePoABlock(cfg networkConfig) error {
	// Creating a new ambrosus bridge client.
	bridge, err := amb.New(&config.AMBConfig{
		Network: config.Network{URL: "https://network.ambrosus.io"},
	})
	if err != nil {
		return err
	}
	defer bridge.Client.Close()

	// Getting block by hash.
	block, err := bridge.Client.BlockByHash(context.Background(), cfg.Block)
	if err != nil {
		return err
	}

	header, err := bridge.HeaderByNumber(block.Number()) // TODO: delete
	if err != nil {
		return err
	}

	data, err := amb.EncodeBlock(header)
	if err != nil {
		return err
	}

	return writeToJSONFile(data, fmt.Sprintf("./assets/testdata/BlockPoA-%d.json", block.Header().Number.Uint64()))
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
