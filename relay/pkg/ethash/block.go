package ethash

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash/merkle"
)

func ProcessDuringRead(datasetPath string, mt *merkle.DatasetTree) error {
	var (
		file *os.File
		err  error
	)

	file, err = os.Open(datasetPath)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	buf := [128]byte{}

	_, err = io.ReadFull(reader, buf[:8])
	if err != nil {
		return err
	}

	var i uint32 = 0

	for {
		n, err := io.ReadFull(reader, buf[:128])
		if n == 0 {
			if err == nil {
				continue
			} else if err == io.EOF {
				break
			}

			return err
		} else if n != 128 {
			return errors.New("error malformed dataset")
		}

		mt.Insert(merkle.Word(buf), i)

		if err != nil && err != io.EOF {
			return err
		}

		i++
	}

	return nil
}
