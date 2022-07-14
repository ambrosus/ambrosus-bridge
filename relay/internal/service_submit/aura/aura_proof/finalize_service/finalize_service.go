package finalize_service

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
)

type FinalizeService struct {
	fileUrl string
	cache   map[uint64]uint64 // signaled => finalize
}

func NewFinalizeService(url string) *FinalizeService {
	return &FinalizeService{
		fileUrl: url,
		cache:   make(map[uint64]uint64),
	}
}

func (s *FinalizeService) GetBlockWhenFinalize(signalBlockNum uint64) (uint64, error) {
	if finalizeAt, ok := s.cache[signalBlockNum]; ok {
		return finalizeAt, nil
	}

	// signalBlockNum not in cache, update cache
	if err := s.getFinalized(); err != nil {
		return 0, fmt.Errorf("getFinalized: %w", err)
	}
	if finalizeAt, ok := s.cache[signalBlockNum]; ok {
		return finalizeAt, nil
	}

	// signalBlockNum still not in cache, it's kinda sus
	return 0, fmt.Errorf("finalize service doesn't know about block %v", signalBlockNum)
}

func (s *FinalizeService) getFinalized() error {
	resp, err := http.Get(s.fileUrl)
	if err != nil {
		return fmt.Errorf("http.Get: %w", err)
	}
	defer resp.Body.Close()

	newCache, err := parseFinalizeFile(resp.Body)
	if err != nil {
		return fmt.Errorf("parseFinalizeFile: %w", err)
	}
	s.cache = newCache
	return nil
}

func parseFinalizeFile(file io.Reader) (map[uint64]uint64, error) {
	signaledAt := uint64(0)
	finalizeAt := uint64(0)
	m := make(map[uint64]uint64)

	for {
		if err := binary.Read(file, binary.LittleEndian, &signaledAt); err == io.EOF {
			break // file ends, returning map
		} else if err != nil {
			return nil, fmt.Errorf("read signaledAt num: %w", err)
		}
		if err := binary.Read(file, binary.LittleEndian, &finalizeAt); err == io.EOF {
			break // file ends, returning map
		} else if err != nil {
			return nil, fmt.Errorf("read finalizeAt num: %w", err)
		}
		m[signaledAt] = finalizeAt
	}

	return m, nil
}
