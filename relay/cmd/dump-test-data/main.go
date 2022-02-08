package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"relay/config"
	"relay/networks/amb"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// Nerwork URL.
	url string = "https://network.ambrosus-dev.io"
	// Default path to save tests data file.
	defaultPath string = "receipts_proof/fixtures/data.json"
)

var (
	// Path to save tests data file.
	pathToSave string
	// Transactions hash.
	txHash common.Hash = common.HexToHash("0xfe802d29486f3648fd082ecad8c8455aa751b18f36be7cff3cfd65245253233a")
)

// Write data structure.
type Data struct {
	Log      *types.Log
	Header   *types.Header
	Receipts []*types.Receipt
}

func init() {
	flag.StringVar(&pathToSave, "path", defaultPath, "Path to save test data dump.")
	flag.Parse()
}

func main() {
	// Creating a new ambrosus bridge.
	bridge := amb.New(&config.Bridge{Url: url})

	// Getting log receipts.
	logReceipt, err := bridge.Client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatalf("error getting log receipts: %s", err.Error())
	}

	logs := logReceipt.Logs[3]

	// Getting block by hash.
	block, err := bridge.Client.BlockByHash(context.Background(), logs.BlockHash)
	if err != nil {
		log.Fatalf("error getting block by hash: %s", err.Error())
	}

	// Getting bridge receipts.
	receipts, err := bridge.GetReceipts(logs.BlockHash)
	if err != nil {
		log.Fatalf("error getting receipts: %s", err.Error())
	}

	data := Data{Log: logs, Header: block.Header(), Receipts: receipts}

	// Wrire data to json file.
	if err := writeToFile(data); err != nil {
		log.Fatalf("error writing data to file: %s", err.Error())
	}
}

// Wrire data to json file.
func writeToFile(data Data) error {
	// Marshal data.
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Fatalf("error marshal data: %s", err.Error())
	}

	// Write data to file.
	return ioutil.WriteFile(pathToSave, file, 0644)
}
