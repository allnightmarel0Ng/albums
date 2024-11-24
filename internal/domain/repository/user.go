package repository

import (
	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	authorizationSql = /* sql */ `SELECT 
							email,
							role,
							created_at,
							password_hash
						FROM public.users
						WHERE email = $1;`
)

type UserRepository interface {
	GetDatabase() postgres.Database
	Authorize(email string) (model.User, string, error)
}

type userRepository struct {
	db postgres.Database
}

func NewUserRepository(db postgres.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) GetDatabase() postgres.Database {
	return u.db
}

func (u *userRepository) Authorize(email string) (model.User, string, error) {
	var result model.User
	var passwordHash string

	err := u.db.QueryRow(authorizationSql, email).Scan(&result.Email, &result.Role, &result.CreatedAt, &passwordHash)

	return result, passwordHash, err
}
