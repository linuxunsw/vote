package middleware

import (
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/go-chi/httplog/v3"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/logger"
)

// Provides a huma middleware that wraps a httplog RequestLogger. This means that
// this middleware is tied to the underlying req/res from humago currently
func humaRouterWithLogging(requestLogger func(http.Handler) http.Handler) func(ctx huma.Context, next func(ctx huma.Context)) {
	return func(ctx huma.Context, next func(ctx huma.Context)) {
		// NOTE: Since httplog is "router-agnostic", it should be simple enough to modify this
		// in the case that a different router is used
		r, w := humago.Unwrap(ctx)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newCtx := humago.NewContext(ctx.Operation(), r, w)
			next(newCtx)
		})

		requestLogger(handler).ServeHTTP(w, r)
	}
}

// checks whether the debug header or the environment config is set
func isLoggingBody(cfg config.LoggerConfig) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		// the debug header allows for individual request debugging
		debugHeaderSet := r.Header.Get(cfg.Debug.Header) == cfg.Debug.HeaderValue
		return cfg.LogHTTPBody || debugHeaderSet
	}
}

// Appends request logger, request id, cors, csrf middleware to a top level huma.API
// TODO: impl cors, csrf, etc
func AddGlobalMiddleware(api huma.API, cfg config.LoggerConfig, slogger *slog.Logger, logFormat *httplog.Schema) error {
	level, err := logger.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}

	// structured request logger, has panic recovery
	requestLogger := httplog.RequestLogger(slogger, &httplog.Options{
		Level:         level,
		Schema:        logFormat,
		RecoverPanics: true,

		LogRequestBody:  isLoggingBody(cfg),
		LogResponseBody: isLoggingBody(cfg),
	})

	api.UseMiddleware(humaRouterWithLogging(requestLogger))

	return nil
}
