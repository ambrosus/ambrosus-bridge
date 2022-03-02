package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/amb"
	"github.com/ambrosus/ambrosus-bridge/relay/networks/eth"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

var (
	// Ambrosus configgurations for generating testing data.
	ambCfg = networkConfig{
		Name: "amb",
		Url:  "https://network.ambrosus-dev.io",
		Tx:   common.HexToHash("0xfe802d29486f3648fd082ecad8c8455aa751b18f36be7cff3cfd65245253233a"),
	}
	// Ethereum configgurations for generating testing data.
	ethCfg = networkConfig{
		Name: "eth",
		Url:  "https://rinkeby.infura.io/v3/01117e8ede8e4f36801a6a838b24f36c",
		Tx:   common.HexToHash("0x3ccead6d9ce5dba1add833c4b2631ad2066aad6d2ef7bfe278eb259a7355d9e2"),
	}

	ErrArgNotFound = errors.New("error argument not found")
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
	case "epoch":
		// Generating epoch data.
		if err := generateEpochData(); err != nil {
			log.Fatal().Err(err).Msg("error generating epoch data")
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
	Name string      // Network name.
	Url  string      // Network URL.
	Tx   common.Hash // Network transaction hash.
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
	receipts, err := ethereum.GetReceipts(bridge.Client, logs.BlockHash)
	if err != nil {
		return err
	}

	data := networkData{Log: logs, Header: block.Header(), Receipts: receipts}

	return writeToJSONFile(data, fmt.Sprintf("./receipts_proof/fixtures/%s-data.json", cfg.Name))
}

type powData struct {
	P1                  string
	ParentOrReceiptHash string
	P2                  string
	Difficulty          string
	P3                  string
	Number              string
	P4                  string
	P5                  string
	Nonce               string
	P6                  string
	DataSetLookUp       []string
	WitnessForLookUp    []string
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

	if len(os.Args) < 3 {
		return ErrArgNotFound
	}

	number, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}

	// Getting block by hash.
	block, err := bridge.Client.BlockByNumber(context.Background(), big.NewInt(int64(number)))
	if err != nil {
		return err
	}

	bd, err := eth.EncodeBlock(block.Header(), true)
	if err != nil {
		return err
	}

	data := powData{
		P1:                  "0x" + common.Bytes2Hex(bd.P1),
		ParentOrReceiptHash: "0x" + common.Bytes2Hex(bd.ParentOrReceiptHash[:]),
		P2:                  "0x" + common.Bytes2Hex(bd.P2),
		Difficulty:          "0x" + common.Bytes2Hex(bd.Difficulty),
		P3:                  "0x" + common.Bytes2Hex(bd.P3),
		Number:              "0x" + common.Bytes2Hex(bd.Number),
		P4:                  "0x" + common.Bytes2Hex(bd.P4),
		P5:                  "0x" + common.Bytes2Hex(bd.P5),
		Nonce:               "0x" + common.Bytes2Hex(bd.Nonce),
		P6:                  "0x" + common.Bytes2Hex(bd.P6),
		DataSetLookUp:       bigIntArrayToStringArray(bd.DataSetLookUp),
		WitnessForLookUp:    bigIntArrayToStringArray(bd.WitnessForLookUp),
	}

	return writeToJSONFile(data, fmt.Sprintf("./assets/testdata/BlockPoW-%d.json", block.Header().Number.Uint64()))
}

type poaData struct {
	P0Seal      string
	P0Bare      string
	P1          string
	ParentHash  string
	P2          string
	ReceiptHash string
	P3          string
	S1          string
	Step        string
	S2          string
	Signature   string
	Type        *big.Int
}

// Encoding PoA block.
func encodePoABlock(cfg networkConfig) error {
	log.Info().Msg("Encoding poa block...")

	// Creating a new ambrosus bridge client.
	bridge, err := amb.New(&config.AMBConfig{
		Network: config.Network{URL: "https://network.ambrosus.io"},
	})
	if err != nil {
		return err
	}
	defer bridge.Client.Close()

	if len(os.Args) < 3 {
		return ErrArgNotFound
	}

	number, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}

	header, err := bridge.HeaderByNumber(big.NewInt(int64(number)))
	if err != nil {
		return err
	}

	bd, err := amb.EncodeBlock(header)
	if err != nil {
		return err
	}

	data := poaData{
		P0Seal:      "0x" + common.Bytes2Hex(bd.P0Seal),
		P0Bare:      "0x" + common.Bytes2Hex(bd.P0Bare),
		P1:          "0x" + common.Bytes2Hex(bd.P1),
		ParentHash:  "0x" + common.Bytes2Hex(bd.ParentHash[:]),
		P2:          "0x" + common.Bytes2Hex(bd.P2),
		ReceiptHash: "0x" + common.Bytes2Hex(bd.ReceiptHash[:]),
		P3:          "0x" + common.Bytes2Hex(bd.P3),
		S1:          "0x" + common.Bytes2Hex(bd.S1),
		Step:        "0x" + common.Bytes2Hex(bd.Step),
		S2:          "0x" + common.Bytes2Hex(bd.S2),
		Signature:   "0x" + common.Bytes2Hex(bd.Signature),
		Type:        bd.Type,
	}

	return writeToJSONFile(data, fmt.Sprintf("./assets/testdata/BlockPoA-%d.json", number))
}

type epochData struct {
	Epoch                   *big.Int
	FullSizeIn128Resolution *big.Int
	BranchDepth             *big.Int
	MerkleNodes             []string
}

// Generating epoch data.
func generateEpochData() error {
	log.Info().Msg("Generating epoch data...")

	if len(os.Args) < 3 {
		return ErrArgNotFound
	}

	number, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}

	ed, err := ethash.GenerateEpochData(uint64(number))
	if err != nil {
		return err
	}

	data := epochData{
		Epoch:                   ed.Epoch,
		FullSizeIn128Resolution: ed.FullSizeIn128Resolution,
		BranchDepth:             ed.BranchDepth,
		MerkleNodes:             bigIntArrayToStringArray(ed.MerkleNodes),
	}

	return writeToJSONFile(data, fmt.Sprintf("./assets/testdata/epoch-%d.json", number))
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

func bigIntArrayToStringArray(intArray []*big.Int) []string {
	var array []string

	for _, n := range intArray {
		array = append(array, n.String())
	}

	return array
}
