package store

import (
	"context"
	"time"
)

type SubmitBallot struct {
	Positions map[string]string `db:"positions"`
}

type Ballot struct {
	Positions map[string]string `db:"positions"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type BallotStore interface {
	// Submit or replace an existing ballot for the given election and voter zid.
	// Does not perform validation of the submission.
	SubmitOrReplaceBallot(ctx context.Context, electionID string, zid string, submission SubmitBallot) error

	// Get a ballot by election ID and voter zID. Returns nil if not found.
	GetBallot(ctx context.Context, electionID string, zid string) (*Ballot, error)

	// Delete a ballot by election ID and voter zID. Does nothing if the ballot doesn't exist.
	TryDeleteBallot(ctx context.Context, electionID string, zid string) error
}
