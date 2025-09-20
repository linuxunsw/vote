package pg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/linuxunsw/vote/backend/internal/store"
)

type PgBallotStore struct {
	// *pgx.Pool
	pool PgxPoolIface

	NowProvider func() time.Time
}

func NewPgBallotStore(pool PgxPoolIface) store.BallotStore {
	return &PgBallotStore{
		pool: pool,

		NowProvider: time.Now,
	}
}

func (b *PgBallotStore) SubmitOrReplaceBallot(ctx context.Context, electionID string, zid string, submission store.SubmitBallot) error {
	submissionObject, err := json.Marshal(submission.Positions)
	if err != nil {
		return err
	}
	
	// Upsert the ballot
	_, err = b.pool.Exec(ctx, `
		insert into ballots (election_id, zid, positions, created_at, updated_at)
		values ($1, $2, $3, $4, $4)
		on conflict (election_id, zid) do update
			set positions = $3, updated_at = $4
	`, electionID, zid, string(submissionObject), b.NowProvider())
	if err != nil {
		return err
	}

	return nil
}

func (b *PgBallotStore) GetBallot(ctx context.Context, electionID string, zid string) (*store.Ballot, error) {
	rows, err := b.pool.Query(ctx, `
		select positions, created_at, updated_at
		from ballots
		where election_id = $1 and zid = $2
	`, electionID, zid)
	if err != nil {
		return nil, err
	}

	ballot, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[store.Ballot])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ballot, nil
}

func (b *PgBallotStore) TryDeleteBallot(ctx context.Context, electionID string, zid string) error {
	_, err := b.pool.Exec(ctx, `
		delete from ballots
		where election_id = $1 and zid = $2
	`, electionID, zid)
	if err != nil {
		return err
	}
	return nil
}
