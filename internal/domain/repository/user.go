package repository

import (
	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	authorizationSql = `SELECT 
							email,
							role,
							created_at
						FROM public.users
						WHERE email = $1 AND password = $2;`
)

type UserRepository interface {
	GetDatabase() postgres.Database
	Authorize(email, passwordHash string) (model.User, error)
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

func (u *userRepository) Authorize(email string, passwordHash string) (model.User, error) {
	var result model.User

	err := u.db.QueryRow(authorizationSql, email, passwordHash).Scan(&result.Email, &result.Role, &result.CreatedAt)

	return result, err
}
