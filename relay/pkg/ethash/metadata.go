package ethash

import (
	"encoding/binary"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"golang.org/x/crypto/sha3"
)

type EthashMetaData struct {
	cachedir     string // Data directory to store the verification caches
	cachesinmem  int    // Number of caches to keep in memory
	cachesondisk int    // Number of caches to keep on disk
	dagdir       string // Data directory to store full mining datasets
	dagsinmem    int    // Number of mining datasets to keep in memory
	dagsondisk   int    // Number of mining datasets to keep on disk

	caches   map[uint64]*cache   // In memory caches to avoid regenerating too often
	fcache   *cache              // Pre-generated cache for the estimated future epoch
	datasets map[uint64]*dataset // In memory datasets to avoid regenerating too often

	// Mining related fields
	update   chan struct{} // Notification channel to update mining parameters
	hashrate metrics.Meter // Meter tracking the average hashrate

	// The fields below are hooks for testing
	tester bool // Flag whether to use a smaller test dataset

	lock sync.Mutex // Ensures thread safety for the in-memory caches and mining fields
}

func New(cachedir string, cachesinmem, cachesondisk int, dagdir string, dagsinmem, dagsondisk int) *EthashMetaData {
	if cachesinmem <= 0 {
		log.Warn("One ethash cache must alwast be in memory", "requested", cachesinmem)
		cachesinmem = 1
	}
	if cachedir != "" && cachesondisk > 0 {
		log.Info("Disk storage enabled for ethash caches", "dir", cachedir, "count", cachesondisk)
	}
	if dagdir != "" && dagsondisk > 0 {
		log.Info("Disk storage enabled for ethash DAGs", "dir", dagdir, "count", dagsondisk)
	}
	return &EthashMetaData{
		cachedir:     cachedir,
		cachesinmem:  cachesinmem,
		cachesondisk: cachesondisk,
		dagdir:       dagdir,
		dagsinmem:    dagsinmem,
		dagsondisk:   dagsondisk,
		caches:       make(map[uint64]*cache),
		datasets:     make(map[uint64]*dataset),
		update:       make(chan struct{}),
		hashrate:     metrics.NewMeter(),
	}
}

func (ethash *EthashMetaData) GetVerificationIndices(blockNumber uint64, hash common.Hash, nonce uint64) []uint32 {
	// Recompute the digest and PoW value and verify against the header
	cache := ethash.cache(blockNumber)

	size := DatasetSize(blockNumber)
	return hashimotoLightIndices(size, cache, hash.Bytes(), nonce)
}

func hashimotoLightIndices(size uint64, cache []uint32, hash []byte, nonce uint64) []uint32 {
	keccak512 := makeHasher(sha3.NewLegacyKeccak512())

	lookup := func(index uint32) []uint32 {
		rawData := generateDatasetItem(cache, index, keccak512)

		data := make([]uint32, len(rawData)/4)
		for i := 0; i < len(data); i++ {
			data[i] = binary.LittleEndian.Uint32(rawData[i*4:])
		}
		return data
	}
	return hashimotoIndices(hash, nonce, size, lookup)
}

func hashimotoIndices(hash []byte, nonce uint64, size uint64, lookup func(index uint32) []uint32) []uint32 {
	result := []uint32{}
	// Calculate the number of thoretical rows (we use one buffer nonetheless)
	rows := uint32(size / mixBytes)

	// Combine header+nonce into a 64 byte seed
	seed := make([]byte, 40)
	copy(seed, hash)
	binary.LittleEndian.PutUint64(seed[32:], nonce)

	seed = crypto.Keccak512(seed)
	seedHead := binary.LittleEndian.Uint32(seed)

	// Start the mix with replicated seed
	mix := make([]uint32, mixBytes/4)
	for i := 0; i < len(mix); i++ {
		mix[i] = binary.LittleEndian.Uint32(seed[i%16*4:])
	}
	// Mix in random dataset nodes
	temp := make([]uint32, len(mix))

	for i := 0; i < loopAccesses; i++ {
		parent := fnv(uint32(i)^seedHead, mix[i%len(mix)]) % rows
		result = append(result, parent)
		for j := uint32(0); j < mixBytes/hashBytes; j++ {
			copy(temp[j*hashWords:], lookup(2*parent+j))
		}

		fnvHash(mix, temp)
	}

	return result
}

func fnvHash(mix []uint32, data []uint32) {
	for i := 0; i < len(mix); i++ {
		mix[i] = mix[i]*0x01000193 ^ data[i]
	}
}

func (ethash *EthashMetaData) cache(block uint64) []uint32 {
	epoch := block / epochLength

	// If we have a PoW for that epoch, use that
	ethash.lock.Lock()

	current, future := ethash.caches[epoch], (*cache)(nil)
	if current == nil {
		// No in-memory cache, evict the oldest if the cache limit was reached
		for len(ethash.caches) > 0 && len(ethash.caches) >= ethash.cachesinmem {
			var evict *cache
			for _, cache := range ethash.caches {
				if evict == nil || evict.used.After(cache.used) {
					evict = cache
				}
			}
			delete(ethash.caches, evict.epoch)
			evict.release()

			log.Trace("Evicted ethash cache", "epoch", evict.epoch, "used", evict.used)
		}
		// If we have the new cache pre-generated, use that, otherwise create a new one
		if ethash.fcache != nil && ethash.fcache.epoch == epoch {
			log.Trace("Using pre-generated cache", "epoch", epoch)
			current, ethash.fcache = ethash.fcache, nil
		} else {
			log.Trace("Requiring new ethash cache", "epoch", epoch)
			current = &cache{epoch: epoch}
		}
		ethash.caches[epoch] = current

		// If we just used up the future cache, or need a refresh, regenerate
		if ethash.fcache == nil || ethash.fcache.epoch <= epoch {
			if ethash.fcache != nil {
				ethash.fcache.release()
			}
			log.Trace("Requiring new future ethash cache", "epoch", epoch+1)
			future = &cache{epoch: epoch + 1}
			ethash.fcache = future
		}
		// New current cache, set its initial timestamp
		current.used = time.Now()
	}
	ethash.lock.Unlock()

	// Wait for generation finish, bump the timestamp and finalize the cache
	current.generate(ethash.cachedir, ethash.cachesondisk, ethash.tester)

	current.lock.Lock()
	current.used = time.Now()
	current.lock.Unlock()

	// If we exhausted the future cache, now's a good time to regenerate it
	if future != nil {
		go future.generate(ethash.cachedir, ethash.cachesondisk, ethash.tester)
	}
	return current.cache
}
