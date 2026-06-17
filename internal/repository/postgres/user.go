package postgres

import (
	"context"
	"fmt"

	"github.com/denidarta/dauth/internal/core"
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

func (r *UserRepo) Create(ctx context.Context, user *core.User) (*core.User, error) {
	query := `
	INSERT INTO user (email, password_hash)
	VALUES ($1, $2)
	RETURNING id, email, password_hash, created_at, updated_at`
	row := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash)
	result := &core.User{}
	err := row.Scan(
		&result.ID,
		&result.Email,
		&result.PasswordHash,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	return result, nil
}
