package repository

import (
	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type AuthorizationRepository interface {
	Authorize(email string) (model.User, string, error)
}

type authorizationRepository struct {
	users repository.UserRepository
}

func NewAuthorizationRepository(users repository.UserRepository) AuthorizationRepository {
	return &authorizationRepository{
		users: users,
	}
}

func (a *authorizationRepository) Authorize(email string) (model.User, string, error) {
	result, hash, err := a.users.Authorize(email)
	if err != nil {
		return model.User{}, "", err
	}

	return result, hash, err
}
