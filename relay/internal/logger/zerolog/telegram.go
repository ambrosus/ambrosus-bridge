package zerolog

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/rs/zerolog"
)

var fieldsFormat = "<b>%s</b>: %s\n"

type tgLogger struct {
	external_logger.ExternalLogger
}

func NewTgLogger(logger external_logger.ExternalLogger) *tgLogger {
	return &tgLogger{logger}
}

func (t *tgLogger) Write(p []byte) (n int, err error) {
	var rawLog rawLog

	if err := t.parseLog(p, &rawLog); err != nil {
		return len(p), err
	}

	prefix := fmt.Sprintf("[<b>%s</b>]", strings.ToUpper(rawLog.Bridge))
	msg := t.makeMsg(&rawLog)

	return len(p), allowedLevels[rawLog.Level](t, prefix, msg)
}

func (t *tgLogger) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if _, ok := allowedLevels[level.String()]; ok {
		return t.Write(p)
	}
	return len(p), nil
}

func (t *tgLogger) parseLog(p []byte, rawLog *rawLog) error {
	if err := json.Unmarshal(p, rawLog); err != nil {
		return err
	}
	if err := json.Unmarshal(p, &rawLog.Rest); err != nil {
		return err
	}
	delete(rawLog.Rest, "level")
	delete(rawLog.Rest, "bridge")
	delete(rawLog.Rest, "message")
	delete(rawLog.Rest, "error")
	return nil
}

func (t *tgLogger) makeMsg(rawLog *rawLog) string {
	var msg string
	var fields = make([]string, 0, len(rawLog.Rest))

	for field := range rawLog.Rest {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	for _, field := range fields {
		msg += fmt.Sprintf(fieldsFormat, field, rawLog.Rest[field])
	}
	if rawLog.Error != "" {
		msg += fmt.Sprintf(fieldsFormat, "error", rawLog.Error)
	}
	if rawLog.Message != "" {
		msg += fmt.Sprintf("\n%s", rawLog.Message)
	}
	return msg
}
