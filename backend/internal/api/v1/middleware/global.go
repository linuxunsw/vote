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

// Provides a huma middleware that wraps a stdlib-compatible middleware. This means that
// this middleware is tied to the underlying req/res from humago currently
func humaGoMiddleware(middleware func(http.Handler) http.Handler) func(ctx huma.Context, next func(ctx huma.Context)) {
	return func(ctx huma.Context, next func(ctx huma.Context)) {
		// NOTE: Since httplog is "router-agnostic", it should be simple enough to modify this
		// in the case that a different router is used
		r, w := humago.Unwrap(ctx)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(ctx.Context()) // https://github.com/danielgtaylor/huma/issues/859

			newCtx := humago.NewContext(ctx.Operation(), r, w)
			next(newCtx)
		})

		middleware(handler).ServeHTTP(w, r)
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

type GlobalMiddlewareOptions struct {
	Logger          *slog.Logger
	LogFormat       *httplog.Schema
	LoggerCfg       config.LoggerConfig
	RateLimitCfg    config.RateLimitConfig
	CrossOrigin     *http.CrossOriginProtection
	RealIPAllowlist []string
}

// Appends request logger, request id, cors, csrf middleware to a top level huma.API
func AddGlobalMiddleware(api huma.API, opts GlobalMiddlewareOptions) error {
	level, err := logger.ParseLevel(opts.LoggerCfg.Level)
	if err != nil {
		return err
	}

	requestLogger := newRequestLogger(opts, level)
	rateLimiter := newRateLimiter(opts.RateLimitCfg, opts.RealIPAllowlist)

	api.UseMiddleware(humaGoMiddleware(requestLogger))
	api.UseMiddleware(humaGoMiddleware(func(next http.Handler) http.Handler {
		return opts.CrossOrigin.Handler(next)
	}))
	api.UseMiddleware(humaGoMiddleware(rateLimiter))

	return nil
}
