package handlers_test

import (
	"encoding/json"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/mailer/mock_mailer"
)

// returns cookie header
func cookiePer(t *testing.T, api humatest.TestAPI, mailer *mock_mailer.MockMailer, zid string, name string, roles []string) string {
	resp := generateOTPSubmit(t, api, mailer, zid)
	res := resp.Result()
	cookie := extractCookieHeader(res.Header)

	resp = api.Put("/api/v1/nomination", cookie, map[string]any{
		"candidate_name":      name,
		"contact_email":       "john@example.com",
		"discord_username":    "johndoe#1234",
		"executive_roles":     roles,
		"candidate_statement": "Skibbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"url":                 nil,
	})
	if resp.Code != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
	return cookie
}

func TestBallotGet(t *testing.T) {
	cfg := config.Load()
	api, mailer := NewAPI(t)

	_ = createElection(t, api, cfg.JWT, TestingDummyJWTAdmin, []string{
		"z0000000",
		"z0000001",
		"z0000002",
		"z0000003",
	})

	if code := transitionElectionState(t, api, cfg.JWT, TestingDummyJWTAdmin, "NOMINATIONS_OPEN"); code != 204 {
		t.Fatalf("expected 204 No Content, got %d", code)
	}
	_ = cookiePer(t, api, mailer, "z0000000", "John Doe", []string{"president", "secretary"})
	_ = cookiePer(t, api, mailer, "z0000001", "Joe Blow", []string{"president"})
	_ = cookiePer(t, api, mailer, "z0000002", "grieve man", []string{"grievance_officer", "treasurer"})

	if code := transitionElectionState(t, api, cfg.JWT, TestingDummyJWTAdmin, "NOMINATIONS_CLOSED"); code != 204 {
		t.Fatalf("expected 204 No Content, got %d", code)
	}
	if code := transitionElectionState(t, api, cfg.JWT, TestingDummyJWTAdmin, "VOTING_OPEN"); code != 204 {
		t.Fatalf("expected 204 No Content, got %d", code)
	}

	expectedMap := map[string](map[string]bool){
		"president":         {"John Doe": false, "Joe Blow": false},
		"secretary":         {"John Doe": false},
		"grievance_officer": {"grieve man": false},
		"treasurer":         {"grieve man": false},
	}

	resp := generateOTPSubmit(t, api, mailer, "z0000003")
	res := resp.Result()
	cookie := extractCookieHeader(res.Header)

	resp = api.Get("/api/v1/ballot", cookie)
	if resp.Code != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
	ballotResp := models.PublicBallot{}
	_ = json.Unmarshal(resp.Body.Bytes(), &ballotResp)

	if ballotResp.HasVoted {
		t.Fatalf("expected HasVoted to be false, got true")
	}
	for role, noms := range ballotResp.Candidates {
		expNoms, ok := expectedMap[role]
		if !ok {
			t.Fatalf("unexpected role %q", role)
		}
		for _, nom := range noms {
			_, ok := expNoms[nom.CandidateName]
			if !ok {
				t.Fatalf("unexpected candidate %q for role %q", nom.CandidateName, role)
			}
			expNoms[nom.CandidateName] = true
		}
		for nom, found := range expNoms {
			if !found {
				t.Fatalf("expected candidate %q for role %q not found", nom, role)
			}
		}
		delete(expectedMap, role)
	}

	// check if expectedMap is now empty
	if len(expectedMap) != 0 {
		t.Fatalf("not all expected roles found, missing: %v", expectedMap)
	}
}
