package logger

import (
	"io"
	"os"

	ext_zerolog "github.com/ambrosus/ambrosus-bridge/relay/internal/logger/zerolog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func NewSubLogger(bridge string, extLogger ext_zerolog.ExternalLogger) zerolog.Logger {
	var writer io.Writer = os.Stderr

	if extLogger != nil {
		writer = zerolog.MultiLevelWriter(writer, extLogger)
	}

	logger := zerolog.New(writer).With().Str("bridge", bridge).Logger()
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
