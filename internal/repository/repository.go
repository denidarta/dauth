package repository

import (
	"context"

	"github.com/denidarta/dauth/internal/core"
)

type UserRepository interface {
	Create(ctx context.Context, user *core.User) (*core.User, error)
	FindByEmail(ctx context.Context, email string) (*core.User, error)
	FindByID(ctx context.Context, id string) (*core.User, error)
}
type ProductRepository interface {
	FindByID(ctx context.Context, id string) (*core.Product, error)
}

type UserProductRepository interface {
	AddMember(ctx context.Context, up *core.UserProduct) error
	IsMember(ctx context.Context, userID, productID string) (bool, error)
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *core.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*core.RefreshToken, error)
	Delete(ctx context.Context, id string) error
}

type InvitationRepository interface {
	Create(ctx context.Context, inv *core.Invitation) error
	FindByHash(ctx context.Context, hash string) (*core.Invitation, error)
	MarkAccepted(ctx context.Context, id string) error
}
