package postgres_test

import (
	"context"
	"testing"

	"github.com/denidarta/dauth/internal/core"
	"github.com/denidarta/dauth/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepo_Create(t *testing.T) {
	repo := postgres.NewUserRepo(setupTestDB(t))

	user, err := repo.Create(context.Background(), &core.User{
		Email:        "alice@example.com",
		PasswordHash: "$2a$12$fakehashfortest",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, user.ID, "should generate a user ID")
	assert.Equal(t, "alice@example.com", user.Email)
	assert.Equal(t, "$2a$12$fakehashfortest", user.PasswordHash)
	assert.False(t, user.CreatedAt.IsZero(), "should set created_at")
	assert.False(t, user.UpdatedAt.IsZero(), "should set updated_at")
}

func TestUserRepo_Create_DuplicateEmail(t *testing.T) {
	repo := postgres.NewUserRepo(setupTestDB(t))

	createTestUser(t, repo, "bob@example.com")

	_, err := repo.Create(context.Background(), &core.User{
		Email:        "bob@example.com",
		PasswordHash: "$2a$12$fakehash2",
	})

	assert.Error(t, err, "duplicate email should be rejected")
}

func TestUserRepo_FindByEmail(t *testing.T) {
	repo := postgres.NewUserRepo(setupTestDB(t))
	created := createTestUser(t, repo, "carol@example.com")

	found, err := repo.FindByEmail(context.Background(), "carol@example.com")

	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, "carol@example.com", found.Email)
	assert.Equal(t, created.PasswordHash, found.PasswordHash)
}

func TestUserRepo_FindByEmail_NotFound(t *testing.T) {
	repo := postgres.NewUserRepo(setupTestDB(t))

	_, err := repo.FindByEmail(context.Background(), "nobody@example.com")

	assert.ErrorIs(t, err, core.ErrUserNotFound)
}

func TestUserRepo_FindByID(t *testing.T) {
	repo := postgres.NewUserRepo(setupTestDB(t))
	created := createTestUser(t, repo, "dave@example.com")

	found, err := repo.FindByID(context.Background(), created.ID)

	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, "dave@example.com", found.Email)
	assert.Equal(t, created.PasswordHash, found.PasswordHash)
}

func TestUserRepo_FindByID_NotFound(t *testing.T) {
	repo := postgres.NewUserRepo(setupTestDB(t))

	_, err := repo.FindByID(context.Background(), "00000000-0000-0000-0000-000000000000")

	assert.ErrorIs(t, err, core.ErrUserNotFound)
}
