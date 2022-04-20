package ethash

import (
	"encoding/binary"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/bitutil"
	"golang.org/x/crypto/sha3"
)

const (
	cacheRounds = 3 // Number of rounds in cache production
)

func (e *Ethash) getCache(epoch uint64) ([]byte, error) {
	if cache, ok := e.caches[epoch]; ok {
		e.logger.Debug().Msg("Loaded old ethash cache from cache :)")
		return cache, nil
	}

	if e.useFs() {
		// Try to load the file from disk and memory map it
		cache, err := os.ReadFile(e.pathToCache(epoch))
		if err == nil {
			e.logger.Debug().Msg("Loaded old ethash cache from disk")
			e.caches[epoch] = cache
			return cache, nil
		}
	}

	// Generate new
	cache, err := e.generateCache(epoch)
	if err != nil {
		return nil, err
	}

	if e.useFs() {
		if err = dumpToFile(e.pathToCache(epoch), cache); err != nil {
			e.logger.Warn().Err(err).Msg("Failed to save cache file")
		}
	}

	e.caches[epoch] = cache
	return cache, nil
}

func (e *Ethash) generateCache(epoch uint64) ([]byte, error) {
	e.logger.Info().Msg("Start generating ethash cache")
	start := time.Now()

	seed := seedHash(epoch)
	size := cacheSize(epoch)

	buffer := make([]byte, size)
	generateCache(buffer, seed)

	e.logger.Info().
		Str("elapsed", common.PrettyDuration(time.Since(start)).String()).
		Msg("Generated ethash verification cache")

	return buffer, nil
}

// generateCache creates a verification cache of a given size for an input seed.
// The cache production process involves first sequentially filling up 32 MB of
// memory, then performing two passes of Sergio Demian Lerner's RandMemoHash
// algorithm from Strict Memory Hard Hashing Functions (2014). The output is a
// set of 524288 64-byte values.
func generateCache(cache []byte, seed []byte) {
	// Calculate the number of thoretical rows (we'll store in one buffer nonetheless)
	size := uint64(len(cache))
	rows := int(size) / hashBytes

	// Create a hasher to reuse between invocations
	keccak512 := makeHasher(sha3.NewLegacyKeccak512())

	// Sequentially produce the initial dataset
	keccak512(cache, seed)
	for offset := uint64(hashBytes); offset < size; offset += hashBytes {
		keccak512(cache[offset:], cache[offset-hashBytes:offset])
	}

	// Use a low-round version of randmemohash
	temp := make([]byte, hashBytes)

	for i := 0; i < cacheRounds; i++ {
		for j := 0; j < rows; j++ {
			srcOff := ((j - 1 + rows) % rows) * hashBytes
			dstOff := j * hashBytes
			xorOff := (binary.LittleEndian.Uint32(cache[dstOff:]) % uint32(rows)) * hashBytes

			bitutil.XORBytes(temp, cache[srcOff:srcOff+hashBytes], cache[xorOff:xorOff+hashBytes])
			keccak512(cache[dstOff:], temp)
		}
	}

}

func seedHash(epoch uint64) []byte {
	block := epoch*epochLength + 1 // todo use only epoch
	seed := make([]byte, 32)
	if block < epochLength {
		return seed
	}

	keccak256 := makeHasher(sha3.NewLegacyKeccak256())
	for i := 0; i < int(block/epochLength); i++ {
		keccak256(seed, seed)
	}

	return seed
}
