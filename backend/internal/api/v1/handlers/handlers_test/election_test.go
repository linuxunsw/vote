package handlers_test

import (
	"fmt"
	"testing"

	"github.com/linuxunsw/vote/backend/internal/config"
)

func TestElectionCreateFailure(t *testing.T) {
	cfg := config.Load()
	api, _ := NewAPI(t)

	zid := "z0000000"

	_ = createElection(t, api, cfg.JWT, TestingDummyJWTAdmin, []string{
		zid,
	})

	// create an election
	adminCookie := fmt.Sprintf("Cookie: %s=%s", cfg.JWT.CookieName, TestingDummyJWTAdmin)
	resp := api.Post("/api/v1/elections", adminCookie, map[string]any{
		"name": "Test Electione!",
	})
	if resp.Code != 400 {
		t.Fatalf("expected 400 Bad Request, got %d", resp.Code)
	}
}
