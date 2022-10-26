//go:build !ci

package fixtures

import (
	"testing"
)

// Generate pre params fixtures
func TestGeneratePreParams(t *testing.T) {
	GetOrGenPreParams("")
}
