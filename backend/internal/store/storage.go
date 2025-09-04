package store

import (
	"context"
	"time"
)

type (
	OTPValidate int
	Store       any
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

type OTPStore interface {
	// Create or replace OTP entry given zid and code. Be sure to check if the user
	// has an active OTP code to avoid them repeatedly generating codes.
	CreateOrReplace(ctx context.Context, zid string, code string) (err error)

	// Gets the active OTP entry for a given zid. Returns valid OTPEntry if exists,
	// returns nil OTPEntry if nonexistent.
	Active(ctx context.Context, zid string) (*OTPEntry, error)

	// Validates a plaintext code with the entry in the database.
	ValidateAndConsume(ctx context.Context, zid string, code string) (valid bool, reason OTPValidate, err error)

	// Consumes an OTP code owned by zid unconditionally.
	ConsumeIfExists(ctx context.Context, zid string) error
}
