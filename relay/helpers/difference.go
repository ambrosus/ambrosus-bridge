package helpers

import "github.com/ethereum/go-ethereum/common"

// DiffAddresses returns the elements in `a` that aren't in `b` for common.Address
func DiffAddresses(a, b []common.Address) []common.Address {
	mb := make(map[common.Address]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []common.Address
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
