package finalize_service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	ts := httptest.NewServer(http.FileServer(http.Dir("./fixtures")))
	defer ts.Close()

	s := NewFinalizeService(ts.URL + "/logs.bin")

	_, err := s.GetBlockWhenFinalize(0)
	assert.Equal(t, fmt.Errorf("finalize service doesn't know about block %v", 0), err)

	assert.Equal(t, uint64(2), s.cache[1])
	assert.Equal(t, uint64(4), s.cache[3])

	finalize, _ := s.GetBlockWhenFinalize(1)
	// todo check that service doesn't fetch the file again
	assert.Equal(t, uint64(2), finalize)

	s.cache = make(map[uint64]uint64) // clear cache, service must refetch info

	finalize, _ = s.GetBlockWhenFinalize(3)
	assert.Equal(t, uint64(4), finalize)

}
