package news_test

import (
	"context"
	"fmt"
	"github.com/boinkkitty/newsapi/internal/news"
	"github.com/boinkkitty/newsapi/internal/postgres"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	pgtc "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/uptrace/bun"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var db *bun.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	pdb, cf, err := createTestDB(ctx)
	if err != nil {
		panic(err)
	}

	db = pdb
	code := m.Run()

	if err := cf(ctx); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestStore_Create(t *testing.T) {
	news.NewStore(db)
}

func createTestContainer(ctx context.Context) (ctr *pgtc.PostgresContainer, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return ctr, fmt.Errorf("working directory: %w", err)
	}

	sqlScripts := filepath.Join(wd, "testdata", "sql", "store.sql")

	ctr, err = pgtc.Run(
		ctx,
		"postgres:16-alpine",
		pgtc.WithInitScripts(sqlScripts),
		pgtc.WithDatabase("postgres"),
		pgtc.WithUsername("postgres"),
		pgtc.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second),
		),
	)

	if err != nil {
		return ctr, fmt.Errorf("run container: %w", err)
	}
	return ctr, nil
}

type DBCleanUpFunc func(ctx context.Context) error

func createTestDB(ctx context.Context) (*bun.DB, DBCleanUpFunc, error) {
	ctr, err := createTestContainer(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("create test container: %w", err)
	}

	p, err := ctr.MappedPort(ctx, nat.Port("5432/tcp"))
	if err != nil {
		return nil, nil, fmt.Errorf("mapped port: %w", err)
	}

	db, err := postgres.NewDB(&postgres.Config{
		Host:     "localhost",
		Debug:    true,
		DBName:   "postgres",
		User:     "postgres",
		Password: "postgres",
		Port:     p.Port(),
		SSLMode:  "disable",
	})
	if err != nil {
		return nil, nil, fmt.Errorf("new db: %w", err)
	}

	cf := func(ctx context.Context) error {
		if err := db.Close(); err != nil {
			return fmt.Errorf("close db: %w", err)
		}
		if err := ctr.Terminate(ctx); err != nil {
			return fmt.Errorf("container terminate: %w", err)
		}
		return nil
	}

	return db, cf, nil
}
