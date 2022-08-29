package ambrosus_explorer

import "math/big"

// BigInt is a wrapper over big.Int to implement only unmarshalText
// for json decoding.
type BigIntString big.Int

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *BigIntString) UnmarshalText(text []byte) (err error) {
	var bigInt = new(big.Int)
	err = bigInt.UnmarshalText(text)
	if err != nil {
		return
	}

	*b = BigIntString(*bigInt)
	return nil
}

// MarshalText implements the encoding.TextMarshaler
func (b *BigIntString) MarshalText() (text []byte, err error) {
	return []byte(b.Int().String()), nil
}

// Int returns b's *big.Int form
func (b *BigIntString) Int() *big.Int {
	return (*big.Int)(b)
}
