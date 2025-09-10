package handlers_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/go-chi/httplog/v3"
	v1 "github.com/linuxunsw/vote/backend/internal/api/v1"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/logger"
	"github.com/linuxunsw/vote/backend/internal/mailer/mock_mailer"
	"github.com/linuxunsw/vote/backend/internal/store/pg"
	"github.com/linuxunsw/vote/backend/internal/store/pg/harness"
)

func TestMain(m *testing.M) {
	harness.HarnessMain(m)
}

// See backend/cmd/api/main.go
func NewAPI(t *testing.T) (humatest.TestAPI, *mock_mailer.MockMailer) {
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

	stores := v1.HandlerDependencies{
		Logger:   logger,
		Cfg:      cfg,
		Mailer:   mailer,
		Checker:  nil,
		OtpStore: otpStore,
		ElectionStore: electionStore,
	}

	v1.Register(api, stores)

	return api, mailer.(*mock_mailer.MockMailer)
}
