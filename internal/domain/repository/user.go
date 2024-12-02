package repository

import (
	"context"

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

	updateBalanceSQL =
	/* sql */ `UPDATE TABLE public.users
				SET balance = balance + ($1)
				WHERE id = $2;`

	callPayForOrderSQL =
	/* sql */ `CALL pay_for_order($1, $2);`
)

type UserRepository interface {
	GetIDPasswordHash(ctx context.Context, email string) (int, string, bool, error)
	GetUser(ctx context.Context, id int) (model.User, error)
	ChangeBalance(ctx context.Context, id, diff int) error
	PayForOrder(ctx context.Context, userID int, orderID int) error
}

type userRepository struct {
	db postgres.Database
}

func NewUserRepository(db postgres.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) GetIDPasswordHash(ctx context.Context, email string) (int, string, bool, error) {
	var (
		id           int
		passwordHash string
		isAdmin      bool
	)

	err := u.db.QueryRow(ctx, selectIDPasswordHashByEmailSQL, email).Scan(&id, &passwordHash, &isAdmin)
	return id, passwordHash, isAdmin, err
}

func (u *userRepository) GetUser(ctx context.Context, id int) (model.User, error) {
	var result model.User

	err := u.db.QueryRow(ctx, selectUserByEmailSQL, id).Scan(&result.ID, &result.Email, &result.IsAdmin, &result.Nickname, &result.Balance, &result.ImageURL)
	if err != nil {
		return model.User{}, err
	}

	return result, err
}

func (u *userRepository) ChangeBalance(ctx context.Context, id, diff int) error {
	return u.db.Exec(ctx, updateBalanceSQL, diff, id)
}

func (u *userRepository) PayForOrder(ctx context.Context, userID int, orderID int) error {
	return u.db.Exec(ctx, callPayForOrderSQL, userID, orderID)
}
