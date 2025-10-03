package handlers_test

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/go-chi/httplog/v3"
	v1 "github.com/linuxunsw/vote/backend/internal/api/v1"
	"github.com/linuxunsw/vote/backend/internal/api/v1/models"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/logger"
	"github.com/linuxunsw/vote/backend/internal/mailer/mock_mailer"
	"github.com/linuxunsw/vote/backend/internal/store/pg"
	"github.com/linuxunsw/vote/backend/internal/store/pg/harness"
)

const (
	TestingDummyJWTAdmin = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ6MTIzNDU2NyIsImlzcyI6InZvdGUtYXBpIiwiaXNBZG1pbiI6dHJ1ZX0.pzeRspChlGtMEc3XuVLjDpFriJzdehXqny7L85VxSW0"
)

func TestMain(m *testing.M) {
	harness.HarnessMain(m)
}

// See backend/cmd/api/main.go
func NewAPIWithNowProvider(t *testing.T, nowProvider func() time.Time) (humatest.TestAPI, *mock_mailer.MockMailer) {
	cfg := config.Load()

	// intial cfg for both logger and httplog middleware
	logFormat := httplog.SchemaOTEL.Concise(cfg.Logger.Concise)
	loggerOpts := &slog.HandlerOptions{
		ReplaceAttr: logFormat.ReplaceAttr,
	}

	// init mailer and db
	mailer := mock_mailer.NewMockMailer()
	pool := harness.EphemeralPool(t)

	logger, err := logger.New(cfg.Logger, loggerOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		t.FailNow()
	}

	// init api
	_, api := humatest.New(t)

	// setup stores
	otpStore := pg.NewPgOTPStore(pool, cfg.OTP)
	electionStore := pg.NewPgElectionStore(pool)
	nominationStore := pg.NewPgNominationStore(pool)
	ballotStore := pg.NewPgBallotStore(pool)

	otpStore.(*pg.PgOTPStore).NowProvider = nowProvider
	electionStore.(*pg.PgElectionStore).NowProvider = nowProvider
	nominationStore.(*pg.PgNominationStore).NowProvider = nowProvider
	ballotStore.(*pg.PgBallotStore).NowProvider = nowProvider

	stores := v1.HandlerDependencies{
		Logger:          logger,
		Cfg:             cfg,
		Mailer:          mailer,
		Checker:         nil,
		OtpStore:        otpStore,
		ElectionStore:   electionStore,
		NominationStore: nominationStore,
		BallotStore:     ballotStore,
	}

	v1.Register(api, stores)

	return api, mailer.(*mock_mailer.MockMailer)
}

func NewAPI(t *testing.T) (humatest.TestAPI, *mock_mailer.MockMailer) {
	return NewAPIWithNowProvider(t, time.Now)
}

func createElection(t *testing.T, api humatest.TestAPI, cfg config.JWTConfig, jwt string, memberList []string) string {
	adminCookie := fmt.Sprintf("Cookie: %s=%s", cfg.CookieName, jwt)

	// create an election
	resp := api.Post("/api/v1/elections", adminCookie, map[string]any{
		"name": "Test Election",
	})
	if resp.Code != 200 {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}

	electionResp := models.CreateElectionResponseBody{}
	err := json.Unmarshal(resp.Body.Bytes(), &electionResp)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	electionId := electionResp.ElectionId

	// put member list
	resp = api.Put("/api/v1/elections/"+electionId+"/members", adminCookie, map[string]any{
		"zids": memberList,
	})

	// "No Content" with PUT
	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}

	return electionId
}

// Assumes there will be one cookie sent back, the JWT cookie.
func extractCookieHeader(headers http.Header) string {
	return "Cookie: " + "SESSION=" + extractJWT(headers)
}

func extractJWT(headers http.Header) string {
	sc := headers.Get("Set-Cookie")
	parts := strings.Split(sc, ";")
	kv := strings.TrimSpace(parts[0]) // "SESSION=abc123"

	// remove = from the left

	kvParts := strings.SplitN(kv, "=", 2)
	if len(kvParts) != 2 {
		log.Fatalf("invalid cookie format: %s", kv)
	}
	return kvParts[1]
}

func generateOTPSubmit(t *testing.T, api humatest.TestAPI, mailer *mock_mailer.MockMailer, zid string) *httptest.ResponseRecorder {
	resp := api.Post("/api/v1/otp/generate", map[string]any{
		"zid": zid,
	})
	// generate returns nothing, it goes to the mailer
	if resp.Code != 204 {
		t.Fatalf("expected 204 OK, got %d", resp.Code)
	}

	code := mailer.MockRetrieveOTP(zid + "@unsw.edu.au")

	return api.Post("/api/v1/otp/submit", map[string]any{
		"zid": zid,
		"otp": code,
	})
}

func compareStructs[T any](expected, actual T) error {
	expectedValue := reflect.ValueOf(expected)
	actualValue := reflect.ValueOf(actual)

	// Ensure both inputs are structs. The generic type T ensures they are of the same underlying type.
	if expectedValue.Kind() != reflect.Struct || actualValue.Kind() != reflect.Struct {
		return fmt.Errorf("CompareStructs: both inputs must be structs, got expected %s and actual %s",
			expectedValue.Kind(), actualValue.Kind())
	}

	structType := expectedValue.Type()
	diffs := []string{}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		expectedField := expectedValue.Field(i)
		actualField := actualValue.Field(i)

		// handle time.Time
		if expectedField.Type() == reflect.TypeOf(time.Time{}) && actualField.Type() == reflect.TypeOf(time.Time{}) {
			et := expectedField.Interface().(time.Time)
			at := actualField.Interface().(time.Time)
			// Treat zero times as equal only if both are zero; otherwise use Equal
			if (!et.IsZero() || !at.IsZero()) && !et.Equal(at) {
				diffs = append(diffs, fmt.Sprintf("  field '%s': expected %v (type %s), got %v (type %s)",
					field.Name, et, expectedField.Type(), at, actualField.Type()))
			}
			continue
		}

		// reflect.DeepEqual handles most built-in types, slices, maps, and
		// nested structs recursively.
		if !reflect.DeepEqual(expectedField.Interface(), actualField.Interface()) {
			diffs = append(diffs, fmt.Sprintf("  field '%s': expected %v (type %s), got %v (type %s)",
				field.Name,
				expectedField.Interface(), expectedField.Type(),
				actualField.Interface(), actualField.Type()))
		}
	}

	if len(diffs) > 0 {
		return fmt.Errorf("CompareStructs: structs differ:\nExpected: %+v\nActual:   %+v\nDifferences:\n%s",
			expected, actual, strings.Join(diffs, "\n"))
	}

	return nil
}
