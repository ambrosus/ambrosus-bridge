package middlewares

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/rs/zerolog"
)

const TgLevel zerolog.Level = 99

var DefaultOnlyLevels = []string{
	zerolog.LevelErrorValue,
	zerolog.LevelWarnValue,
	zerolog.LevelFatalValue,
	TgLevel.String(),
}

type OnlyLevels struct {
	Levels []string
	Hook   logger.Hook
}

func NewOnlyLevels(hook logger.Hook, levels []string) *OnlyLevels {
	return &OnlyLevels{
		Levels: levels,
		Hook:   hook,
	}
}

func (t *OnlyLevels) Log(l *logger.ExtLog) {
	for _, level := range t.Levels {
		if l.Level == level {
			t.Hook.Log(l)
		}
	}
}
