package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	selectAddAlbumProcedure =
	/* sql */ `SELECT add_album_to_user_order($1, $2);`
	selectDeleteAlbumProcedure =
	/* sql */ `SELECT delete_album_from_user_order($1, $2);`
)

type OrderRepository interface {
	AddAlbumToUserOrder(ctx context.Context, userID, albumID int) error
	DeleteAlbumFromUserOrder(ctx context.Context, userID, albumID int) error
}

type orderRepository struct {
	db postgres.Database
}

func NewOrderRepository(db postgres.Database) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (o *orderRepository) AddAlbumToUserOrder(ctx context.Context, userID, albumID int) error {
	_, err := o.db.Query(ctx, selectAddAlbumProcedure, userID, albumID)
	return err
}

func (o *orderRepository) DeleteAlbumFromUserOrder(ctx context.Context, userID, albumID int) error {
	_, err := o.db.Query(ctx, selectDeleteAlbumProcedure, userID, albumID)
	return err
}
