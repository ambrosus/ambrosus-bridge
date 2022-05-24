package telegram

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/rs/zerolog"
)

const fieldsFormat = "<b>%s</b>: %s\n"

func (t *tgLogger) Log(l *logger.ExtLog) {
	var msg string

	msg += fmt.Sprintf("[<b>%s</b>] ", strings.ToUpper(l.Bridge))

	switch l.Level {
	case zerolog.LevelErrorValue:
		msg += fmt.Sprintf("<b>%s</b>\n", "We got an unexpected error:")
	case zerolog.LevelWarnValue:
		msg += fmt.Sprintf("<b>%s</b>\n", "Warning!")
	default:
		return
	}

	var fields = make([]string, 0, len(l.Rest))
	for field := range l.Rest {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	// escape telegram markup symbols
	replacer := strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;")

	for _, field := range fields {
		msg += fmt.Sprintf(fieldsFormat, field, replacer.Replace(l.Rest[field].(string)))
	}
	if l.Error != "" {
		msg += fmt.Sprintf(fieldsFormat, "error", replacer.Replace(l.Error))
	}
	if l.Message != "" {
		msg += fmt.Sprintf("\n%s", replacer.Replace(l.Message))
	}

	_ = t.send(msg)
}
