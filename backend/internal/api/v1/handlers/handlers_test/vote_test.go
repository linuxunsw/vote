package handlers_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/config"
)

func TestVoteSubmit(t *testing.T) {
	nowProvider := new(func() time.Time)
	*nowProvider = func() time.Time { return time.Now() } // dummy

	cfg := config.Load()
	api, mailer := NewAPIWithNowProvider(t, func() time.Time {
		return (*nowProvider)()
	})

	zid0 := "z0000000"
	zid1 := "z0000001"
	_ = createElection(t, api, cfg.JWT, TestingDummyJWTAdmin, []string{
		zid0,
		zid1,
	})

	if code := transitionElectionState(t, api, cfg.JWT, TestingDummyJWTAdmin, "NOMINATIONS_OPEN"); code != 204 {
		t.Fatalf("expected 204 No Content, got %d", code)
	}

	res0 := generateOTPSubmit(t, api, mailer, zid0).Result()
	cookie0 := extractCookieHeader(res0.Header)
	resp := api.Put("/api/v1/nomination", cookie0, map[string]any{
		"candidate_name":      "John Doe",
		"contact_email":       "john@example.com",
		"discord_username":    "johndoe#1234",
		"executive_roles":     []string{"president", "secretary"},
		"candidate_statement": "Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50Deez50",
		"url":                 nil,
	})
	if resp.Code != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
	body0 := models.SubmitNominationResponseBody{}
	_ = json.Unmarshal(resp.Body.Bytes(), &body0)
	nominationId0 := body0.NominationID

	res1 := generateOTPSubmit(t, api, mailer, zid1).Result()
	cookie1 := extractCookieHeader(res1.Header)

	now0 := time.Now().Truncate(time.Second)
	*nowProvider = func() time.Time { return now0 }

	if code := transitionElectionState(t, api, cfg.JWT, TestingDummyJWTAdmin, "NOMINATIONS_CLOSED"); code != 204 {
		t.Fatalf("expected 204 No Content, got %d", code)
	}
	if code := transitionElectionState(t, api, cfg.JWT, TestingDummyJWTAdmin, "VOTING_OPEN"); code != 204 {
		t.Fatalf("expected 204 No Content, got %d", code)
	}

	resp = api.Put("/api/v1/vote", cookie1, map[string]any{
		"positions": map[string]string{
			"president": nominationId0,
			"secretary": nominationId0,
		},
	})
	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}

	// test the time not being updated
	now1 := time.Now().Truncate(time.Second)
	*nowProvider = func() time.Time { return now1 }

	resp = api.Put("/api/v1/vote", cookie1, map[string]any{
		"positions": map[string]string{
			"president": nominationId0,
			"bogus":     nominationId0,
		},
	})
	// HTTP/1.1 422 Unprocessable Entity
	if resp.Code != 422 {
		// bogus
		t.Fatalf("expected 422 Unprocessable Entity, got %d", resp.Code)
	}

	resp = api.Put("/api/v1/vote", cookie1, map[string]any{
		"positions": map[string]string{
			"president":         nominationId0,
			"grievance_officer": nominationId0,
		},
	})
	// HTTP/1.1 400 Bad Request
	if resp.Code != 400 {
		// grievance_officer not open for election
		t.Fatalf("expected 400 Bad Request, got %d", resp.Code)
	}

	resp = api.Get("/api/v1/vote", cookie1)
	if resp.Code != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
	body := models.Vote{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	if err != nil {
		t.Fatalf("expected body, got error %v", err)
	}

	body.CreatedAt = normaliseTime(body.CreatedAt)
	body.UpdatedAt = normaliseTime(body.UpdatedAt)

	unchangedResp := models.Vote{
		Positions: map[string]string{
			"president": nominationId0,
			"secretary": nominationId0,
		},
		CreatedAt: normaliseTime(now0),
		UpdatedAt: normaliseTime(now0),
	}
	if err := compareStructs(unchangedResp, body); err != nil {
		t.Fatal(err)
	}

	// now change the vote
	resp = api.Put("/api/v1/vote", cookie1, map[string]any{
		"positions": map[string]string{
			"president": nominationId0,
			// allow putting no vote
		},
	})
	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}
	resp = api.Get("/api/v1/vote", cookie1)
	if resp.Code != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
	body = models.Vote{}
	err = json.Unmarshal(resp.Body.Bytes(), &body)
	if err != nil {
		t.Fatalf("expected body, got error %v", err)
	}
	changedResp := models.Vote{
		Positions: map[string]string{
			"president": nominationId0,
		},
		CreatedAt: now0,
		UpdatedAt: now1,
	}
	if err := compareStructs(changedResp, body); err != nil {
		t.Fatal(err)
	}

	// delete vote
	resp = api.Delete("/api/v1/vote", cookie1)
	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}
	resp = api.Get("/api/v1/vote", cookie1)
	if resp.Code != 404 {
		t.Fatalf("expected 404 Not Found, got %d", resp.Code)
	}
}
