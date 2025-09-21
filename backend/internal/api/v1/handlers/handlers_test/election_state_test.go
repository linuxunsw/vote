package handlers_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/config"
)

type electionStateTestCase struct {
	states        []string
	expect400Last bool
	expect422Last bool
}

// The new set of test cases mirroring the previous examples and logic.
var testCases = []electionStateTestCase{
	// --- Valid Lifecycle Sequences ---
	// A full, valid sequence from the beginning to the end.
	{states: []string{"CLOSED", "NOMINATIONS_OPEN", "NOMINATIONS_CLOSED", "VOTING_OPEN", "VOTING_CLOSED", "RESULTS", "END"}},
	// A shorter, valid sequence.
	{states: []string{"CLOSED", "NOMINATIONS_OPEN", "NOMINATIONS_CLOSED", "VOTING_OPEN"}},

	// --- Valid Backtracking within a Group ---
	// After opening nominations, close them, then reopen them. The final step is valid.
	{states: []string{"CLOSED", "NOMINATIONS_OPEN", "NOMINATIONS_CLOSED", "NOMINATIONS_OPEN"}},
	// After opening voting, close it, then reopen it. The final step is valid.
	{states: []string{"CLOSED", "NOMINATIONS_OPEN", "NOMINATIONS_CLOSED", "VOTING_OPEN", "VOTING_CLOSED", "VOTING_OPEN"}},

	// --- Valid Transitions to Self ---
	// Transition from a state to itself.
	{states: []string{"NOMINATIONS_OPEN", "NOMINATIONS_OPEN"}},
	// Transition to self in the middle of a sequence.
	{states: []string{"CLOSED", "NOMINATIONS_OPEN", "NOMINATIONS_OPEN"}},

	// --- Error Case: Invalid Jumps ---
	// Jump from CLOSED directly to VOTING_OPEN.
	{states: []string{"CLOSED", "VOTING_OPEN"}, expect400Last: true},
	// Jump from NOMINATIONS_OPEN (an entry state) to the next group's entry state.
	{states: []string{"CLOSED", "NOMINATIONS_OPEN", "VOTING_OPEN"}, expect400Last: true},
	// Jump over the VOTING group entirely.
	{states: []string{"CLOSED", "NOMINATIONS_OPEN", "NOMINATIONS_CLOSED", "RESULTS"}, expect400Last: true},

	// --- Error Case: Invalid Regressions ---
	// Regress from VOTING_OPEN back to NOMINATIONS_CLOSED.
	{states: []string{"CLOSED", "NOMINATIONS_OPEN", "NOMINATIONS_CLOSED", "VOTING_OPEN", "NOMINATIONS_CLOSED"}, expect400Last: true},
	// Regress from RESULTS back to VOTING_CLOSED.
	{states: []string{"CLOSED", "NOMINATIONS_OPEN", "NOMINATIONS_CLOSED", "VOTING_OPEN", "VOTING_CLOSED", "RESULTS", "VOTING_CLOSED"}, expect400Last: true},

	// --- Error Case: Invalid States ---
	// Transition from a valid state to a non-existent state.
	{states: []string{"CLOSED", "NOMINATIONS_OPEN", "SOME_INVALID_STATE"}, expect422Last: true},
}

func TestElectionState_MegaFuzz(t *testing.T) {
	// 1. create election
	// 2. transition up to a certain state, expect no errors, then transition to the next state
	//    - create elections on each iteration, check 400 error code for rejection
	//    - verify valid transitions, code 204, on each and last iteration
	// 3. if the current state is invalid, expect 400 error code
	// 4. if the current state if END, create election again, expect 200 OK

	for _, tc := range testCases {
		t.Run(strings.Join(tc.states, "_"), func(t *testing.T) {
			megafuzzMain(t, tc)
		})
	}
}

// Returns (electionID, code)
func createElectionResp(t *testing.T, api humatest.TestAPI, cfg config.JWTConfig, jwt string) (string, int) {
	adminCookie := fmt.Sprintf("Cookie: %s=%s", cfg.CookieName, jwt)

	// create an election
	resp := api.Post("/api/v1/elections", adminCookie, map[string]any{
		"name": "Test Election",
	})
	if resp.Code != 200 {
		return "", resp.Code
	}

	electionResp := models.CreateElectionResponseBody{}
	err := json.Unmarshal(resp.Body.Bytes(), &electionResp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	return electionResp.ElectionId, resp.Code
}

func getCurrentElectionState(t *testing.T, api humatest.TestAPI, cfg config.JWTConfig, jwt string) models.GetElectionStateResponseBody {
	adminCookie := fmt.Sprintf("Cookie: %s=%s", cfg.CookieName, jwt)

	resp := api.Get("/api/v1/state", adminCookie)
	if resp.Code != 200 {
		t.Fatalf("expected 200 OK on get election, got %d", resp.Code)
	}

	electionStateResp := models.GetElectionStateResponseBody{}
	err := json.Unmarshal(resp.Body.Bytes(), &electionStateResp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	return electionStateResp
}

func transitionElectionState(t *testing.T, api humatest.TestAPI, cfg config.JWTConfig, jwt string, newState string) int {
	adminCookie := fmt.Sprintf("Cookie: %s=%s", cfg.CookieName, jwt)

	resp := api.Put("/api/v1/state", adminCookie, map[string]any{
		"state": newState,
	})

	return resp.Code
}

//var stateToField = map[string]pgx.Identifier{}

func normaliseTime(t time.Time) time.Time {
	return t.Truncate(time.Second)
}

func megafuzzMain(t *testing.T, tc electionStateTestCase) {
	nowProvider := new(func() time.Time)

	cfg := config.Load()
	api, _ := NewAPIWithNowProvider(t, func() time.Time {
		return (*nowProvider)()
	})

	stateMap := map[string]time.Time{}

	// TODO, we do not check `created_at` times
	*nowProvider = func() time.Time { return time.Now() }
	_, code := createElectionResp(t, api, cfg.JWT, TestingDummyJWTAdmin)
	if code != 200 {
		t.Fatalf("expected 200 OK on initial election creation, got %d", code)
	}

	// skip first
	for i, state := range tc.states {
		if i == 0 {
			continue // skip first
		}
		last := i == len(tc.states)-1

		// (leadup) check backtracking doesn't actually mess up the previous timestamps
		now := time.Now()
		*nowProvider = func() time.Time { return now }
		if _, ok := stateMap[state]; !ok {
			stateMap[state] = normaliseTime(now)
		}

		// (transition) transition to the next state
		code = transitionElectionState(t, api, cfg.JWT, TestingDummyJWTAdmin, state)
		if last && (tc.expect400Last || tc.expect422Last) {
			if tc.expect400Last && code != 400 {
				t.Fatalf("expected 400 Bad Request on last transition to invalid state %s, got %d", state, code)
			} else if tc.expect422Last && code != 422 {
				t.Fatalf("expected 422 Unprocessable Entity on last transition to invalid state %s, got %d", state, code)
			}
			return // end test case, failed successfully
		}
		if code != 204 {
			t.Fatalf("expected 204 No Content on transition to state %s, got %d", state, code)
		}

		// (check) verify the current state
		electionStateResp := getCurrentElectionState(t, api, cfg.JWT, TestingDummyJWTAdmin)

		if state == "END" {
			if electionStateResp.State != "NO_ELECTION" {
				t.Fatalf("expected state NO_ELECTION after END, got %s", electionStateResp.State)
			}
		} else {
			if electionStateResp.State != state {
				t.Fatalf("expected state %s, got %s", state, electionStateResp.State)
			}

			// (check) verify the time, against backtracks
			timeReturned, err := time.Parse(time.RFC3339Nano, electionStateResp.StateCreatedAt)
			if err != nil {
				t.Fatalf("failed to parse returned time: %v", err)
			}
			timeReturned = normaliseTime(timeReturned)
			if !timeReturned.Equal(stateMap[state]) {
				// backtracking shouldn't change the timestamp
				t.Fatalf("expected state created at %s, got %s", stateMap[state], timeReturned)
			}
		}

		// (check) create an election on each iteration
		_, code := createElectionResp(t, api, cfg.JWT, TestingDummyJWTAdmin)

		if state == "END" {
			if code != 200 {
				t.Fatalf("expected 200 OK on election creation after END, got %d", code)
			}
		} else if code == 200 {
			t.Fatalf("expected non-200 code on election creation when state is not CLOSED or END, got %d", code)
		}
	}

	// (check) verify last
	electionStateResp := getCurrentElectionState(t, api, cfg.JWT, TestingDummyJWTAdmin)
	lastState := tc.states[len(tc.states)-1]

	if lastState == "END" {
		// should be closed, we made a new one
		if electionStateResp.State != "CLOSED" {
			t.Fatalf("expected state CLOSED after END, got %s", electionStateResp.State)
		}
	}
}
