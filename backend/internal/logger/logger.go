package logger

import (
	"log/slog"
	"os"

	"github.com/golang-cz/devslog"
	"github.com/linuxunsw/vote/backend/internal/config"
)

func ParseLevel(s string) (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(s))
	return level, err
}

// Creates a new log/slog logger, using devslog in dev environments for pretty printing.
// Accepts config as well as initial handler options, which can be used to provide ReplaceAttrs.
func New(cfg config.LoggerConfig, opts *slog.HandlerOptions) (*slog.Logger, error) {
	parsedLevel, err := ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	opts.Level = parsedLevel
	opts.AddSource = cfg.AddSource

	var h slog.Handler
	if cfg.PrettyPrint {
		// devslog provides colorful text handler implementing slog.Handler
		h = devslog.NewHandler(os.Stdout, &devslog.Options{
			SortKeys:           true,
			MaxErrorStackTrace: 5,
			MaxSlicePrintSize:  20,
			HandlerOptions:     opts,
		})
	} else {
		h = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(h)

	return logger, nil
}
