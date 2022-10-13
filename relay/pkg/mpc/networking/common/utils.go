package common

import (
	"context"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
)

var (
	KeygenOperation    = []byte("keygen")
	HeaderTssID        = "X-TSS-ID"
	HeaderTssOperation = "X-TSS-Operation"
	ResultPrefix       = []byte("result")
	EndpointFullMsg    = "/tx"
)

type OpError struct {
	Type string
	Err  error
}

type OperationFunc func(ctx context.Context, inCh chan []byte, outCh chan *tss_wrap.Message) ([]byte, error)
