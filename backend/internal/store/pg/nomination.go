package pg

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/linuxunsw/vote/backend/internal/store"
)

type PgNominationStore struct {
	// *pgx.Pool
	pool PgxPoolIface

	NowProvider func() time.Time
	v7Provider  func() (uuid.UUID, error)
}

func NewPgNominationStore(pool PgxPoolIface) store.NominationStore {
	return &PgNominationStore{
		pool: pool,

		NowProvider: time.Now,
		v7Provider:  uuid.NewV7,
	}
}

func (st *PgNominationStore) SubmitOrReplaceNomination(ctx context.Context, electionID string, candidateZid string, submission store.SubmitNomination) (string, error) {
	now := st.NowProvider()
	nominationId, err := st.v7Provider()
	if err != nil {
		return "", err
	}

	_, err = st.pool.Exec(ctx, `
		insert into nominations (
			nomination_id,
		
			candidate_zid,
			election_id,

			candidate_name,
			contact_email,
			discord_username,

			executive_roles,
			candidate_statement,
			url,

			created_at,
			updated_at
		) values ($10, $1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
		on conflict (candidate_zid, election_id) do update set
			candidate_name = EXCLUDED.candidate_name,
			contact_email = EXCLUDED.contact_email,
			discord_username = EXCLUDED.discord_username,
			executive_roles = EXCLUDED.executive_roles,
			candidate_statement = EXCLUDED.candidate_statement,
			url = EXCLUDED.url,
			updated_at = EXCLUDED.updated_at;
	`, candidateZid, electionID,
		submission.CandidateName, submission.ContactEmail, submission.DiscordUsername,
		submission.ExecutiveRoles, submission.CandidateStatement, submission.URL,
		now, nominationId,
	)

	if err != nil {
		return "", err
	}

	return nominationId.String(), nil
}

func (st *PgNominationStore) GetNomination(ctx context.Context, electionID string, candidateZid string) (*store.Nomination, error) {
	rows, err := st.pool.Query(ctx, `
		select * from nominations
		where election_id = $1 and candidate_zid = $2
	`, electionID, candidateZid)

	if err != nil {
		return nil, err
	}

	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[store.Nomination])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (st *PgNominationStore) GetNominationByPublicId(ctx context.Context, nominationId string) (*store.Nomination, error) {
	rows, err := st.pool.Query(ctx, `
		select * from nominations
		where nomination_id = $1
	`, nominationId)

	if err != nil {
		return nil, err
	}

	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[store.Nomination])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (st *PgNominationStore) GetElectionNominations(ctx context.Context, electionId string) ([]store.Nomination, error) {
	rows, err := st.pool.Query(ctx, `
		select * from nominations
		where election_id = $1
		order by created_at asc
	`, electionId)

	if err != nil {
		return nil, err
	}

	entries, err := pgx.CollectRows(rows, pgx.RowToStructByName[store.Nomination])
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func (st *PgNominationStore) TryDeleteNomination(ctx context.Context, electionID string, candidateZid string) error {
	_, err := st.pool.Exec(ctx, `
		delete from nominations
		where election_id = $1 and candidate_zid = $2
	`, electionID, candidateZid)
	return err
}
