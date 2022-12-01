package helpers

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type BigInt struct {
	BigInt *big.Int
}

func (i *BigInt) UnmarshalJSON(b []byte) error {
	var val string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}

	bi, ok := new(big.Int).SetString(val, 10)
	if !ok {
		return fmt.Errorf("invalid number: %s", val)
	}
	i.BigInt = bi
	return nil
}
