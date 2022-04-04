package ethash

import (
	"encoding/binary"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/bitutil"
	"golang.org/x/crypto/sha3"
)

const (
	cacheRounds      = 3       // Number of rounds in cache production
	cacheGrowthBytes = 1 << 17 // Cache growth per epoch
	cacheInitBytes   = 1 << 24 // Bytes in cache at genesis
)

func (e *Ethash) getCache(epoch uint64) ([]byte, error) {
	if cache, ok := e.caches[epoch]; ok {
		e.logger.Debug("Loaded old ethash cache from cache :)")
		return cache, nil
	}

	if e.useFs() {
		// Try to load the file from disk and memory map it
		cache, err := readFile(e.pathToCache(epoch))
		if err == nil {
			e.logger.Debug("Loaded old ethash cache from disk")
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
			e.logger.Warn("Failed to save cache file", "err", err)
		}
	}

	e.caches[epoch] = cache
	return cache, nil
}

func (e *Ethash) generateCache(epoch uint64) ([]byte, error) {
	e.logger.Info("Start generating ethash cache")
	start := time.Now()
	defer func() {
		e.logger.Info("Generated ethash verification cache",
			"elapsed", common.PrettyDuration(time.Since(start)))
	}()

	seed := seedHash(epoch)
	size := cacheSize(epoch)

	buffer := make([]byte, size)
	generateCache(buffer, seed)
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

	// Start a monitoring goroutine to report progress on low end devices
	var progress uint32

	// Create a hasher to reuse between invocations
	keccak512 := makeHasher(sha3.NewLegacyKeccak512())

	// Sequentially produce the initial dataset
	keccak512(cache, seed)
	for offset := uint64(hashBytes); offset < size; offset += hashBytes {
		keccak512(cache[offset:], cache[offset-hashBytes:offset])
		atomic.AddUint32(&progress, 1)
	}
	// Use a low-round version of randmemohash
	temp := make([]byte, hashBytes)

	for i := 0; i < cacheRounds; i++ {
		for j := 0; j < rows; j++ {
			var (
				srcOff = ((j - 1 + rows) % rows) * hashBytes
				dstOff = j * hashBytes
				xorOff = (binary.LittleEndian.Uint32(cache[dstOff:]) % uint32(rows)) * hashBytes
			)
			bitutil.XORBytes(temp, cache[srcOff:srcOff+hashBytes], cache[xorOff:xorOff+hashBytes])
			keccak512(cache[dstOff:], temp)

			atomic.AddUint32(&progress, 1)
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

// todo use this
/*



	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(3 * time.Second):
				logger.Info("Generating ethash verification cache", "percentage", atomic.LoadUint32(&progress)*100/uint32(rows)/4, "elapsed", common.PrettyDuration(time.Since(start)))
			}
		}
	}()
*/
