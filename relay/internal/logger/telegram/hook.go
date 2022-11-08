package telegram

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/rs/zerolog"
)

const fieldsFormat = "<b>%s</b>: %s\n"

func (t *TgLogger) Log(l *logger.ExtLog) {
	msg := BuildMessage(l)
	if msg == "" {
		return
	}

	_, _, _ = t.Send(msg)
}

func BuildMessage(l *logger.ExtLog) string {
	var msg string

	if l.Bridge != "" && l.SideBridge != "" {
		msg += fmt.Sprintf("[<b>%s->%s</b>] ", strings.ToUpper(l.Bridge), strings.ToUpper(l.SideBridge))
	}

	switch l.Level {
	case zerolog.LevelErrorValue:
		msg += fmt.Sprintf("<b>%s</b>\n", "We got an unexpected error:")
	case zerolog.LevelWarnValue:
		msg += fmt.Sprintf("<b>%s</b>\n", "Warning!")
	case zerolog.LevelFatalValue:
		msg += fmt.Sprintf("<b>%s</b>\n", "Fatal!")
	default:
		msg += fmt.Sprintf("<b>%s</b>\n", "Log")
	}

	var fields = make([]string, 0, len(l.Rest))
	for field := range l.Rest {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	// escape telegram markup symbols
	replacer := strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")

	if l.Service != "" {
		msg += fmt.Sprintf(fieldsFormat, "service", replacer.Replace(l.Service))
	}
	for _, field := range fields {
		msg += fmt.Sprintf(fieldsFormat, field, replacer.Replace(fmt.Sprintf("%v", l.Rest[field])))
	}
	if l.Error != "" {
		msg += fmt.Sprintf(fieldsFormat, "error", replacer.Replace(l.Error))
	}
	if l.Message != "" {
		msg += fmt.Sprintf("\n%s", replacer.Replace(l.Message))
	}

	return msg
}
