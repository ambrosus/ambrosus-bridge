package ethash

import (
	"encoding/binary"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

const (
	datasetInitBytes   = 1 << 30 // Bytes in dataset at genesis
	datasetGrowthBytes = 1 << 23 // Dataset growth per epoch
	datasetParents     = 256     // Number of parents of each dataset element
)

func (e *Ethash) getDag(epoch uint64) ([]byte, error) {
	e.dagKLock.Lock(epoch)
	defer e.dagKLock.Unlock(epoch)

	if dag, ok := e.dags[epoch]; ok {
		e.logger.Debug("Loaded old ethash dag from cache")
		return dag, nil
	}

	if e.useFs() {
		// Try to load the file from disk and memory map it
		dag, err := readFile(e.pathToDag(epoch))
		if err == nil {
			e.logger.Debug("Loaded old ethash dag from disk")
			e.dags[epoch] = dag
			return dag, nil
		}
	}

	// Generate new
	dag, err := e.generateDag(epoch)
	if err != nil {
		return nil, err
	}

	if e.useFs() {
		if err = dumpToFile(e.pathToDag(epoch), dag); err != nil {
			e.logger.Warn("Failed to save dag file", "err", err)
		}
	}

	e.dags[epoch] = dag
	return dag, nil
}

func (e *Ethash) generateDag(epoch uint64) ([]byte, error) {
	e.logger.Info("Start generating ethash dag")

	start := time.Now()
	progress := uint64(0)
	size := datasetSize(epoch)

	cache, err := e.getCache(epoch)
	if err != nil {
		return nil, err
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		for ; progress <= size/64; <-ticker.C {
			e.logger.Info("Generating DAG in progress", "percentage", progress*hashBytes*100/size,
				"progress", progress, "size", size,
				"elapsed", common.PrettyDuration(time.Since(start)))
		}
	}()

	buffer := make([]byte, size)
	generateDataset(buffer, bytesToUint32Slice(cache), &progress)

	e.logger.Info("Generated ethash dag", "elapsed", common.PrettyDuration(time.Since(start)))
	return buffer, nil
}

// generateDataset generates the entire ethash dataset for mining.
// This method places the result into dest in machine byte order.
func generateDataset(dataset []byte, cache []uint32, progress *uint64) {
	// Generate the dataset on many goroutines since it takes a while
	threads := uint64(runtime.NumCPU())
	size := uint64(len(dataset))
	batchSize := size/(hashBytes*threads) + 1

	var pend sync.WaitGroup
	pend.Add(int(threads))

	for i := uint64(0); i < threads; i++ {
		go func(id uint64) {
			defer pend.Done()

			// Create a hasher to reuse between invocations
			keccak512 := makeHasher(sha3.NewLegacyKeccak512())

			// Calculate the data segment this thread should generate
			first := id * batchSize
			limit := first + batchSize
			if limit > size/hashBytes {
				limit = size / hashBytes
			}

			// Calculate the dataset segment
			for index := first; index < limit; index++ {
				item := generateDatasetItem(cache, uint32(index), keccak512)
				copy(dataset[index*hashBytes:], item)

				atomic.AddUint64(progress, 1)
			}
		}(i)
	}

	// Wait for all the generators to finish and return
	pend.Wait()
}

// generateDatasetItem combines data from 256 pseudorandomly selected cache nodes,
// and hashes that to compute a single dataset node.
func generateDatasetItem(cache []uint32, index uint32, keccak512 hasher) []byte {
	// Calculate the number of thoretical rows (we use one buffer nonetheless)
	rows := uint32(len(cache) / hashWords)
	//r := (index % rows) * hashWords

	// Initialize the mix
	mix := make([]byte, hashBytes)

	binary.LittleEndian.PutUint32(mix, cache[(index%rows)*hashWords]^index)
	for i := 1; i < hashWords; i++ {
		binary.LittleEndian.PutUint32(mix[i*4:], cache[(index%rows)*hashWords+uint32(i)])
	}
	keccak512(mix, mix)

	// Convert the mix to uint32s to avoid constant bit shifting
	intMix := make([]uint32, hashWords)
	for i := 0; i < len(intMix); i++ {
		intMix[i] = binary.LittleEndian.Uint32(mix[i*4:])
	}

	// fnv it with a lot of random cache nodes based on index
	for i := uint32(0); i < datasetParents; i++ {
		parent := fnv(index^i, intMix[i%16]) % rows
		fnvHash(intMix, cache[parent*hashWords:])
	}

	// Flatten the uint32 mix into a binary one and return
	for i, val := range intMix {
		binary.LittleEndian.PutUint32(mix[i*4:], val)
	}
	keccak512(mix, mix)

	return mix
}
