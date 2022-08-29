package antidouble

import (
	"fmt"
	"sync"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/telegram"
	zerolog "github.com/rs/zerolog/log"
)

var nowFunc = time.Now

const (
	logExpireTime = 15 * time.Minute
	counterFormat = "%s\n\n<i>This log was repeated %d times</i>"
)

//go:generate mockgen -source=antidouble.go -destination=../mocks/antidouble.go
type TgLoggerI interface {
	Send(text string) (ids []uint64, parts []string, err error)
	EditText(msgId uint64, text string) error
}

type AntiDoubleTgLogger struct {
	tgLogger TgLoggerI

	mu   sync.Mutex
	logs map[string]LogMsg
}

func NewAntiDoubleTgLogger(tgLogger TgLoggerI) *AntiDoubleTgLogger {
	return &AntiDoubleTgLogger{tgLogger: tgLogger, logs: make(map[string]LogMsg)}
}

type LogMsg struct {
	// lastLogMsg      string
	msgParts    []string
	tgId        uint64
	sentCount   uint64
	firstSentAt time.Time
}

func (l *LogMsg) isExpired() bool {
	return l.firstSentAt.Add(logExpireTime).Before(nowFunc())
}

func (f *AntiDoubleTgLogger) Log(l *logger.ExtLog) {
	msg := telegram.BuildMessage(l)
	if msg == "" {
		return
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.clearExpiredLogs()
	log, ok := f.logs[msg]
	if ok {
		log.sentCount++
		f.logs[msg] = log // cannot assign to struct field in map, can only assign the entire struct

		// TODO: add a check for the message length and do something with that (mb just send new message with the counter)
		editedMsg := fmt.Sprintf(counterFormat, log.msgParts[len(log.msgParts)-1], log.sentCount)
		err := f.tgLogger.EditText(log.tgId, editedMsg)
		if err != nil {
			zerolog.Error().Err(err).Msg("")
		}
	} else {
		// send the message
		ids, parts, err := f.tgLogger.Send(msg)
		if err != nil {
			return
		}

		// save the msg, tgId and counter
		f.logs[msg] = LogMsg{
			msgParts:    parts,
			tgId:        ids[len(ids)-1],
			sentCount:   1,
			firstSentAt: nowFunc(),
		}
	}
}

func (f *AntiDoubleTgLogger) clearExpiredLogs() {
	for msg, log := range f.logs {
		if log.isExpired() {
			delete(f.logs, msg)
		}
	}
}
