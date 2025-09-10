package handlers_test

import (
	"testing"

	"github.com/linuxunsw/vote/backend/internal/config"
)

func TestOTPConsumeOnce(t *testing.T) {
	cfg := config.Load()
	api, mailer := NewAPI(t)

	zid := "z0000000"
	_ = createElection(t, api, cfg.JWT, []string{
		zid,
	})

	resp := api.Post("/api/v1/otp/generate", map[string]any{
		"zid": zid,
	})
	// generate returns nothing, it goes to the mailer
	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}

	code := mailer.MockRetrieveOTP(zid + "@unsw.edu.au")

	resp = api.Post("/api/v1/otp/submit", map[string]any{
		"zid": zid,
		"otp": code,
	})
	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}

	res := resp.Result()
	if res.Header.Get("Set-Cookie") == "" {
		t.Fatalf("expected Set-Cookie header, got none")
	}
}
