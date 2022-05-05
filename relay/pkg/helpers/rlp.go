package helpers

func RlpPrefix(bytesCount int) []byte {
	if bytesCount <= 0x37 {
		return []byte{byte(0x80 + bytesCount)}
	}
	rlpBytesLen := 0
	for i := 1; bytesCount/i != 0; {
		rlpBytesLen++
		i *= 256
	}

	res := make([]byte, rlpBytesLen+1)
	res[0] = byte(0x80 + rlpBytesLen + 0x37)
	for i := rlpBytesLen; i > 0; i-- {
		res[i] = byte(bytesCount % 256)
		bytesCount /= 256
	}
	return res
}
