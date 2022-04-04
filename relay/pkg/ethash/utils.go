package ethash

import (
	"encoding/binary"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"unsafe"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash/merkle"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

const (
	hashWords    = 16
	mixBytes     = 128
	hashBytes    = 64
	loopAccesses = 64

	epochLength = 30000
	maxEpoch    = 2048
)

type hasher func(dest []byte, data []byte)

func makeHasher(h hash.Hash) hasher {
	return func(dest []byte, data []byte) {
		h.Write(data)
		h.Sum(dest[:0])
		h.Reset()
	}
}

func bytesToUint32Slice(b []byte) []uint32 {
	sh := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh.Len /= 4
	sh.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&sh))
}

func (e *Ethash) populateMerkle(epoch uint64, mt *merkle.DatasetTree) (int, int, error) {
	e.logger.Debug("Populating merkle tree")
	defer e.logger.Debug("Finish populating merkle tree")

	dag, err := e.getDag(epoch)
	if err != nil {
		return 0, 0, err
	}

	fullSize := len(dag)
	fullSizeIn128Resolution := fullSize / 128
	branchDepth := len(fmt.Sprintf("%b", fullSizeIn128Resolution-1))

	mt.RegisterStoredLevel(uint32(branchDepth), 10)

	var buf [128]byte
	for i := 0; i < fullSizeIn128Resolution; i++ {
		copy(buf[:], dag[i*128:(i+1)*128])
		mt.Insert(merkle.Word(buf), uint32(i))
	}

	mt.Finalize()
	return fullSize, branchDepth, nil
}

func (e *Ethash) getVerificationIndices(blockNumber uint64, hash common.Hash, nonce uint64) ([]uint32, error) {
	// Recompute the digest and PoW value and verify against the header
	cache, err := e.getCache(epoch(blockNumber))
	if err != nil {
		return nil, err
	}

	size := datasetSize(blockNumber)
	return hashimotoLightIndices(size, bytesToUint32Slice(cache), hash.Bytes(), nonce), nil
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
	var result []uint32
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

func fnv(a, b uint32) uint32 {
	return a*0x01000193 ^ b
}

func fnvHash(mix []uint32, data []uint32) {
	for i := 0; i < len(mix); i++ {
		mix[i] = fnv(mix[i], data[i])
	}
}

func readFile(path string) ([]byte, error) {
	// read all file at once in pre-allocated buffer (10x faster than os.ReadFile)
	fInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, fInfo.Size())
	_, err = f.Read(buf)
	return buf, err
}

func dumpToFile(path string, buf []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(buf)
	return err
}

func epoch(blockNumber uint64) uint64 {
	return blockNumber / epochLength
}

// separate mutex for each epoch
type Kmutex struct {
	mapMutex   sync.Mutex
	mapOfMutex map[uint64]*sync.Mutex
}

func NewKmutex() *Kmutex {
	return &Kmutex{mapOfMutex: map[uint64]*sync.Mutex{}}
}
func (m *Kmutex) Lock(i uint64) {
	m.mapMutex.Lock()
	defer m.mapMutex.Unlock()
	if _, ok := m.mapOfMutex[i]; !ok {
		m.mapOfMutex[i] = &sync.Mutex{}
	}
	m.mapOfMutex[i].Lock()
}
func (m *Kmutex) Unlock(i uint64) {
	m.mapMutex.Lock()
	defer m.mapMutex.Unlock()
	m.mapOfMutex[i].Unlock()
	delete(m.mapOfMutex, i)
}
