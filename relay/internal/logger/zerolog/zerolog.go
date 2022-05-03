package zerolog

import (
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/rs/zerolog"
)

var allowedLevels = map[string]func(external_logger.ExternalLogger, string, string) error{
	zerolog.LevelErrorValue: external_logger.ExternalLogger.LogError,
	zerolog.LevelWarnValue:  external_logger.ExternalLogger.LogWarning,
}

type ExternalLogger interface {
	external_logger.ExternalLogger

	Write(p []byte) (n int, err error)
	WriteLevel(l zerolog.Level, p []byte) (n int, err error)
}

type rawLog struct {
	Level   string                 `json:"level"`
	Bridge  string                 `json:"bridge"`
	Message string                 `json:"message"`
	Error   string                 `json:"error"`
	Rest    map[string]interface{} `json:"-"`
}
