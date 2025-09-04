package models

import (
	"time"

	"github.com/alexliesenfeld/health"
)

// HealthOutput mirrors the health library's CheckerResult for Huma responses.
type HealthOutput struct {
	HealthStatus string         `json:"health_status" enum:"up,down,unknown" example:"up"`
	Info         map[string]any `json:"info,omitempty"`
	// Details uses the library's CheckResult JSON marshalling
	Details map[string]health.CheckResult `json:"details,omitempty"`
	Checked time.Time                     `json:"checked"`
}

type HealthResponse struct {
	Body HealthOutput
}
