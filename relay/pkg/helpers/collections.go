package helpers

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// SortedValuesWithIndices return list from map values sorted by map keys with keys to indices map
func SortedValuesWithIndices[K constraints.Integer, V any](m map[K]V) ([]V, map[K]int) {
	values := make([]V, 0, len(m))
	keyToIndex := map[K]int{}
	for i, k := range SortedKeys(m) {
		values = append(values, m[k])
		keyToIndex[k] = i
	}
	return values, keyToIndex
}

// SortedValues return list from map values sorted by map keys
func SortedValues[K constraints.Integer, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, k := range SortedKeys(m) {
		values = append(values, m[k])
	}
	return values
}

// SortedKeys return list of sorted map keys
func SortedKeys[K constraints.Integer, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return Sorted(keys)
}

// Sorted is like sort.Ints, but generi.
func Sorted[K constraints.Integer](l []K) []K {
	sort.Slice(l, func(i, j int) bool { return l[i] < l[j] })
	return l
}

func Unique[T comparable](slice []T) []T {
	mapSet := map[T]bool{}
	for _, v := range slice {
		mapSet[v] = true
	}

	var sliceSet []T
	for v := range mapSet {
		sliceSet = append(sliceSet, v)
	}

	return sliceSet
}
