package pg

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/linuxunsw/vote/backend/internal/store"
)

type pgNominationStore struct {
	// *pgx.Pool
	pool PgxPoolIface

	nowProvider func() time.Time
}

func NewPgNominationStore(pool PgxPoolIface) store.NominationStore {
	return &pgNominationStore{
		pool: pool,

		nowProvider: time.Now,
	}
}

func (st *pgNominationStore) SubmitOrReplaceNomination(ctx context.Context, electionID string, candidateZid string, submission store.SubmitNomination) error {
	now := st.nowProvider()

	_, err := st.pool.Exec(ctx, `
		insert into nominations (
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
		) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)
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
		now,
	)

	return err
}

func (st *pgNominationStore) GetNomination(ctx context.Context, electionID string, candidateZid string) (*store.Nomination, error) {
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
