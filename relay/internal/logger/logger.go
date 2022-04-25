package logger

import (
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func NewSubLogger(bridge string, extLogger external_logger.ExternalLogger) zerolog.Logger {
	logger := log.With().Str("bridge", bridge).Logger()
	if extLogger != nil {
		logger = logger.Hook(Hook{extLogger})
	}
	return logger
}

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if os.Getenv("DEBUG") != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
