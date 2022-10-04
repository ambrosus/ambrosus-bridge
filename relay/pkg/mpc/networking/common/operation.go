package common

import "github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"

type Operation struct {
	Started bool
	// will be sent to tss
	InCh chan []byte
	// received from tss or another client; will be sent to own tss or another client
	OutCh   chan *tss_wrap.OutputMessage
	SignMsg []byte
}

const chSize = 10

func NewOperation(msg []byte) Operation {
	return Operation{
		Started: true,
		InCh:    make(chan []byte, chSize),
		OutCh:   make(chan *tss_wrap.OutputMessage, chSize),
		SignMsg: msg,
	}
}

func (o Operation) Stop() {
	o.Started = false
	close(o.InCh)
	close(o.OutCh)
}
