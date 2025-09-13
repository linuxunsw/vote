package harness

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/linuxunsw/vote/backend/internal/config"
	"github.com/linuxunsw/vote/backend/internal/store/migrations"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// these are globals set up by the testing harness. these stay throughout over all tests

var controlDBConn *pgx.Conn
var connectionUrlConfig *pgx.ConnConfig

func configToURL(conf *pgx.ConnConfig) string {
	// build a URL of the form:
	// postgresql://user:password@host:port/database?params
	u := &url.URL{
		Scheme: "postgresql",
		User:   url.UserPassword(conf.User, conf.Password),
		Host:   conf.Host,
		Path:   "/" + conf.Database,
	}

	// if a port was specified in conf.Port (which is an int), append it to the Host.
	// if not, the Host string might already include the port
	if conf.Port != 0 {
		u.Host = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	}

	q := url.Values{}
	for k, v := range conf.RuntimeParams {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func HarnessMain(m *testing.M) {
	var err error
	
	cfg := config.Load()

	controlDBConn, err = pgx.Connect(context.Background(), cfg.Database.Address)
	if err != nil {
		log.Fatalf("failed to connect to control database: %v", err)
	}

	connectionUrlConfig, err = pgx.ParseConfig(cfg.Database.Address)
	if err != nil {
		log.Fatalf("failed to parse database config: %v", err)
	}

	os.Exit(m.Run())
}

func EphemeralPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := t.Context()

	safeTestName := strings.ReplaceAll(t.Name(), "/", "_") // Sanitize test name for DB name
	safeTestName = strings.ReplaceAll(safeTestName, "\\", "_")
	testDBName := fmt.Sprintf("test-%s-%d", strings.ToLower(safeTestName), time.Now().UnixNano())

	if len(testDBName) > 31 {
		testDBName = testDBName[:31]
	}
	testDBName += fmt.Sprintf("_%d", os.Getpid())

	createDBSQL := fmt.Sprintf("CREATE DATABASE %s;", pgx.Identifier{testDBName}.Sanitize())
	_, err := controlDBConn.Exec(ctx, createDBSQL)

	if err != nil {
		t.Fatalf("failed to create ephemeral database %s: %v", testDBName, err)
	}

	connectionConfig := connectionUrlConfig.Copy()
	connectionConfig.Database = testDBName

	if err != nil {
		t.Fatalf("failed to connect to ephemeral database %s: %v", testDBName, err)
	}

	sqlDB, err := sql.Open("pgx", configToURL(connectionConfig))
	if err != nil {
		t.Fatalf("failed to open sql.DB for ephemeral database %s: %v", testDBName, err)
	}

	// migrations
	migrationsProvider, err := goose.NewProvider(goose.DialectPostgres, sqlDB, migrations.Migrations)
	if err != nil {
		t.Fatalf("failed to create migrations provider for ephemeral database %s: %v", testDBName, err)
	}
	_, err = migrationsProvider.Up(ctx) // apply all migrations
	if err != nil {
		t.Fatalf("failed to run migrations on ephemeral database %s: %v", testDBName, err)
	}
	_ = sqlDB.Close()

	pool, err := pgxpool.New(ctx, configToURL(connectionConfig))

	if os.Getenv("NO_CLEANUP") == "" {
		t.Cleanup(func () {
			if pool != nil {
				pool.Close()
			}
			dropDBSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE);", pgx.Identifier{testDBName}.Sanitize())
			_, dropErr := controlDBConn.Exec(context.Background(), dropDBSQL)
			if dropErr != nil {
				t.Errorf("failed to drop ephemeral database %s: %v", testDBName, dropErr)
			}
		})
	} else {
		log.Printf("EphemeralPool: Test '%v' Database '%v'", t.Name(), testDBName)
		t.Cleanup(func () {
			if pool != nil {
				pool.Close()
			}
		})
	}
	if err != nil {
		t.Fatalf("failed to create connection pool to ephemeral database %s: %v", testDBName, err)
	}

	return pool
}
