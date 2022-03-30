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
	if level == zerolog.ErrorLevel {
		if err := h.extLogger.LogError(message); err != nil {
			log.Error().Err(err).Msg("error send external log")
		}
	}
}

func NewSubLogger(bridge string, extLogger external_logger.ExternalLogger) zerolog.Logger {
	return log.Hook(LoggerHook{extLogger: extLogger}).With().Str("bridge", bridge).Logger()
}
