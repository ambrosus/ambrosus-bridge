package middlewares

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
)

type Hooks struct {
	Hooks []logger.Hook
}

func NewHooks(hooks []logger.Hook) *Hooks {
	return &Hooks{Hooks: hooks}
}

func (t *Hooks) Log(l *logger.ExtLog) {
	for _, hook := range t.Hooks {
		hook.Log(l)
	}
}
