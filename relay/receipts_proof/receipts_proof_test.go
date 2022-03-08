package receipts_proof

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
)

// Data file structure.
type Data struct {
	Log      *types.Log
	Header   *types.Header
	Receipts []*types.Receipt
}

// Testing receipts proof.
func TestReceiptsProof(t *testing.T) {
	// Tests structures.
	tests := []struct{ name, fixtures string }{
		{name: "AMB", fixtures: "fixtures/amb-data.json"},
		{name: "ETH", fixtures: "fixtures/eth-data.json"},
	}

	// Conducting tests in various structures.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Loading data file.
			data, err := loadDataFile(tt.fixtures)
			if err != nil {
				t.Fatalf("error loading data file: %s", err.Error())
			}

			logElements := [][]byte{
				data.Log.Address.Bytes(),
				data.Log.Topics[1].Bytes(),
				data.Log.Data,
			}

			// Calculate proof.
			proof, err := CalcProof(data.Receipts, data.Log, logElements)
			if err != nil {
				t.Errorf("error calculate proof: %s", err.Error())
			}

			// Check for similarity of a receipt hash.
			receiptHash := CheckProof(proof, logElements)
			if receiptHash != data.Header.ReceiptHash {
				t.Error("error proof check failed")
			}
		})
	}
}

// Loading data file.
func loadDataFile(path string) (*Data, error) {
	var data Data

	// Reading data file.
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Unmarshal data file.
	if err = json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
