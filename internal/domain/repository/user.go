package repository

import (
	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	selectIDPasswordHashByEmailSQL =
	/* sql */ `SELECT
					id,
					password_hash,
					is_admin
				FROM public.users
				WHERE email = $1;`
	selectUserByEmailSQL =
	/* sql */ `SELECT
					id,
					email,
					is_admin,
					nickname,
					balance,
					image_url
				FROM public.users
				WHERE id = $1;`
)

type UserRepository interface {
	GetIDPasswordHash(email string) (int, string, bool, error)
	GetUser(id int) (model.User, error)
}

type userRepository struct {
	db postgres.Database
}

func NewUserRepository(db postgres.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) GetIDPasswordHash(email string) (int, string, bool, error) {
	var (
		id           int
		passwordHash string
		isAdmin      bool
	)

	err := u.db.QueryRow(selectIDPasswordHashByEmailSQL, email).Scan(&id, &passwordHash, &isAdmin)
	return id, passwordHash, isAdmin, err
}

func (u *userRepository) GetUser(id int) (model.User, error) {
	var result model.User

	err := u.db.QueryRow(selectUserByEmailSQL, id).Scan(&result.ID, &result.Email, &result.IsAdmin, &result.Nickname, &result.Balance, &result.ImageURL)
	if err != nil {
		return model.User{}, err
	}

	return result, err
}
