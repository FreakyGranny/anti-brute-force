package logger

import (
	"errors"
	"strings"

	"github.com/rs/zerolog"
)

// SetLogLevel sets logging level.
func SetLogLevel(level string) error {
	logLevel, err := parseStringLogLevel(level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(logLevel)

	return nil
}

// parseStringLogLevel returns zerolog level.
func parseStringLogLevel(level string) (zerolog.Level, error) {
	switch strings.ToLower(level) {
	case "info":
		return zerolog.InfoLevel, nil
	case "debug":
		return zerolog.DebugLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	case "warn":
		return zerolog.WarnLevel, nil
	case "warning":
		return zerolog.WarnLevel, nil
	default:
		return zerolog.ErrorLevel, errors.New("unknown log level")
	}
}
