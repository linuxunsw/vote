package pg

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/linuxunsw/vote/backend/internal/store"
)

type pgElectionStore struct {
	// *pgx.Pool
	pool     PgxPoolIface

	nowProvider func() time.Time
	v7Provider  func() (uuid.UUID, error)
}

func NewPgElectionStore(pool PgxPoolIface) store.ElectionStore {
	return &pgElectionStore{
		pool: pool,

		nowProvider: time.Now,
		v7Provider:  uuid.NewV7,
	}
}

var (
	zIDRegex = regexp.MustCompile(`^z[0-9]{7}$`)
)

func deduplicateAndVerifyEntries(entries []string, electionId string) (deduplicated []store.ElectionMemberEntry, valid bool) {
	dedupEntries := make(map[string]store.ElectionMemberEntry)
	for _, zid := range entries {
		if !zIDRegex.MatchString(zid) {
			return nil, false
		}
		
		if _, exists := dedupEntries[zid]; !exists {
			dedupEntries[zid] = store.ElectionMemberEntry{
				Zid: zid,
				ElectionID: electionId,
			}
		}
	}

	dedupEntriesList := make([]store.ElectionMemberEntry, 0, len(dedupEntries))
	for _, entry := range dedupEntries {
		dedupEntriesList = append(dedupEntriesList, entry)
	}

	return dedupEntriesList, true
}

func (st *pgElectionStore) SetMembers(ctx context.Context, electionId string, entries []string) error {
	tx, err := st.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	deduplicatedEntries, valid := deduplicateAndVerifyEntries(entries, electionId)
	if !valid {
		return store.ErrElectionSetMembersFailedValidation
	}

	// we create a primary key here purely for the index, we already deduplicate
	_, err = tx.Exec(ctx, `
		create temp table temp_election_member_list (
			zid text,
			election_id uuid,

			primary key ("election_id", "zid")
		) on commit drop;
	`)
	if err != nil {
		return err
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"temp_election_member_list"},
		[]string{"election_id", "zid"},
		pgx.CopyFromSlice(len(deduplicatedEntries), func(i int) ([]interface{}, error) {
			entry := deduplicatedEntries[i]
			return []interface{}{
				electionId,
				entry.Zid,
			}, nil
		}),
	)
	if err != nil {
		return err
	}

	// compute add and remove sets

	// the alternative to this would be to copy all the data back to the client and perform
	// the set difference there. that would be more expensive. it would be much better to perform
	// nlog(n) set "difference" in the database as the data needs to go there anyway

	_, err = tx.Exec(ctx, `
		delete from election_member_list
		where election_id = $1
		and zid not in (select zid from temp_election_member_list where election_id = $1);
	`, electionId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
		insert into election_member_list (zid, election_id)
		select zid, election_id from temp_election_member_list
		on conflict do nothing;
	`)
	if err != nil {
		return err
	}

	// commit transaction, the temporary table will be dropped
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (st *pgElectionStore) GetMember(ctx context.Context, electionId string, zid string) (*store.ElectionMemberEntry, error) {
	rows, err := st.pool.Query(ctx, `
		select * from election_member_list
		where election_id = $1 and zid = $2
	`, electionId, zid)
	if err != nil {
		return nil, err
	}

	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[store.ElectionMemberEntry])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (st *pgElectionStore) CreateElection(ctx context.Context, name string) (string, error) {
	now := st.nowProvider()

	current, err := st.CurrentElection(ctx)
	if err != nil {
		return "", err
	}
	if current != nil {
		return "", store.ErrElectionCreateAlreadyRunning
	}

	electionId, err := st.v7Provider()
	if err != nil {
		return "", err
	}

	_, err = st.pool.Exec(ctx, `
		insert into elections (election_id, name, state, created_at)
		values ($1, $2, 'CLOSED', $3)
	`, electionId, name, now)
	if err != nil {
		return "", err
	}
	return electionId.String(), nil
}

func (st *pgElectionStore) CurrentElection(ctx context.Context) (*store.Election, error) {
	// get latest election
	rows, err := st.pool.Query(ctx, `
		select * from elections
		order by created_at desc
		limit 1;
	`)
	if err != nil {
		return nil, err
	}

	election, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[store.Election])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	if election.State == "END" {
		// all elections are finalised, no current election
		return nil, nil
	}

	return &election, nil
}
