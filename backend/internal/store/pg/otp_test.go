package pg

import (
	"testing"
	"time"

	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/store"

	"github.com/pashagolub/pgxmock/v4"
)

func shimNow(st *pgOTPStore, testNow time.Time) {
	nowProvider := func () time.Time {
		return testNow
	}
	st.nowProvider = nowProvider
}

func TestCreateOrReplace(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*pgOTPStore)

	testNow := time.Now()
	zid := "z0000000"
	code := "123123"

	mock.ExpectExec(`insert into otp.*on conflict \(zid\) do update set`).
		WithArgs(zid, st.hashCode(code), testNow).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	shimNow(st, testNow)
	if err := st.CreateOrReplace(ctx, zid, code); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

var (
	otpRows = []string{"zid", "code_hash", "retry_amount", "created_at"}
)

func TestValidateAndConsumeNoEntry(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := t.Context()
	st := NewPgOTPStore(mock, config.Load().OTP).(*pgOTPStore)

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
	st := NewPgOTPStore(mock, config.Load().OTP).(*pgOTPStore)

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
	st := NewPgOTPStore(mock, config.Load().OTP).(*pgOTPStore)

	testNowBegin := time.Now()
	zid := "z0000000"
	code := "123123"

	mock.ExpectBegin()
	mock.ExpectQuery(`select \* from otp.*where zid = \$1.*for update`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRows).
			AddRow(zid, st.hashCode(code), st.maxRetry, testNowBegin))

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
	st := NewPgOTPStore(mock, config.Load().OTP).(*pgOTPStore)

	testNowBegin := time.Now()
	zid := "z0000000"
	code := "123123"
	wrongCode := "321321"

	mock.ExpectBegin()
	mock.ExpectQuery(`select \* from otp.*where zid = \$1.*for update`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRows).
			AddRow(zid, st.hashCode(code), 0, testNowBegin))

	mock.ExpectExec(`update otp set attempts`).
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
	st := NewPgOTPStore(mock, config.Load().OTP).(*pgOTPStore)

	testNowBegin := time.Now()
	zid := "z0000000"
	code := "123123"

	mock.ExpectBegin()
	mock.ExpectQuery(`select \* from otp.*where zid = \$1.*for update`).
		WithArgs(zid).
		WillReturnRows(pgxmock.NewRows(otpRows).
			AddRow(zid, st.hashCode(code), 0, testNowBegin))

	mock.ExpectExec(`delete from otp where zid`).
		WithArgs(zid). // deleted
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
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
	st := NewPgOTPStore(mock, config.Load().OTP).(*pgOTPStore)

	zid := "z0000000"

	mock.ExpectExec(`delete from otp where zid`).
		WithArgs(zid). // deleted
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	if err := st.ConsumeIfExists(ctx, zid); err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
