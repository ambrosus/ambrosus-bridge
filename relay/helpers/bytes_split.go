package helpers

import (
	"bytes"
	"fmt"
)

func BytesSplit(input []byte, seps [][]byte) ([][]byte, error) {
	// example:
	//  input = "0xdead00beef2242
	//  seps = [0x00, 0x22]
	//  result = [0xDEAD. 0xbeef, 0x42]

	result := make([][]byte, len(seps)+1)

	for i, se := range seps {
		r := bytes.SplitN(input, se, 2)
		if len(r) != 2 {
			return nil, fmt.Errorf("no '%v' in '%v'", se, input)
		}

		result[i] = r[0]
		input = r[1]
	}
	result[len(seps)] = input
	return result, nil
}
