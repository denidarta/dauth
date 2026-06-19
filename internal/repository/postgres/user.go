package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/denidarta/dauth/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	var repo UserRepo
	repo.db = db
	return &repo
}

func scanUser(row pgx.Row) (*core.User, error) {
	result := &core.User{}
	err := row.Scan(
		&result.ID,
		&result.Email,
		&result.PasswordHash,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UserRepo) Create(ctx context.Context, user *core.User) (*core.User, error) {
	query := `
	INSERT INTO users (email, password_hash)
	VALUES ($1, $2)
	RETURNING id, email, password_hash, created_at, updated_at`
	result, err := scanUser(r.db.QueryRow(ctx, query, user.Email, user.PasswordHash))
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return result, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*core.User, error) {
	query := `
	SELECT id, email, password_hash, created_at, updated_at
	FROM users
	WHERE id = $1`
	result, err := scanUser(r.db.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.ErrUserNotFound
		}
		return nil, fmt.Errorf("finding user by id: %w", err)
	}
	return result, nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*core.User, error) {
	query := `
	SELECT id, email, password_hash, created_at, updated_at
	FROM users
	WHERE email = $1`
	result, err := scanUser(r.db.QueryRow(ctx, query, email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, core.ErrUserNotFound
		}
		return nil, fmt.Errorf("finding user by email: %w", err)
	}
	return result, nil
}
