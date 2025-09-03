package middleware

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/httplog/v3"
)

// creates a structured request logger, has panic recovery
func newRequestLogger(opts GlobalMiddlewareOptions, level slog.Level) func(http.Handler) http.Handler {
	return httplog.RequestLogger(opts.Logger, &httplog.Options{
		Level:         level,
		Schema:        opts.LogFormat,
		RecoverPanics: true,

		LogRequestBody:  isLoggingBody(opts.LoggerCfg),
		LogResponseBody: isLoggingBody(opts.LoggerCfg),
	})
}
