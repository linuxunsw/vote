package middleware

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v3"
	"github.com/linuxunsw/vote/backend/internal/config"
)

// checks whether the debug header or the environment config is set
func isLoggingBody(cfg config.LoggerConfig) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		// the debug header allows for individual request debugging
		debugHeaderSet := r.Header.Get(cfg.Debug.Header) == cfg.Debug.HeaderValue
		return cfg.LogHTTPBody || debugHeaderSet
	}
}

// creates a structured request logger, has panic recovery
func newRequestLogger(opts GlobalMiddlewareOptions, level slog.Level) func(http.Handler) http.Handler {
	return httplog.RequestLogger(opts.Logger, &httplog.Options{
		Level:              level,
		Schema:             opts.LogFormat,
		RecoverPanics:      true,
		LogResponseHeaders: []string{opts.LoggerCfg.RequestIDHeader},

		LogRequestBody:  isLoggingBody(opts.LoggerCfg),
		LogResponseBody: isLoggingBody(opts.LoggerCfg),
	})
}
