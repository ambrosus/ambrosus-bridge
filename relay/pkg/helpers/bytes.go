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

func BytesToBytes32(bytes []byte) (bytes32 [32]byte) {
	copy(bytes32[:], bytes)
	return
}

func BytesToBytes3(bytes []byte) (bytes32 [3]byte) {
	copy(bytes32[:], bytes)
	return
}
func BytesToBytes4(bytes []byte) (bytes32 [4]byte) {
	copy(bytes32[:], bytes)
	return
}

func BytesConcat(slices ...[]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}

	res := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(res[i:], s)
	}
	return res
}
