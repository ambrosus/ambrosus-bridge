package middlewares

import (
	"strings"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
)

var DefaultSkipSubstrings = []string{
	"websocket: close 1006 (abnormal closure): unexpected EOF",
	"websocket: close 1011 (internal server error)",
	"connection reset by peer",
	"cannot retry err [http2: Transport received Server's graceful shutdown GOAWAY] after Request.Body was written; define Request.GetBody to avoid this error",
	"i/o timeout",
	"CloudFlare WebSocket proxy restarting",
}

type SkipSubstrings struct {
	Substrings []string
	Hook       logger.Hook
}

func NewSkipSubstrings(hook logger.Hook, substrings []string) *SkipSubstrings {
	return &SkipSubstrings{
		Substrings: substrings,
		Hook:       hook,
	}
}

func (t *SkipSubstrings) Log(l *logger.ExtLog) {
	for _, s := range t.Substrings {
		if strings.Contains(l.Error, s) {
			return
		}
	}
	t.Hook.Log(l)
}
