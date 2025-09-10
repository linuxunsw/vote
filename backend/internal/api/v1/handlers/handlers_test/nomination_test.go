package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/store"
)

// Assumes there will be one cookie sent back, the JWT cookie.
func intoJWTCookie(headers http.Header) string {
	sc := headers.Get("Set-Cookie")
	// Example: "SESSION=abc123; Path=/; HttpOnly; Secure"
	parts := strings.Split(sc, ";")
	kv := strings.TrimSpace(parts[0]) // "SESSION=abc123"
	return "Cookie: " + kv
}

func TestNominationSubmit(t *testing.T) {
	cfg := config.Load()
	api, mailer := NewAPI(t)

	zid := "z0000000"
	electionId := createElection(t, api, cfg.JWT, []string{
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

	cookie := intoJWTCookie(res.Header)
	resp = api.Put("/api/v1/nomination", cookie, map[string]any{
		"candidate_name":      "John Doe",
		"contact_email":       "john@example.com",
		"discord_username":    "johndoe#1234",
		"executive_roles":     []string{"president", "secretary"},
		"candidate_statement": "Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50",
		"url":                 nil,
	})

	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}

	resp = api.Get("/api/v1/nomination", cookie)
	if resp.Code != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}

	outputNom := store.Nomination{}
	_ = json.Unmarshal(resp.Body.Bytes(), &outputNom)

	nominationResp := store.Nomination{
		ElectionID:         electionId,
		CandidateZID:       zid,
		CandidateName:      "John Doe",
		ContactEmail:       "john@example.com",
		DiscordUsername:    "johndoe#1234",
		ExecutiveRoles:     []string{"president", "secretary"},
		CandidateStatement: "Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50",
		URL:                nil,
		CreatedAt: 			outputNom.CreatedAt,
		UpdatedAt: 			outputNom.UpdatedAt,
	}

	if err := compareStructs(nominationResp, outputNom); err != nil {
		t.Fatal(err)
	}
}
