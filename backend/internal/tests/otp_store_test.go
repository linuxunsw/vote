package tests

import (
	"testing"

	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/store"
	"github.com/linuxunsw/vote/backend/internal/store/pg"
	"github.com/linuxunsw/vote/backend/internal/tests/harness"
)

func TestOTPConsumeOnce(t *testing.T) {
	pool := harness.EphemeralPool(t)
	otpStore := pg.NewPgOTPStore(pool, config.Load().OTP)
	ctx := t.Context()

	zid := "z0000000"
	code := "123123"

	if err := otpStore.CreateOrReplace(ctx, zid, code); err != nil {
		t.Fatalf("CreateOrReplace failed: %v", err)
	}
	entry, err := otpStore.Active(ctx, zid)
	if err != nil {
		t.Fatalf("Active failed: %v", err)
	}
	if entry == nil {
		t.Fatalf("Active returned nil entry")
	}

	ok, reason, err := otpStore.ValidateAndConsume(ctx, zid, code)
	if err != nil {
		t.Fatalf("ValidateAndConsume failed: %v", err)
	}

	if !ok {
		t.Fatalf("ValidateAndConsume returned not ok: %v", reason)
	}
	if reason != store.OTPValidateSuccess {
		t.Fatalf("ValidateAndConsume returned reason: %v", reason)
	}

	entry, err = otpStore.Active(ctx, zid)
	if err != nil {
		t.Fatalf("Active failed: %v", err)
	}
	if entry != nil {
		t.Fatalf("Active returned entry after consumption: %+v", entry)
	}
}

func TestOTPReachRetryLimit(t *testing.T) {
	pool := harness.EphemeralPool(t)
	otpConfig := config.Load().OTP
	otpStore := pg.NewPgOTPStore(pool, otpConfig)
	ctx := t.Context()

	zid := "z0000000"
	code := "123123"
	wrongCode := "321321"

	if err := otpStore.CreateOrReplace(ctx, zid, code); err != nil {
		t.Fatalf("CreateOrReplace failed: %v", err)
	}

	for i := 0; i < otpConfig.MaxRetry + 2; i++ {
		valid, reason, err := otpStore.ValidateAndConsume(ctx, zid, wrongCode)
		if err != nil {
			t.Fatalf("ValidateAndConsume failed at iteration %d: %v", i, err)
		}
		if valid {
			t.Fatalf("ValidateAndConsume returned valid at iteration %d", i)
		}
		if i < otpConfig.MaxRetry && reason != store.OTPValidateMismatch {
			t.Fatalf("ValidateAndConsume reason not OTPValidateMismatch at iteration %d: %v", i, reason)
		} else if i == otpConfig.MaxRetry && reason != store.OTPValidateAttemptsExceeded {
			t.Fatalf("ValidateAndConsume reason not OTPValidateAttemptsExceeded at iteration %d: %v", i, reason)
		} else if i > otpConfig.MaxRetry && reason != store.OTPValidateNotFoundOrExpired {
			t.Fatalf("ValidateAndConsume reason not OTPValidateNotFoundOrExpired at iteration %d: %v", i, reason)
		}
	}
}

func TestOTPReachRatelimit(t *testing.T) {
	pool := harness.EphemeralPool(t)
	otpConfig := config.Load().OTP
	otpStore := pg.NewPgOTPStore(pool, otpConfig)
	ctx := t.Context()

	zid := "z0000000"
	code := "123123"

	for i := 0; i < otpConfig.RatelimitCount + 2; i++ {
		err := otpStore.CreateOrReplace(ctx, zid, code)
		if i <= otpConfig.RatelimitCount && err != nil {
			t.Fatalf("CreateOrReplace failed at iteration %d: %v", i, err)
		} else if i > otpConfig.RatelimitCount && err == nil {
			t.Fatalf("CreateOrReplace succeeded at iteration %d but should have failed due to rate limit", i)
		}
	}

	// should still exist
	if entry, err := otpStore.Active(ctx, zid); err != nil {
		t.Fatalf("Active failed: %v", err)
	} else if entry == nil {
		t.Fatalf("Active returned nil entry")
	}

	// validate now
	valid, reason, err := otpStore.ValidateAndConsume(ctx, zid, code)
	if err != nil {
		t.Fatalf("ValidateAndConsume failed: %v", err)
	}
	if !valid {
		t.Fatalf("ValidateAndConsume returned not valid: %v", reason)
	}
	if reason != store.OTPValidateSuccess {
		t.Fatalf("ValidateAndConsume returned reason: %v", reason)
	}

	// shouldn't exist anymore
	if entry, err := otpStore.Active(ctx, zid); err != nil {
		t.Fatalf("Active failed: %v", err)
	} else if entry != nil {
		t.Fatalf("Active returned entry after consumption: %+v", entry)
	}
}
