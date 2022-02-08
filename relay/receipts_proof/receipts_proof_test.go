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
	// Testing args.
	type args struct{ path string }

	// Tests structures.
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "OK", args: args{path: "fixtures/data.json"}},
	}

	// Conducting tests in various structures.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Loading data file.
			data, err := loadDataFile(tt.args.path)
			if err != nil {
				t.Fatalf("error loading data file: %s", err.Error())
			}

			// Calculate proof.
			proof, err := CalcProof(data.Receipts, data.Log)
			if err != nil {
				t.Errorf("error calculate proof: %s", err.Error())
			}

			// Check for similarity of a receipt hash.
			receipt := CheckProof(proof, data.Log)
			if receipt != data.Header.ReceiptHash {
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

	// Unmarsahl data file.
	if err = json.Unmarshal([]byte(file), &data); err != nil {
		return nil, err
	}

	return &data, nil
}
