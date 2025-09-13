package pg

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/linuxunsw/vote/backend/internal/store"
)

type PgElectionStore struct {
	// *pgx.Pool
	pool PgxPoolIface

	NowProvider func() time.Time
	v7Provider  func() (uuid.UUID, error)
}

type query interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func NewPgElectionStore(pool PgxPoolIface) store.ElectionStore {
	return &PgElectionStore{
		pool: pool,

		NowProvider: time.Now,
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
				Zid:        zid,
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

func (st *PgElectionStore) assertElectionExists(ctx context.Context, tx query, electionId string) error {
	rows, err := tx.Query(ctx, `
		select 1 as exists from elections
		where election_id = $1
	`, electionId)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[struct{ Exists int }])
	if errors.Is(err, pgx.ErrNoRows) {
		return store.ErrElectionNotFound
	} else if err != nil {
		return store.ErrElectionNotFound
	}

	return nil
}

func (st *PgElectionStore) SetMembers(ctx context.Context, electionId string, entries []string) error {
	tx, err := st.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// ErrElectionNotFound
	if err := st.assertElectionExists(ctx, tx, electionId); err != nil {
		return err
	}

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

func (st *PgElectionStore) GetMember(ctx context.Context, electionId string, zid string) (*store.ElectionMemberEntry, error) {
	// ErrElectionNotFound
	if err := st.assertElectionExists(ctx, st.pool, electionId); err != nil {
		return nil, err
	}

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

func (st *PgElectionStore) CreateElection(ctx context.Context, name string) (string, error) {
	now := st.NowProvider()

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

func (st *PgElectionStore) currentElection(ctx context.Context, tx query) (*store.Election, error) {
	// get latest election
	rows, err := tx.Query(ctx, `
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

func (st *PgElectionStore) CurrentElection(ctx context.Context) (*store.Election, error) {
	return st.currentElection(ctx, st.pool)
}

/* func (st *pgElectionStore) GetElection(ctx context.Context, electionId string) (*store.Election, error) {
	rows, err := st.pool.Query(ctx, `
		select * from elections
		where election_id = $1
	`, electionId)
	if err != nil {
		return nil, err
	}

	election, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[store.Election])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &election, nil
} */

// timestamp fields will only be set once
var timestampTransitionAway = map[store.ElectionState]pgx.Identifier{
	// created_at is always set on creation
	"CLOSED":             {"nominations_open_at"},
	"NOMINATIONS_OPEN":   {"nominations_close_at"},
	"NOMINATIONS_CLOSED": {"voting_open_at"},
	"VOTING_OPEN":        {"voting_close_at"},
	"VOTING_CLOSED":      {"results_published_at"},
	"RESULTS":            {"ended_at"},
	"END":                nil, // final state, no timestamp
}

func (st *PgElectionStore) CurrentElectionSetState(ctx context.Context, newStateString string) error {
	tx, err := st.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	current, err := st.currentElection(ctx, tx)
	if err != nil {
		return err
	}
	if current == nil {
		return store.ErrElectionNotFound
	}

	currentState := current.State
	newState, err := currentState.TryTransition(newStateString)
	if err != nil {
		return err
	}
	if newState == currentState {
		// no-op
		return nil
	}

	// transition from currentState -> newState
	now := st.NowProvider()
	timestampIdent := timestampTransitionAway[currentState]

	_, err = tx.Exec(ctx, `
		update elections
		set state = $1
		where election_id = $2
	`, newState, current.ElectionID)
	if err != nil {
		return err
	}

	if timestampIdent != nil {
		// pgx.Identifier does perform sanitisation
		sql := fmt.Sprintf(`
			update elections
			set %s = $1
			where election_id = $2 and %s is null
		`, timestampIdent.Sanitize(), timestampIdent.Sanitize())
		_, err = tx.Exec(ctx, sql, now, current.ElectionID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
