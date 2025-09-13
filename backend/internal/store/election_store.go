package store

import (
	"context"
	"errors"
	"time"
)

type ElectionMemberEntry struct {
	ElectionID string `db:"election_id"`
	Zid        string `db:"zid"`
}

type ElectionState string

const (
	StateClosed            ElectionState = "CLOSED"
	StateNominationsOpen   ElectionState = "NOMINATIONS_OPEN"
	StateNominationsClosed ElectionState = "NOMINATIONS_CLOSED"
	StateVotingOpen        ElectionState = "VOTING_OPEN"
	StateVotingClosed      ElectionState = "VOTING_CLOSED"
	StateResults           ElectionState = "RESULTS"
	StateEnd               ElectionState = "END"
)

const (
	entry int = 1
	exit  int = 2
)

type stateMapEntry struct {
	group int
	kind int
}

// statemap assumes there are at most two per group, one entry and one exit
var stateMap = map[ElectionState]stateMapEntry{
	StateClosed:            {group: 0, kind: entry | exit},
	//
	StateNominationsOpen:   {group: 1, kind: entry},
	StateNominationsClosed: {group: 1, kind: exit},
	//
	StateVotingOpen:        {group: 2, kind: entry},
	StateVotingClosed:      {group: 2, kind: exit},
	//
	StateResults:           {group: 3, kind: entry | exit},
	StateEnd:               {group: 4, kind: entry | exit},
}

var (
	ErrElectionTransitionInvalidState = errors.New("invalid election state")
	ErrElectionTransitionCannotJump = errors.New("cannot jump election states")
	ErrElectionTransitionCannotRegress = errors.New("cannot regress election states")
)

func (e ElectionState) TryTransition(toState string) (ElectionState, error) {
	to := ElectionState(toState)
	
	if e == to {
		return e, nil
	}

	fromEntry, ok1 := stateMap[e]
	toEntry, ok2 := stateMap[to]
	if !ok1 || !ok2 {
		return "", ErrElectionTransitionInvalidState
	}

	if toEntry.group < fromEntry.group {
		return "", ErrElectionTransitionCannotRegress
	}

	// two pairs, can go back and forth
	if toEntry.group == fromEntry.group {
		return to, nil
	}

	// jumping to the next group
	// from exit to next entry
	if toEntry.group == fromEntry.group + 1 {
		if fromEntry.kind & exit == 0 {
			return "", ErrElectionTransitionCannotJump
		}
		if toEntry.kind & entry == 0 {
			return "", ErrElectionTransitionCannotJump
		}
		return to, nil
	}

	// jumping forward at least two groups
	return "", ErrElectionTransitionCannotJump
}

type Election struct {
	ElectionID string        `db:"election_id"`
	Name       string        `db:"name"`
	State      ElectionState `db:"state"`
	CreatedAt  time.Time     `db:"created_at"`

	NominationsOpenAt  *time.Time `db:"nominations_open_at"`
	NominationsCloseAt *time.Time `db:"nominations_close_at"`
	VotingOpenAt       *time.Time `db:"voting_open_at"`
	VotingCloseAt      *time.Time `db:"voting_close_at"`
	ResultsPublishedAt *time.Time `db:"results_published_at"`
	EndedAt            *time.Time `db:"ended_at"`
}

func (e Election) StateCreatedAt() time.Time {
	switch e.State {
	case "CLOSED":
		return e.CreatedAt
	case "NOMINATIONS_OPEN":
		return *e.NominationsOpenAt
	case "NOMINATIONS_CLOSED":
		return *e.NominationsCloseAt
	case "VOTING_OPEN":
		return *e.VotingOpenAt
	case "VOTING_CLOSED":
		return *e.VotingCloseAt
	case "RESULTS":
		return *e.ResultsPublishedAt
	case "END":
		return *e.EndedAt
	}
	return e.CreatedAt
}

var ErrElectionSetMembersFailedValidation = errors.New("zID failed to validate")
var ErrElectionCreateAlreadyRunning = errors.New("election is already running")

var ErrElectionNotFound = errors.New("election not found")

type ElectionStore interface {
	// Sets members for the current election. This is a string array of zIDs.
	// Returns error ErrElectionSetMembersFailedValidation and aborts if any zID
	// fails to validate.
	// Returns error ErrElectionNotFound if the election referenced by ID does not exist.
	SetMembers(ctx context.Context, electionId string, entries []string) error

	// Get data for member by zid in this election. Returns nil with no error if not found.
	// Returns error ErrElectionNotFound if the election referenced by ID does not exist.
	GetMember(ctx context.Context, electionId string, zid string) (*ElectionMemberEntry, error)

	// Create a new election and return its ID. Returns error ErrElectionCreateAlreadyRunning if
	// there is already an election in progress (not in RESULTS state).
	CreateElection(ctx context.Context, name string) (string, error)

	// Get the current election, or nil if none exists.
	CurrentElection(ctx context.Context) (*Election, error)

	// a string type for newState has been chosen to better communicate intent

	// Sets the state of the current to newState. Does not assume newState is a valid
	// transition. Performs the validation and state transition inside.
	// Returns error ErrElectionNotFound if the election no election is running.
	//
	// Can return the following errors on invalid transitions:
	//   ErrElectionTransitionInvalidState
	//   ErrElectionTransitionCannotJump
	//   ErrElectionTransitionCannotRegress
	CurrentElectionSetState(ctx context.Context, newState string) error
}
