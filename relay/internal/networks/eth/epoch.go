package eth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/rs/zerolog/log"
)

const (
	epochDataPath     string = "./assets/epoch/"
	epochDataFilePath        = epochDataPath + "%d.json"
)

var ErrEpochDataFileNotFound = errors.New("error epoch data file not found")

func (b *Bridge) SetEpochData(epochData *ethash.EpochData) error {
	var nodes []*big.Int
	start := big.NewInt(0)

	for i, node := range epochData.MerkleNodes {
		nodes = append(nodes, node)

		if len(nodes) == 40 || i == len(epochData.MerkleNodes) {
			merkelNodesNumber := big.NewInt(int64(len(nodes)))

			if i < 440 && epochData.Epoch.Uint64() == 128 {
				start.Add(start, merkelNodesNumber)
				nodes = []*big.Int{}

				continue
			}

			err := b.sideBridge.SubmitEpochData(
				epochData.Epoch, epochData.FullSizeIn128Resolution, epochData.BranchDepth,
				nodes, start, merkelNodesNumber,
			)
			if err != nil {
				return err
			}

			start.Add(start, merkelNodesNumber)
			nodes = []*big.Int{}
		}
	}

	return nil
}

func (b *Bridge) loadEpochDataFile(epoch uint64) (*ethash.EpochData, error) {
	data, err := os.ReadFile(fmt.Sprintf(epochDataFilePath, epoch))
	if err != nil {
		return nil, err
	}

	var epochData *ethash.EpochData

	if err := json.Unmarshal(data, &epochData); err != nil {
		return nil, err
	}

	return epochData, nil
}

func (b *Bridge) createEpochDataFile(epoch uint64) (*ethash.EpochData, error) {
	log.Debug().Msgf("creating '%d.json' epoch data file...", epoch)

	data, err := ethash.GenerateEpochData(epoch)
	if err != nil {
		return nil, err
	}

	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return data, err
	}

	if err := os.WriteFile(fmt.Sprintf(epochDataFilePath, epoch), file, 0644); err != nil {
		return data, err
	}

	path := ethash.PathToDataset(epoch, ethash.DefaultDir)
	if err := ethash.DeleteDatasetFile(path); err != nil {
		return data, err
	}

	return data, nil
}

func (b *Bridge) deleteEpochDataFile(epoch uint64) error {
	log.Debug().Msgf("Deleting '%d.json' epoch data file...", epoch)

	if err := os.Remove(fmt.Sprintf(epochDataFilePath, epoch)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	return nil
}

func (b *Bridge) checkEpochDataFile(epoch uint64) error {
	log.Debug().Msgf("Checking '%d.json' epoch data file...", epoch)

	if _, err := os.Stat(fmt.Sprintf(epochDataFilePath, epoch)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrEpochDataFileNotFound
		}

		return err
	}

	return nil
}

func (b *Bridge) getGeneratedEpochNumbers() ([]int, error) {
	files, err := ioutil.ReadDir(epochDataPath)
	if err != nil {
		return nil, err
	}

	var epochs []int

	for _, file := range files {
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
	log.Debug().Msg("Checking epoch data dir...")

	epochs, err := b.getGeneratedEpochNumbers()
	if err != nil {
		return err
	}

	for _, i := range epochs {
		if uint64(i) < epoch {
			if err := b.deleteEpochDataFile(uint64(i)); err != nil {
				return err
			}
		}
	}

	for i := epoch; i <= length; i++ {
		if err := b.checkEpochDataFile(i); err != nil {
			if errors.Is(err, ErrEpochDataFileNotFound) {
				if _, err := b.createEpochDataFile(epoch); err != nil {
					return err
				}

				continue
			}

			return err
		}
	}

	return nil
}
