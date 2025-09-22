package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/httplog/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	v1 "github.com/linuxunsw/vote/backend/internal/api/v1"
	"github.com/linuxunsw/vote/backend/internal/api/v1/handlers"
	"github.com/linuxunsw/vote/backend/internal/api/v1/middleware"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/logger"
	"github.com/linuxunsw/vote/backend/internal/mailer"
	"github.com/linuxunsw/vote/backend/internal/store/migrations"
	"github.com/linuxunsw/vote/backend/internal/store/pg"
	"github.com/pressly/goose/v3"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/spf13/cobra"
)

// For usage, see https://huma.rocks/features/cli/#passing-options
type Options struct {
	Debug bool   `doc:"Enable debug logging"`
	Host  string `doc:"Hostname to listen on."`
	Port  int    `doc:"Port to listen on." short:"p" default:"8888"`
}

func main() {
	cfg := config.Load()

	// intial cfg for both logger and httplog middleware
	logFormat := httplog.SchemaOTEL.Concise(cfg.Logger.Concise)
	loggerOpts := &slog.HandlerOptions{
		ReplaceAttr: logFormat.ReplaceAttr,
	}

	// init logger
	logger, err := logger.New(cfg.Logger, loggerOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}

	// set useful log attrs in prod
	if !cfg.Logger.Concise {
		logger = logger.With(
			slog.String("app", "vote-api"),
			slog.String("version", cfg.API.Version),
		)
	}

	// initialise database
	pool, err := pgxpool.New(context.Background(), cfg.Database.Address)
	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}
	defer pool.Close()

	// start healthcheck
	health := handlers.NewChecker(logger, nil)
	defer health.Stop()

	// init api
	router := http.NewServeMux()
	humaCfg := huma.DefaultConfig("Vote API", cfg.API.Version)
	humaCfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"cookieAuth": {
			Type: "apiKey",
			In:   "cookie",
			Name: cfg.JWT.CookieName,
		},
	}
	api := humago.New(router, humaCfg)

	// init crossorigin
	crossOrigin := http.NewCrossOriginProtection()
	if cfg.Server.AllowedOrigins != nil {
		for _, origin := range cfg.Server.AllowedOrigins {
			if err := crossOrigin.AddTrustedOrigin(origin); err != nil {
				logger.Error("Invalid allowed origin", "origin", origin, "error", err)
				os.Exit(1)
			}
		}
	}

	// global middleware
	opts := middleware.GlobalMiddlewareOptions{
		Logger:          logger,
		LogFormat:       logFormat,
		LoggerCfg:       cfg.Logger,
		CrossOrigin:     crossOrigin,
		RateLimitCfg:    cfg.Server.RateLimit,
		RealIPAllowlist: cfg.Server.RealIPAllowlist,
	}
	err = middleware.AddGlobalMiddleware(api, opts)
	if err != nil {
		logger.Error("Unable to add global middleware", "error", err)
		os.Exit(1)
	}

	// init mailer
	var mail mailer.Mailer
	if cfg.Mailer.ConsoleMailer {
		mail = mailer.NewConsoleMailer(logger)
	} else {
		mail = mailer.NewResendMailer(cfg)
	}

	// setup stores
	otpStore := pg.NewPgOTPStore(pool, cfg.OTP)
	electionStore := pg.NewPgElectionStore(pool)
	nominationStore := pg.NewPgNominationStore(pool)
	ballotStore := pg.NewPgBallotStore(pool)

	deps := v1.HandlerDependencies{
		Logger:          logger,
		Cfg:             cfg,
		Mailer:          mail,
		Checker:         health,
		OtpStore:        otpStore,
		ElectionStore:   electionStore,
		NominationStore: nominationStore,
		BallotStore:     ballotStore,
	}
	v1.Register(api, deps)

	// cli & env parsing for high level config and commands
	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", opts.Port),
			Handler: router,
		}

		hooks.OnStart(func() {
			logger.Info("Starting server", "Port", opts.Port)
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				logger.Error("Server error", "error", err)
				os.Exit(1)
			}
		})

		hooks.OnStop(func() {
			// Graceful shutdown :)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				logger.Error("Graceful shutdown failed", "error", err)
			}
		})
	})

	// Define cli info
	cmd := cli.Root()
	cmd.Use = "vote-api"
	cmd.Version = cfg.API.Version

	// Add commands to cli
	cmd.AddCommand(createOpenAPICommand(api))
	cmd.AddCommand(createMigrateCommand(logger, cfg))

	// TODO: register more commands
	// i.e running tests(?), (de)registering admins(?)

	// When no commands are passed, this starts the server!
	cli.Run()
}

// FIXME: the healthcheck outputs when running this command
func createOpenAPICommand(api huma.API) *cobra.Command {
	return &cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := api.OpenAPI().YAML()
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		},
	}
}

func createMigrateCommand(log *slog.Logger, cfg config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			// Goose requires a standard *sql.DB object instead of pgx stuff.
			db, err := sql.Open("pgx", cfg.Database.Address)
			if err != nil {
				log.Error("Failed to connect to DB for migration", "error", err)
				os.Exit(1)
			}
			defer func() {
				_ = db.Close()
			}()

			migrationsProvider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.Migrations)
			if err != nil {
				log.Error("Failed to create migrations provider", "error", err)
				os.Exit(1)
			}

			fmt.Println("Running migrations...")
			migrations, err := migrationsProvider.Up(context.Background())
			if err != nil {
				log.Error("Failed to run migrations", "error", err)
				os.Exit(1)
			}
			log.Info("migrations applied successfully", "migrations", strconv.Itoa(len(migrations)))
		},
	}
}
