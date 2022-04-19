package logger

import (
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerHook struct {
	extLogger external_logger.ExternalLogger
}

func (h LoggerHook) Run(event *zerolog.Event, level zerolog.Level, message string) {
	var logFunc func(string) error

	switch level {
	case zerolog.ErrorLevel:
		logFunc = h.extLogger.LogError
	case zerolog.WarnLevel:
		logFunc = h.extLogger.LogWarning
	default:
		return
	}

	if err := logFunc(message); err != nil {
		log.Error().Err(err).Msg("error send external log")
	}
}

func NewSubLogger(bridge string, extLogger external_logger.ExternalLogger) zerolog.Logger {
	return log.Hook(LoggerHook{extLogger: extLogger}).With().Str("bridge", bridge).Logger()
}
