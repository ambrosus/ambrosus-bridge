package helpers

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// SortedKeys used for 'ordered' map
func SortedKeys[K constraints.Integer, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
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
