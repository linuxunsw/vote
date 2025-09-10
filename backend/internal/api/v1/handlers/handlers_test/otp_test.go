package handlers_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/config"
)

const (
	TestingDummyJWTAdmin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ6MTIzNDU2NyIsImlzcyI6InZvdGUtYXBpIiwiaXNBZG1pbiI6dHJ1ZX0.pzeRspChlGtMEc3XuVLjDpFriJzdehXqny7L85VxSW0"
)

func TestOTPConsumeOnce(t *testing.T) {
	cfg := config.Load()
	api, mailer := NewAPI(t)

	// TODO remove TestingDummyJWTAdmin when auth implemented
	adminCookie := fmt.Sprintf("Cookie: %s=%s", cfg.JWT.CookieName, TestingDummyJWTAdmin)

	// create an election
	resp := api.Post("/api/v1/elections", adminCookie, map[string]any{
		"name": "Test Election",
	})
	if resp.Code != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}

	electionResp := models.CreateElectionResponseBody{}
	err := json.Unmarshal(resp.Body.Bytes(), &electionResp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	zid := "z0000000"

	// put member list
	resp = api.Put("/api/v1/elections/"+electionResp.ElectionId+"/members", adminCookie, map[string]any{
		"zids": []string{zid},
	})

	// "No Content" with PUT
	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}

	resp = api.Post("/api/v1/otp/generate", map[string]any{
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
