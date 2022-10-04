package common

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
)

// todo move to another file? or rename this file? or pohuy?
var (
	KeygenOperation    = []byte("keygen")
	HeaderTssID        = "X-TSS-ID"
	HeaderTssOperation = "X-TSS-Operation"
)

type Operation struct {
	Started bool
	// will be sent to tss
	InCh chan []byte
	// received from tss or another client; will be sent to own tss or another client
	OutCh   chan *tss_wrap.OutputMessage
	SignMsg []byte
}

const chSize = 10

func NewOperation() Operation {
	return Operation{
		Started: false,
		InCh:    make(chan []byte, chSize),
		OutCh:   make(chan *tss_wrap.OutputMessage, chSize),
	}
}

func (o *Operation) Start(msg []byte) error {
	if o.Started {
		return fmt.Errorf("already started")
	}
	o.SignMsg = msg
	o.Started = true
	return nil
}

func (o *Operation) Stop() {
	o.Started = false
}
