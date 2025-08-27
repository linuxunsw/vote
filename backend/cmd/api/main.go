package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	v1 "github.com/linuxunsw/vote/backend/internal/api/v1"
	"github.com/linuxunsw/vote/backend/internal/api/v1/handlers"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/logger"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

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

	// logging
	log, err := logger.New(cfg.Logger.Level, cfg.Logger.Encoding)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	slog := log.Sugar()
	defer slog.Sync()

	// db
	// FIXME: implement
	// pool, err := db.Connect(cfg.Database.Address)
	// if err != nil {
	// 	log.Fatal("Unable to connect to database", zap.Error(err))
	// }
	// defer pool.Close()

	// init repository and email client
	// store := TODO:
	// emailClient := TODO:

	// start healthcheck
	health := handlers.NewChecker(slog, nil)
	defer health.Stop()

	// init api
	router := http.NewServeMux()
	api := humago.New(router, huma.DefaultConfig("Vote API", cfg.API.Version))

	v1.Register(api, nil, nil, health)

	// cli & env parsing for high level config and commands
	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", opts.Port),
			Handler: router,
		}

		hooks.OnStart(func() {
			slog.Infow("Starting server", "Port", opts.Port)
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				slog.Fatal("Server error", zap.Error(err))
			}
		})

		hooks.OnStop(func() {
			// Graceful shutdown :)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				slog.Error("Graceful shutdown failed", zap.Error(err))
			}
		})
	})

	// Define cli info
	cmd := cli.Root()
	cmd.Use = "vote-api"
	cmd.Version = cfg.API.Version

	// Add commands to cli
	cmd.AddCommand(createOpenAPICommand(api))
	cmd.AddCommand(createMigrateCommand(slog, cfg))

	// TODO: register more commands
	// i.e running tests(?), (de)registering admins(?)

	// When no commands are passed, this starts the server!
	cli.Run()
}

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

func createMigrateCommand(log *zap.SugaredLogger, cfg config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			// Goose requires a standard *sql.DB object instead of pgx stuff.
			db, err := sql.Open("pgx", cfg.Database.Address)
			if err != nil {
				log.Fatalf("Failed to connect to DB for migration: %v\n", err)
				os.Exit(1)
			}
			defer db.Close()

			if err := goose.SetDialect("postgres"); err != nil {
				log.Errorf("Failed to set Goose dialect: %v\n", err)
				os.Exit(1)
			}

			// NOTE: maybe use embeded migrations? see docs
			// FIXME: use correct directory
			fmt.Println("Running migrations...")
			if err := goose.Up(db, "db/migrations"); err != nil {
				log.Errorf("Migration failed: %v\n", err)
				os.Exit(1)
			}
			log.Info("Migrations applied successfully!")
		},
	}
}
