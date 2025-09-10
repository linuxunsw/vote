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

type Election struct {
	ElectionID string    `db:"election_id"`
	Name       string    `db:"name"`
	State      string    `db:"state"`
	CreatedAt  time.Time `db:"created_at"`

	NominationsOpenAt  *time.Time `db:"nominations_open_at"`
	NominationsCloseAt *time.Time `db:"nominations_close_at"`
	VotingOpenAt       *time.Time `db:"voting_open_at"`
	VotingCloseAt      *time.Time `db:"voting_close_at"`
	ResultsPublishedAt *time.Time `db:"results_published_at"`
}

var ErrElectionSetMembersFailedValidation = errors.New("zID failed to validate")
var ErrElectionCreateAlreadyRunning = errors.New("election is already running")

type ElectionStore interface {
	// Sets members for the current election. Returns error ErrElectionSetMembersFailedValidation
	// and aborts if any zid fails to validate.
	SetMembers(ctx context.Context, electionId string, entries []ElectionMemberEntry) error

	// Get data for member by zid in this election. Returns nil with no error if not found.
	GetMember(ctx context.Context, electionId string, zid string) (*ElectionMemberEntry, error)

	// Create a new election and return its ID. Returns error ErrElectionCreateAlreadyRunning if
	// there is already an election in progress (not in RESULTS state).
	CreateElection(ctx context.Context, name string) (string, error)

	// Get the current election, or nil if none exists.
	CurrentElection(ctx context.Context) (*Election, error)
}
