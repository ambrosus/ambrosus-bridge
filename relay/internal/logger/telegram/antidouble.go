package telegram

import (
	"fmt"
	"sync"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
)

const (
	counterFormat = "%s\n\n<i>This log was repeated %d times</i>"
)

type AntiDoubleTgLogger struct {
	tgLogger *TgLogger

	mu              sync.Mutex
	lastLogMsg      string
	lastLogMsgParts []string
	lastLogTgId     uint64
	logSentCount    uint64
}

func NewAntiDoubleTgLogger(tgLogger *TgLogger) *AntiDoubleTgLogger {
	return &AntiDoubleTgLogger{tgLogger: tgLogger}
}

func (f *AntiDoubleTgLogger) Log(l *logger.ExtLog) {
	msg := BuildMessage(l)
	if msg == "" {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	// if msg == last sent message, add the counter to the end
	if msg == f.lastLogMsg {
		f.logSentCount++

		// TODO: add a check for the message length and do something with that (mb just send new message with the counter)
		editedMsg := fmt.Sprintf(counterFormat, f.lastLogMsgParts[len(f.lastLogMsgParts)-1], f.logSentCount)
		_ = f.tgLogger.EditText(f.lastLogTgId, editedMsg)
	} else {
		// send the message
		ids, parts, err := f.tgLogger.Send(msg)
		if err != nil {
			return
		}

		// save the msg, tgId and counter
		f.lastLogMsg = msg
		f.lastLogMsgParts = parts
		f.lastLogTgId = ids[len(ids)-1]
		f.logSentCount = 1
	}
}
