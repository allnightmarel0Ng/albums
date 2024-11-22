package repository

import (
	"errors"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type AuthorizationRepository interface {
	Authorize(email, passwordHash string) (model.User, error)
}

type authorizationRepository struct {
	users repository.UserRepository
}

func NewAuthorizationRepository(users repository.UserRepository) AuthorizationRepository {
	return &authorizationRepository{
		users: users,
	}
}

func (a *authorizationRepository) Authorize(email, passwordHash string) (model.User, error) {
	result, err := a.users.Authorize(email, passwordHash)
	if err != nil {
		err = errors.New("unable to find such user in database")
	}

	return result, err
}
