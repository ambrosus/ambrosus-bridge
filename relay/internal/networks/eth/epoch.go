package eth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
)

const (
	epochDataPath     string = "./assets/epoch/"
	epochDataFilePath        = epochDataPath + "%d.json"
)

var ErrEpochDataFileNotFound = errors.New("error epoch data file not found")

func (b *Bridge) loadEpochDataFile(epoch uint64) (*ethash.EpochData, error) {
	b.Logger.Debug().Msgf("Loading '%d.json' epoch data file...", epoch)

	data, err := os.ReadFile(fmt.Sprintf(epochDataFilePath, epoch))
	if err != nil {
		return nil, fmt.Errorf("read file %v: %w", epochDataFilePath, err)
	}

	var epochData *ethash.EpochData

	if err := json.Unmarshal(data, &epochData); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return epochData, nil
}

func (b *Bridge) createEpochDataFile(epoch uint64) (*ethash.EpochData, error) {
	b.Logger.Info().Msgf("Creating '%d.json' epoch data file...", epoch)

	data, err := ethash.GenerateEpochData(epoch)
	if err != nil {
		return nil, fmt.Errorf("generate epoch data: %w", err)
	}

	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return data, fmt.Errorf("marshal: %w", err)
	}

	if err := os.WriteFile(fmt.Sprintf(epochDataFilePath, epoch), file, 0644); err != nil {
		return data, fmt.Errorf("write file %v: %w", epochDataFilePath, err)
	}

	b.Logger.Info().Msgf("Finish creating '%d.json' epoch data file...", epoch)

	return data, nil
}

func (b *Bridge) deleteEpochDataFile(epoch uint64) error {
	b.Logger.Debug().Msgf("Deleting '%d.json' epoch data file...", epoch)
	err := os.Remove(fmt.Sprintf(epochDataFilePath, epoch))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (b *Bridge) checkEpochDataFile(epoch uint64) error {
	b.Logger.Debug().Msgf("Checking '%d.json' epoch data file...", epoch)
	_, err := os.Stat(fmt.Sprintf(epochDataFilePath, epoch))
	if errors.Is(err, os.ErrNotExist) {
		return ErrEpochDataFileNotFound
	}
	return err
}

func (b *Bridge) getGeneratedEpochNumbers() ([]int, error) {
	files, err := ioutil.ReadDir(epochDataPath)
	if err != nil {
		return nil, err
	}

	var epochs []int

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		name := strings.Trim(file.Name(), ".json")

		epoch, err := strconv.Atoi(name)
		if err != nil {
			return nil, err
		}

		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

func (b *Bridge) checkEpochDataDir(epoch uint64, length uint64) error {
	b.Logger.Debug().Msg("Checking epoch data dir...")

	epochs, err := b.getGeneratedEpochNumbers()
	if err != nil {
		return fmt.Errorf("get epoch numbers: %w", err)
	}

	for _, i := range epochs {
		if uint64(i) < epoch {
			if err := b.deleteEpochDataFile(uint64(i)); err != nil {
				return fmt.Errorf("delete epoch (%v) data file: %w", i, err)
			}
		}
	}

	for i := epoch; i <= length+epoch; i++ {
		if err := b.checkEpochDataFile(i); err != nil {
			if errors.Is(err, ErrEpochDataFileNotFound) {
				if _, err := b.createEpochDataFile(i); err != nil {
					return fmt.Errorf("create epoch (%v) data file: %w", i, err)
				}

				continue
			}

			return fmt.Errorf("check epoch (%v) data file: %w", i, err)
		}
	}

	return nil
}
