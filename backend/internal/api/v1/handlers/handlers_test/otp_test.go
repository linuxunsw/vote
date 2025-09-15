package handlers_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/config"
)

func TestOTPConsumeMember(t *testing.T) {
	cfg := config.Load()
	api, mailer := NewAPI(t)

	then := time.Now().Add(-time.Second) // one second before
	zid := "z0000000"
	_ = createElection(t, api, cfg.JWT, TestingDummyJWTAdmin, []string{
		zid,
	})

	resp := generateOTPSubmit(t, api, mailer, zid)
	if resp.Code != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
	body := models.SubmitOTPResponseBody{};
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	if err != nil {
		t.Fatalf("expected body, got error %v", err)
	}

	if body.Zid != zid {
		t.Fatalf("expected zid")
	}
	if body.IsAdmin {
		t.Fatalf("expected not to be admin")
	}
	if body.Expiry.Before(then.Add(cfg.JWT.Duration)) {
		t.Fatalf("expected expiry duration within 1 second")
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
