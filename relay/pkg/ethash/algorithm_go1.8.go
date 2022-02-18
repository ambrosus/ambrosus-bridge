package ethash

import "math/big"

func cacheSize(block uint64) uint64 {
	// If we have a pre-generated value, use that
	epoch := int(block / epochLength)
	if epoch < len(cacheSizes) {
		return cacheSizes[epoch]
	}
	// No known cache size, calculate manually (sanity branch only)
	size := uint64(cacheInitBytes + cacheGrowthBytes*uint64(epoch) - hashBytes)
	for !new(big.Int).SetUint64(size / hashBytes).ProbablyPrime(1) { // Always accurate for n < 2^64
		size -= 2 * hashBytes
	}
	return size
}

func datasetSize(block uint64) uint64 {
	// If we have a pre-generated value, use that
	epoch := int(block / epochLength)
	if epoch < len(datasetSizes) {
		return datasetSizes[epoch]
	}
	// No known dataset size, calculate manually (sanity branch only)
	size := uint64(datasetInitBytes + datasetGrowthBytes*uint64(epoch) - mixBytes)
	for !new(big.Int).SetUint64(size / mixBytes).ProbablyPrime(1) { // Always accurate for n < 2^64
		size -= 2 * mixBytes
	}
	return size
}
