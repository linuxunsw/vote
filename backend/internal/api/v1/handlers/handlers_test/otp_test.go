package handlers_test

import (
	"testing"

	"github.com/linuxunsw/vote/backend/internal/config"
)

func TestOTPConsumeMember(t *testing.T) {
	cfg := config.Load()
	api, mailer := NewAPI(t)

	zid := "z0000000"
	_ = createElection(t, api, cfg.JWT, TestingDummyJWTAdmin, []string{
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

	_ = createElection(t, api, cfg.JWT, TestingDummyJWTAdmin, []string{})

	zid := "z0000000"
	resp := generateOTPSubmit(t, api, mailer, zid)
	if resp.Code != 403 {
		t.Fatalf("expected 403 Forbidden, got %d", resp.Code)
	}
}

func TestOTPAdmins(t *testing.T) {
	t.Setenv("ADMIN_ZIDS", "z0000000,z0000001,z0000002")

	zid := "z0000001"

	cfg := config.Load()
	api, mailer := NewAPI(t)

	resp := generateOTPSubmit(t, api, mailer, zid)
	res := resp.Result()

	if res.Header.Get("Set-Cookie") == "" {
		t.Fatalf("expected Set-Cookie header, got none")
	}
	jwt := extractJWT(res.Header)
	_ = createElection(t, api, cfg.JWT, jwt, []string{})
}
