package postgres_test

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/denidarta/dauth/internal/core"
	"github.com/denidarta/dauth/internal/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func migrationsPath() string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "migrations")
}

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("dauth_test"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		tcpostgres.WithInitScripts(filepath.Join(migrationsPath(), "000001_init.up.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err, "starting postgres container")
	t.Cleanup(func() { pgContainer.Terminate(context.Background()) })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "getting connection string")

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err, "creating connection pool")
	t.Cleanup(func() { pool.Close() })

	return pool
}

func createTestUser(t *testing.T, repo *postgres.UserRepo, email string) *core.User {
	t.Helper()
	user, err := repo.Create(context.Background(), &core.User{
		Email:        email,
		PasswordHash: "$2a$12$fakehashfortest",
	})
	require.NoError(t, err)
	return user
}
