package middlewares

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram"
)

const (
	counterFormat = "%s\n\n<i>This log was repeated %d times</i>"
)

type AntiDoubleMiddleware struct {
	tgLogger *telegram.TgLogger

	lastLogMsg      string
	lastLogMsgParts []string
	lastLogTgId     uint64
	logSentCount    uint64
}

func NewAntiDoubleMiddleware(tgLogger *telegram.TgLogger) *AntiDoubleMiddleware {
	return &AntiDoubleMiddleware{tgLogger: tgLogger}
}

func (f *AntiDoubleMiddleware) Log(l *logger.ExtLog) {
	msg := telegram.BuildMessage(l)

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
