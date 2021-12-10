package mytrie

func hexToCompact(hex []byte) []byte {
	terminator := byte(0)
	if hasTerm(hex) {
		terminator = 1
		hex = hex[:len(hex)-1]
	}
	buf := make([]byte, len(hex)/2+1)
	buf[0] = terminator << 5 // the flag byte
	if len(hex)&1 == 1 {
		buf[0] |= 1 << 4 // odd flag
		buf[0] |= hex[0] // first nibble is contained in the first byte
		hex = hex[1:]
	}
	decodeNibbles(hex, buf[1:])
	return buf
}

// hexToCompactInPlace places the compact Key in input buffer, returning the length
// needed for the representation
func hexToCompactInPlace(hex []byte) int {
	var (
		hexLen    = len(hex) // length of the hex input
		firstByte = byte(0)
	)
	// Check if we have a terminator there
	if hexLen > 0 && hex[hexLen-1] == 16 {
		firstByte = 1 << 5
		hexLen-- // last part was the terminator, ignore that
	}
	var (
		binLen = hexLen/2 + 1
		ni     = 0 // index in hex
		bi     = 1 // index in bin (compact)
	)
	if hexLen&1 == 1 {
		firstByte |= 1 << 4 // odd flag
		firstByte |= hex[0] // first nibble is contained in the first byte
		ni++
	}
	for ; ni < hexLen; bi, ni = bi+1, ni+2 {
		hex[bi] = hex[ni]<<4 | hex[ni+1]
	}
	hex[0] = firstByte
	return binLen
}

func keybytesToHex(str []byte) []byte {
	l := len(str)*2 + 1
	var nibbles = make([]byte, l)
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	nibbles[l-1] = 16
	return nibbles
}

func decodeNibbles(nibbles []byte, bytes []byte) {
	for bi, ni := 0, 0; ni < len(nibbles); bi, ni = bi+1, ni+2 {
		bytes[bi] = nibbles[ni]<<4 | nibbles[ni+1]
	}
}

// hasTerm returns whether a hex Key has the terminator flag.
func hasTerm(s []byte) bool {
	return len(s) > 0 && s[len(s)-1] == 16
}
