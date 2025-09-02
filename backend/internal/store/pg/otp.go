package pg

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/store"

	"github.com/jackc/pgx/v5"
)

type pgOTPStore struct {
	// *pgx.Pool
	pool     PgxPoolIface
	secret   string
	maxRetry int
	expiry   time.Duration

	ratelimitCount  int
	rateLimitWithin time.Duration

	// cannot create more than `rateLimitCount` OTPs within `rateLimitWithin` duration

	nowProvider func() time.Time
}

func NewPgOTPStore(pool PgxPoolIface, cfg config.OTPConfig) store.OTPStore {
	return &pgOTPStore{
		pool:     pool,
		secret:   cfg.Secret,
		maxRetry: cfg.MaxRetry,
		expiry:   cfg.Duration,

		ratelimitCount:  cfg.RatelimitCount,
		rateLimitWithin: cfg.RatelimitWithin,

		nowProvider: time.Now,
	}
}

func (st *pgOTPStore) hashCode(code string) string {
	mac := hmac.New(sha256.New, []byte(st.secret))
	mac.Write([]byte(code))
	return hex.EncodeToString(mac.Sum(nil))
}

func (st *pgOTPStore) hashCompare(code string, expectedCodeHash string) bool {
	given := st.hashCode(code)
	return hmac.Equal([]byte(expectedCodeHash), []byte(given))
}

type otpRatelimit struct {
	Zid         string    `db:"zid"`
	Count       int       `db:"count"`
	WindowStart time.Time `db:"window_start"`
}

func (st *pgOTPStore) CreateOrReplace(ctx context.Context, zid string, code string) error {
	tx, err := st.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	codeHash := st.hashCode(code)
	now := st.nowProvider()

	ratelimitRows, err := st.pool.Query(ctx, `
		select * from otp_ratelimit where zid = $1
	`, zid)
	if err != nil {
		return err
	}

	ratelimitEntry, err := pgx.CollectOneRow(ratelimitRows, pgx.RowToStructByName[otpRatelimit])
	if errors.Is(err, pgx.ErrNoRows) {
		// no rate limit, reset the window
		_, err = st.pool.Exec(ctx, `
			insert into otp_ratelimit (zid, count, window_start)
			values ($1, 0, $2)
		`, zid, now)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// rate limiting applies
		if ratelimitEntry.WindowStart.Add(st.rateLimitWithin).After(now) {
			if ratelimitEntry.Count >= st.ratelimitCount {
				return store.OTPRateLimitExceeded
			}

			ratelimitEntry.Count++
			_, err := st.pool.Exec(ctx, `
				update otp_ratelimit set count = $2 where zid = $1
			`, zid, ratelimitEntry.Count)
			if err != nil {
				return err
			}
		} else {
			// reset the ratelimit window
			_, err = st.pool.Exec(ctx, `
				update otp_ratelimit
				set count = 0, window_start = $2
				where zid = $1
			`, zid, now)
			if err != nil {
				return err
			}
		}
	}

	_, err = st.pool.Exec(ctx, `
		insert into otp (zid, code_hash, retry_amount, created_at)
		values ($1, $2, 0, $3)
		on conflict (zid) do update set
			code_hash = EXCLUDED.code_hash,
			retry_amount = EXCLUDED.retry_amount,
			created_at = EXCLUDED.created_at;
	`, zid, codeHash, now)
	if err != nil {
		return err
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (st *pgOTPStore) Active(ctx context.Context, zid string) (*store.OTPEntry, error) {
	rows, err := st.pool.Query(ctx, `
		select * from otp
		where zid = $1
	`, zid)

	if err != nil {
		return nil, err
	}

	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[store.OTPEntry])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (st *pgOTPStore) activeForUpdate(ctx context.Context, zid string) (*store.OTPEntry, error) {
	rows, err := st.pool.Query(ctx, `
		select * from otp
		where zid = $1
		for update
	`, zid)

	if err != nil {
		return nil, err
	}

	entry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[store.OTPEntry])
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (st *pgOTPStore) ValidateAndConsume(ctx context.Context, zid string, code string) (valid bool, reason store.OTPValidate, err error) {
	tx, err := st.pool.Begin(ctx)
	if err != nil {
		return false, store.OTPValidateInternalError, err
	}
	defer tx.Rollback(ctx)

	otp, err := st.activeForUpdate(ctx, zid)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, store.OTPValidateNotFoundOrExpired, nil
	} else if err != nil {
		return false, store.OTPValidateInternalError, err
	}

	now := st.nowProvider()
	if otp.CreatedAt.Add(st.expiry).Before(now) {
		return false, store.OTPValidateNotFoundOrExpired, nil
	}

	if otp.RetryAmount >= st.maxRetry {
		return false, store.OTPValidateAttemptsExceeded, nil
	}

	ok := st.hashCompare(code, otp.CodeHash)
	otp.RetryAmount++

	if ok {
		// consume
		err = st.ConsumeIfExists(ctx, zid)
	} else {
		// increment retries
		_, err = tx.Exec(ctx, `
			update otp set attempts = $2
			where zid = $1
		`, zid, otp.RetryAmount)
	}

	if err != nil {
		return false, store.OTPValidateInternalError, err
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		return false, store.OTPValidateInternalError, err
	}
	if ok {
		return true, store.OTPValidateSuccess, nil
	}
	return false, store.OTPValidateMismatch, nil
}

func (st *pgOTPStore) ConsumeIfExists(ctx context.Context, zid string) error {
	tx, err := st.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = st.pool.Exec(ctx, `
		delete from otp where zid = $1;
	`, zid)
	if err != nil {
		return err
	}
	_, err = st.pool.Exec(ctx, `
		delete from otp_ratelimit where zid = $1;
	`, zid)
	if err != nil {
		return err
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
