package ethash

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/consensus/ethash"
)

func CheckDatasetPath(path string, epoch uint64, dir string) error {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			MakeDAG(epoch*epochLength, dir)

			return nil
		}

		return err
	}

	return nil
}

func DeleteDatasetFile(path string) error {
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	return nil
}

func PathToDataset(epoch uint64, dir string) string {
	seed := ethash.SeedHash(epoch*epochLength + 1)

	var endian string

	if !isLittleEndian() {
		endian = ".be"
	}

	return filepath.Join(dir, fmt.Sprintf("full-R%d-%x%s", 23, seed[:8], endian))
}
