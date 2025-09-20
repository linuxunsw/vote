package handlers_test

import (
	"encoding/json"
	"testing"

	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/store"
)

func TestNominationSubmit(t *testing.T) {
	cfg := config.Load()
	api, mailer := NewAPI(t)

	zid := "z0000000"
	electionId := createElection(t, api, cfg.JWT, TestingDummyJWTAdmin, []string{
		zid,
	})

	resp := generateOTPSubmit(t, api, mailer, zid)
	res := resp.Result()
	cookie := extractCookieHeader(res.Header)
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
		CreatedAt:          outputNom.CreatedAt,
		UpdatedAt:          outputNom.UpdatedAt,
	}

	if err := compareStructs(nominationResp, outputNom); err != nil {
		t.Fatal(err)
	}

	resp = api.Delete("/api/v1/nomination", cookie)
	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}
	resp = api.Delete("/api/v1/nomination", cookie)
	if resp.Code != 204 {
		t.Fatalf("expected 204 Not Found, got %d", resp.Code)
	}
}
