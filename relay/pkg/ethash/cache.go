package ethash

import (
	"encoding/binary"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common/bitutil"
	"golang.org/x/crypto/sha3"
)

const (
	cacheRounds      = 3       // Number of rounds in cache production
	cacheGrowthBytes = 1 << 17 // Cache growth per epoch
	cacheInitBytes   = 1 << 24 // Bytes in cache at genesis
)

func (e *Ethash) getOrGenerateCache(epoch uint64) ([]byte, error) {
	path := e.pathToCache(epoch)

	// Try to load the file from disk and memory map it
	cache, err := readFile(path)
	if err == nil {
		e.logger.Debug("Loaded old ethash cache from disk")
		return cache, nil
	}
	e.logger.Debug("Failed to load old cache cache", "err", err)

	// Generate new
	cache, err = e.generateCache(epoch)
	if err != nil {
		return nil, err
	}
	err = dumpToFile(path, cache)
	return cache, err
}

func (e *Ethash) generateCache(epoch uint64) ([]byte, error) {
	e.logger.Info("Start generating ethash cache")

	seed := seedHash(epoch*epochLength + 1)
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
	size := uint64(cap(cache))
	rows := int(size) / hashBytes

	// Start a monitoring goroutine to report progress on low end devices
	var progress uint32

	done := make(chan struct{})
	defer close(done)

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

// todo use this
/*

start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		logger.Info("Generated ethash verification cache", "elapsed", common.PrettyDuration(elapsed))
	}()


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
