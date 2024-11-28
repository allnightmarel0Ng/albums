package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type AuthorizationRepository interface {
	GetIDPasswordHash(ctx context.Context, email string) (int, string, bool, error)
}

type authorizationRepository struct {
	users repository.UserRepository
}

func NewAuthorizationRepository(users repository.UserRepository) AuthorizationRepository {
	return &authorizationRepository{
		users: users,
	}
}

func (a *authorizationRepository) GetIDPasswordHash(ctx context.Context, email string) (int, string, bool, error) {
	return a.users.GetIDPasswordHash(ctx, email)
}
