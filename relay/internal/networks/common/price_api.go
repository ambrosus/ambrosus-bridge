package common

import (
	"bytes"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto"
)

func (b *CommonBridge) Sign(tokenPrice float64, tokenAddress string) ([]byte, error) {
	var data bytes.Buffer
	if err := binary.Write(&data, binary.LittleEndian, tokenPrice); err != nil {
		return nil, err
	}
	data.WriteString(tokenAddress)
	return crypto.Sign(crypto.Keccak256(data.Bytes()), b.Pk)
}

func (b *CommonBridge) GetPrice(tokenAddress string) (float64, error) {
	return 13.37, nil
}
