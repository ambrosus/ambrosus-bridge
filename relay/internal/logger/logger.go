package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func NewLoggerWithHook(extLogger Hook) zerolog.Logger {
	var writer io.Writer = os.Stderr

	if extLogger != nil {
		writer = zerolog.MultiLevelWriter(writer, hook{extLogger})
	}

	logger := zerolog.New(writer)
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
