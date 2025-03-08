package logger

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	cfg "github.com/spf13/pflag"
)

var (
	levelStr = cfg.String("log.level", "info", "Logging level")
)

// nolint: gochecknoinits
func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = marshalStack

	writer := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.TimeOnly}

	log.Logger = log.Output(writer).With().Caller().Stack().Logger()
}

func Setup() error {
	level, err := zerolog.ParseLevel(*levelStr)
	if err != nil {
		return errors.Wrapf(err, "invalid log level(%s)", *levelStr)
	}

	log.Logger = log.Logger.Level(level)

	return nil
}
