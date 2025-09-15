package store

import (
	"errors"
	"testing"
)

func TestElectionState_TryTransition(t *testing.T) {
	testCases := []struct {
		name          string
		currentState  string
		nextState     string
		expectedError error
	}{
		// --- Normal, valid transitions from the lowest state ---
		{"From Closed to Nominations Open", "CLOSED", "NOMINATIONS_OPEN", nil},
		{"From Nominations Open to Nominations Closed", "NOMINATIONS_OPEN", "NOMINATIONS_CLOSED", nil},
		{"From Nominations Closed to Voting Open", "NOMINATIONS_CLOSED", "VOTING_OPEN", nil},
		{"From Voting Open to Voting Closed", "VOTING_OPEN", "VOTING_CLOSED", nil},
		{"From Voting Closed to Results", "VOTING_CLOSED", "RESULTS", nil},
		{"From Results to End", "RESULTS", "END", nil},

		// --- Test transition to self ---
		{"Transition to Self: Closed", "CLOSED", "CLOSED", nil},
		{"Transition to Self: Nominations Open", "NOMINATIONS_OPEN", "NOMINATIONS_OPEN", nil},
		{"Transition to Self: End", "END", "END", nil},

		// --- Test backtracking within a group (should be allowed) ---
		{"Backtrack in Group: Nominations Closed to Open", "NOMINATIONS_CLOSED", "NOMINATIONS_OPEN", nil},
		{"Backtrack in Group: Voting Closed to Open", "VOTING_CLOSED", "VOTING_OPEN", nil},

		// --- Error States: Invalid States ---
		{"Invalid State: From a known state to an unknown state", "NOMINATIONS_OPEN", "UNKNOWN_STATE", ErrElectionTransitionInvalidState},
		{"Invalid State: From an unknown state to a known state", "NON_EXISTENT_STATE", "NOMINATIONS_OPEN", ErrElectionTransitionInvalidState},
		{"Invalid State: From a known state to an empty string", "VOTING_CLOSED", "", ErrElectionTransitionInvalidState},
		{"Invalid State: From an empty string to a known state", "", "NOMINATIONS_OPEN", ErrElectionTransitionInvalidState},

		// --- Error States: Cannot Regress ---
		{"Regression: From Nominations Open to Closed", "NOMINATIONS_OPEN", "CLOSED", ErrElectionTransitionCannotRegress},
		{"Regression: From Voting Open to Nominations Closed", "VOTING_OPEN", "NOMINATIONS_CLOSED", ErrElectionTransitionCannotRegress},
		{"Regression: From Results to Voting Closed", "RESULTS", "VOTING_CLOSED", ErrElectionTransitionCannotRegress},
		{"Regression: From End to Results", "END", "RESULTS", ErrElectionTransitionCannotRegress},

		// --- Error States: Cannot Jump ---
		{"Jump: From Closed to Voting Open", "CLOSED", "VOTING_OPEN", ErrElectionTransitionCannotJump},
		{"Jump: From Nominations Open to Voting Open", "NOMINATIONS_OPEN", "VOTING_OPEN", ErrElectionTransitionCannotJump},
		{"Jump: From Nominations Closed to Results", "NOMINATIONS_CLOSED", "RESULTS", ErrElectionTransitionCannotJump},
		{"Jump: From Voting Open to Results", "VOTING_OPEN", "RESULTS", ErrElectionTransitionCannotJump},
		{"Jump: From Closed to End", "CLOSED", "END", ErrElectionTransitionCannotJump},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Cast the current state string to ElectionState before calling the method
			_, err := ElectionState(tc.currentState).TryTransition(tc.nextState)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("For transition from %s to %s, expected error: '%v', but got: '%v'", tc.currentState, tc.nextState, tc.expectedError, err)
			}
		})
	}
}
