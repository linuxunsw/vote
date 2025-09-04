package handlers

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
)

// NewChecker builds and starts a health.Checker. db can be nil (no db check).
func NewChecker(log *slog.Logger, db *sql.DB) health.Checker {
	opts := []health.CheckerOption{
		health.WithCacheDuration(1 * time.Second),
		health.WithTimeout(5 * time.Second),
		health.WithCheck(health.Check{
			Name:  "app",
			Check: func(ctx context.Context) error { return nil },
		}),
		health.WithStatusListener(func(ctx context.Context, state health.CheckerState) {
			log.Info("Completed health check",
				"status", state.Status,
				"state checks", state.CheckState,
			)
		}),
	}

	if db != nil {
		opts = append(opts, health.WithCheck(health.Check{
			Name:    "database",
			Timeout: 2 * time.Second,
			Check:   db.PingContext,
		}))
	}

	checker := health.NewChecker(opts...)
	// health.NewChecker auto-starts by default. Return it so caller can Stop() on shutdown.
	return checker
}

// Huma health check handler
// It expects a started health.Checker to be provided via closure when registering.
func GetHealth(checker health.Checker) func(ctx context.Context, _ *struct{}) (*models.HealthResponse, error) {
	return func(ctx context.Context, _ *struct{}) (*models.HealthResponse, error) {
		res := checker.Check(ctx) // runs synchronous checks and returns CheckerResult

		out := &models.HealthOutput{
			HealthStatus: string(res.Status),
			Info:         res.Info,
			Details:      res.Details,
			Checked:      time.Now().UTC(),
		}

		response := &models.HealthResponse{Body: *out}
		return response, nil
	}
}
