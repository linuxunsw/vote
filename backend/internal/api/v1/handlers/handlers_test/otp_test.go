package handlers_test

import (
	"testing"

	"github.com/linuxunsw/vote/backend/internal/config"
)

func TestOTPConsumeMember(t *testing.T) {
	cfg := config.Load()
	api, mailer := NewAPI(t)

	zid := "z0000000"
	_ = createElection(t, api, cfg.JWT, []string{
		zid,
	})

	resp := generateOTPSubmit(t, api, mailer, zid)
	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}

	res := resp.Result()
	if res.Header.Get("Set-Cookie") == "" {
		t.Fatalf("expected Set-Cookie header, got none")
	}
}

func TestOTPConsumeNonMember(t *testing.T) {
	cfg := config.Load()
	api, mailer := NewAPI(t)

	_ = createElection(t, api, cfg.JWT, []string{})

	zid := "z0000000"
	resp := generateOTPSubmit(t, api, mailer, zid)
	if resp.Code != 403 {
		t.Fatalf("expected 403 Forbidden, got %d", resp.Code)
	}
}
