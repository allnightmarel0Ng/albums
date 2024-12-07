package repository

import (
	"context"
	"fmt"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	selectIDPasswordHashByEmailSQL =
	/* sql */ `SELECT
					u.id,
					c.password_hash,
					u.is_admin
				FROM public.users AS u
				JOIN public.credentials AS c ON c.user_id = u.id
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
	/* sql */ `UPDATE public.users
				SET balance = balance + ($1)
				WHERE id = $2;`

	callPayForOrderSQL =
	/* sql */ `CALL pay_for_order($1, $2);`

	insertNewUserSQL =
	/* sql */ `INSERT INTO public.users (email, is_admin, nickname, image_url)
				VALUES ($1, $2, $3, $4)
				RETURNING id;`

	insertNewCredentialSQL =
	/* sql */ `INSERT INTO public.credentials (user_id, password_hash)
				VALUES ($1, $2);`

	findEmailSQL =
	/* sql */ `SELECT 
					CASE 
						WHEN EXISTS (SELECT 1 FROM public.users WHERE email = $1) 
						THEN TRUE 
						ELSE FALSE 
					END AS email_exists;`

	selectAlbumOwnersIdsSQL =
	/* sql */ `SELECT user_id
				FROM public.purchased_albums
				WHERE album_id = $1;`
)

type UserRepository interface {
	GetIDPasswordHash(ctx context.Context, email string) (int, string, bool, error)
	GetUser(ctx context.Context, id int) (model.User, error)
	ChangeBalance(ctx context.Context, id int, diff uint) error
	PayForOrder(ctx context.Context, userID int, orderID int) error
	AddNewUser(ctx context.Context, email, password_hash string, isAdmin bool, nickname, imageURL string) error
	FindUserByEmail(ctx context.Context, email string) (bool, error)
	GetAlbumOwnersIds(ctx context.Context, albumID int) ([]int, error)
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

func (u *userRepository) ChangeBalance(ctx context.Context, id int, diff uint) error {
	return u.db.Exec(ctx, updateBalanceSQL, diff, id)
}

func (u *userRepository) PayForOrder(ctx context.Context, userID int, orderID int) error {
	return callWillSerialization(u.db, ctx, callPayForOrderSQL, userID, orderID)
}

func (u *userRepository) AddNewUser(ctx context.Context, email, password_hash string, isAdmin bool, nickname, imageURL string) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v", r)
		}

		if err != nil {
			err = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	var id int
	err = tx.QueryRow(ctx, insertNewUserSQL, email, isAdmin, nickname, imageURL).Scan(&id)
	if err != nil {
		return err
	}

	return tx.Exec(ctx, insertNewCredentialSQL, id, password_hash)
}

func (u *userRepository) FindUserByEmail(ctx context.Context, email string) (bool, error) {
	var result bool
	err := u.db.QueryRow(ctx, findEmailSQL, email).Scan(&result)
	return result, err
}

func (u *userRepository) GetAlbumOwnersIds(ctx context.Context, albumID int) ([]int, error) {
	rows, err := u.db.Query(ctx, selectAlbumOwnersIdsSQL, albumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []int

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		result = append(result, id)
	}
	return result, nil
}
