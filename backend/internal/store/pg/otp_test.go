package pg

import (
	"testing"
	"time"

	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/store"

	"github.com/pashagolub/pgxmock/v4"
)

func shimNow(st *PgOTPStore, testNow time.Time) {
	nowProvider := func() time.Time {
		return testNow
	}
	st.NowProvider = nowProvider
}

func TestCreateOrReplaceEmptyRatelimit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	testNow := time.Now()
	zid := "z0000000"
	code := "123123"

	mock.ExpectBegin()
	mock.ExpectQuery(`select.*from otp_ratelimit`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRatelimitRows))

	mock.ExpectExec(`insert into otp_ratelimit`).
		WithArgs(zid, testNow).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectExec(`insert into otp.*on conflict \(zid\) do update set`).
		WithArgs(zid, st.hashCode(code), testNow).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	shimNow(st, testNow)
	if err := st.CreateOrReplace(ctx, zid, code); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateOrReplaceSuccessRatelimit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	testNow := time.Now()
	zid := "z0000000"
	code := "123123"
	if st.ratelimitCount < 1 {
		t.Fatalf("ratelimit count is less than 1")
	}
	ratelimitCount := st.ratelimitCount - 1

	mock.ExpectBegin()
	mock.ExpectQuery(`select.*from otp_ratelimit`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRatelimitRows).
			AddRow(zid, ratelimitCount, testNow))

	mock.ExpectExec(`update otp_ratelimit set count = \$2`).
		WithArgs(zid, ratelimitCount + 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectExec(`insert into otp.*on conflict \(zid\) do update set`).
		WithArgs(zid, st.hashCode(code), testNow).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	shimNow(st, testNow)
	if err := st.CreateOrReplace(ctx, zid, code); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateOrReplaceSuccessAfterRatelimit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	testNow := time.Now()
	testAfterRatelimit := testNow.Add(st.rateLimitWithin).Add(1 * time.Second)

	zid := "z0000000"
	code := "123123"
	if st.ratelimitCount < 1 {
		t.Fatalf("ratelimit count is less than 1")
	}
	ratelimitCount := st.ratelimitCount // will fail if not after

	mock.ExpectBegin()
	mock.ExpectQuery(`select.*from otp_ratelimit`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRatelimitRows).
			AddRow(zid, ratelimitCount, testNow))

	// set new window_start
	mock.ExpectExec(`update otp_ratelimit set count = 0`).
		WithArgs(zid, testAfterRatelimit).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectExec(`insert into otp.*on conflict \(zid\) do update set`).
		WithArgs(zid, st.hashCode(code), testAfterRatelimit).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	shimNow(st, testAfterRatelimit)
	if err := st.CreateOrReplace(ctx, zid, code); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

var (
	otpRows          = []string{"zid", "code_hash", "retry_amount", "created_at"}
	otpRatelimitRows = []string{"zid", "count", "window_start"}
)

func TestActiveExists(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	testNowBegin := time.Now()
	zid := "z0000000"
	code := "123123"

	mock.ExpectQuery(`select \* from otp.*where zid = \$1`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRows).
			AddRow(zid, st.hashCode(code), 0, testNowBegin))

	otp, err := st.Active(ctx, zid)
	if err != nil {
		t.Fatal(err)
	}
	if otp == nil {
		t.Fatalf("OTP should exist")
	}
	expected := store.OTPEntry{
		Zid:         zid,
		CodeHash:    st.hashCode(code),
		RetryAmount: 0,
		CreatedAt:   testNowBegin,
	}
	if *otp != expected {
		t.Fatalf("otp != expected, got %v", *otp)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestActiveNoExist(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	zid := "z0000000"

	mock.ExpectQuery(`select \* from otp.*where zid = \$1`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRows))

	otp, err := st.Active(ctx, zid)
	if err != nil {
		t.Fatal(err)
	}
	if otp != nil {
		t.Fatalf("OTP shouldn't exist")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestValidateAndConsumeNoEntry(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	zid := "z0000000"

	mock.ExpectBegin()
	mock.ExpectQuery(`select \* from otp.*where zid = \$1.*for update`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRows))
	mock.ExpectRollback()

	valid, reason, err := st.ValidateAndConsume(ctx, zid, "123123")
	if err != nil {
		t.Fatal(err)
	}
	if valid == true {
		t.Fatalf("should be invalid")
	}
	if reason != store.OTPValidateNotFoundOrExpired {
		t.Fatalf("incorrect reason")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestValidateAndConsumeExpired(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	testNowBegin := time.Now()
	testNowAfterExpiry := testNowBegin.Add(st.expiry).Add(1 * time.Second)

	zid := "z0000000"
	code := "123123"

	mock.ExpectBegin()
	mock.ExpectQuery(`select \* from otp.*where zid = \$1.*for update`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRows).
			AddRow(zid, st.hashCode(code), 0, testNowBegin))

	shimNow(st, testNowAfterExpiry)
	valid, reason, err := st.ValidateAndConsume(ctx, zid, code)
	if err != nil {
		t.Fatal(err)
	}
	if valid == true {
		t.Fatalf("should be invalid")
	}
	if reason != store.OTPValidateNotFoundOrExpired {
		t.Fatalf("incorrect reason")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestValidateAndConsumeAttemptsExceeded(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	testNowBegin := time.Now()
	zid := "z0000000"
	code := "123123"

	mock.ExpectBegin()
	mock.ExpectQuery(`select \* from otp.*where zid = \$1.*for update`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRows).
			AddRow(zid, st.hashCode(code), st.maxRetry, testNowBegin))

	mock.ExpectBegin()
	mock.ExpectExec(`delete from otp where zid`).
		WithArgs(zid). // deleted
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(`delete from otp_ratelimit where zid`).
		WithArgs(zid). // deleted
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	shimNow(st, testNowBegin)
	valid, reason, err := st.ValidateAndConsume(ctx, zid, code)
	if err != nil {
		t.Fatal(err)
	}
	if valid == true {
		t.Fatalf("should be invalid")
	}
	if reason != store.OTPValidateAttemptsExceeded {
		t.Fatalf("incorrect reason")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestValidateAndConsumeMismatch(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	testNowBegin := time.Now()
	zid := "z0000000"
	code := "123123"
	wrongCode := "321321"

	mock.ExpectBegin()
	mock.ExpectQuery(`select \* from otp.*where zid = \$1.*for update`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRows).
			AddRow(zid, st.hashCode(code), 0, testNowBegin))

	mock.ExpectExec(`update otp set retry_amount`).
		WithArgs(zid, 1). // incremented
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()

	shimNow(st, testNowBegin)
	valid, reason, err := st.ValidateAndConsume(ctx, zid, wrongCode)
	if err != nil {
		t.Fatal(err)
	}
	if valid == true {
		t.Fatalf("should be invalid")
	}
	if reason != store.OTPValidateMismatch {
		t.Fatalf("incorrect reason")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestValidateAndConsumeSuccess(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	testNowBegin := time.Now()
	zid := "z0000000"
	code := "123123"

	mock.ExpectBegin()
	mock.ExpectQuery(`select \* from otp.*where zid = \$1.*for update`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRows).
			AddRow(zid, st.hashCode(code), 0, testNowBegin))

	mock.ExpectBegin()
	mock.ExpectExec(`delete from otp where zid`).
		WithArgs(zid). // deleted
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(`delete from otp_ratelimit where zid`).
		WithArgs(zid). // deleted
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	mock.ExpectCommit()

	shimNow(st, testNowBegin)
	valid, reason, err := st.ValidateAndConsume(ctx, zid, code)
	if err != nil {
		t.Fatal(err)
	}
	if valid == false {
		t.Fatalf("should be valid")
	}
	if reason != store.OTPValidateSuccess {
		t.Fatalf("incorrect reason")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestConsumeIfExists(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*PgOTPStore)

	zid := "z0000000"

	mock.ExpectBegin()
	mock.ExpectExec(`delete from otp where zid`).
		WithArgs(zid). // deleted
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(`delete from otp_ratelimit where zid`).
		WithArgs(zid). // deleted
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	if err := st.ConsumeIfExists(ctx, zid); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
