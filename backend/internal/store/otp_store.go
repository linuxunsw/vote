package store

import (
	"context"
	"errors"
	"time"
)

type (
	OTPValidate int
)

const (
	OTPValidateSuccess OTPValidate = iota
	OTPValidateInternalError
	OTPValidateNotFoundOrExpired
	OTPValidateAttemptsExceeded
	OTPValidateMismatch
)

type OTPEntry struct {
	Zid         string    `db:"zid"`
	CodeHash    string    `db:"code_hash"`
	RetryAmount int       `db:"retry_amount"`
	CreatedAt   time.Time `db:"created_at"`
}

var ErrOTPRateLimitExceeded = errors.New("rate limit exceeded")

type OTPStore interface {
	// Create or replace OTP entry given zid and code. This will manage ratelimits
	// and return error ErrOTPRateLimitExceeded if the ratelimit is exceeded.
	CreateOrReplace(ctx context.Context, zid string, code string) (err error)

	// Gets the active OTP entry for a given zid. Returns valid OTPEntry if exists,
	// returns nil OTPEntry if nonexistent.
	Active(ctx context.Context, zid string) (*OTPEntry, error)

	// Validates a plaintext code with the entry in the database.
	ValidateAndConsume(ctx context.Context, zid string, code string) (valid bool, reason OTPValidate, err error)

	// Consumes an OTP code owned by zid unconditionally. Clears ratelimits.
	ConsumeIfExists(ctx context.Context, zid string) error
}

func (v OTPValidate) ToString() string {
	switch v {
	case OTPValidateSuccess:
		return "success"
	case OTPValidateInternalError:
		return "internal_error"
	case OTPValidateNotFoundOrExpired:
		return "not found or expired"
	case OTPValidateAttemptsExceeded:
		return "attempts exceeded"
	case OTPValidateMismatch:
		return "mismatch"
	default:
		return "unknown"
	}
}
