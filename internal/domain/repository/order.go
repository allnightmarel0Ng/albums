package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	callAddAlbumProcedureSQL =
	/* sql */ `CALL add_album_to_user_order($1, $2);`
	callDeleteAlbumProcedureSQL =
	/* sql */ `CALL delete_album_from_user_order($1, $2);`
	selectUserOrdersSQL =
	/* sql */ `SELECT
					o.id,
					u.id,
					u.email,
					u.is_admin,
					u.nickname,
					u.balance,
					u.image_url,
					o.date,
					o.total_price,
					o.is_paid,
					a.id,
					a.name,
					ar.id,
					ar.name,
					ar.genre,
					ar.image_url,
					a.image_url,
					a.price
				FROM public.orders AS o
				JOIN public.users AS u ON o.user_id = u.id
				JOIN public.order_items AS oi ON oi.order_id = o.id
				RIGHT JOIN public.albums AS a ON oi.album_id = a.id
				JOIN public.artists AS ar ON a.artist_id = ar.id
				WHERE u.id = $1`
	isPaidFilterSQL = " AND o.is_paid = FALSE"
	orderSQL        = "\tORDER BY o.is_paid, o.id"
)

type OrderRepository interface {
	AddAlbumToUserOrder(ctx context.Context, userID, albumID int) error
	DeleteAlbumFromUserOrder(ctx context.Context, userID, albumID int) error
	GetUserOrders(ctx context.Context, userID int, unpaidOnly bool) ([]model.Order, error)
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
	return callWillSerialization(o.db, ctx, callAddAlbumProcedureSQL, userID, albumID)
}

func (o *orderRepository) DeleteAlbumFromUserOrder(ctx context.Context, userID, albumID int) error {
	return callWillSerialization(o.db, ctx, callDeleteAlbumProcedureSQL, userID, albumID)
}

func (o *orderRepository) GetUserOrders(ctx context.Context, userID int, unpaidOnly bool) ([]model.Order, error) {
	var sb strings.Builder
	_, err := sb.WriteString(selectUserOrdersSQL)
	if err != nil {
		return nil, err
	}

	if unpaidOnly {
		_, err := sb.WriteString(isPaidFilterSQL)
		if err != nil {
			return nil, err
		}
	}

	_, err = sb.WriteString(orderSQL)
	if err != nil {
		return nil, err
	}

	_, err = sb.WriteRune(';')
	if err != nil {
		return nil, err
	}

	rows, err := o.db.Query(ctx, sb.String(), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[int]*model.Order)

	for rows.Next() {
		var (
			order  model.Order
			album  model.Album
			author model.Artist
		)

		err := rows.Scan(&order.ID, &order.Orderer.ID, &order.Orderer.Email, &order.Orderer.IsAdmin, &order.Orderer.Nickname, &order.Orderer.Balance, &order.Orderer.ImageURL, &order.Date, &order.TotalPrice, &order.IsPaid, &album.ID, &album.Name, &author.ID, &author.Name, &author.Genre, &author.ImageURL, &album.ImageURL, &album.Price)
		if err != nil {
			return nil, err
		}

		album.Author = &author

		_, ok := ordersMap[order.ID]
		if !ok {
			order.Albums = []model.Album{album}
			ordersMap[order.ID] = &order
		} else {
			ordersMap[order.ID].Albums = append(ordersMap[order.ID].Albums, album)
		}
	}

	result := make([]model.Order, len(ordersMap))
	index := 0
	for _, v := range ordersMap {
		result[index] = *v
		index++
	}

	return result, nil
}

func callWillSerialization(db postgres.Database, ctx context.Context, sql string, params ...interface{}) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic! %v", r)
		}

		if err == nil {
			err = tx.Commit(ctx)
		} else {
			err = tx.Rollback(ctx)
		}
	}()

	err = tx.Exec(ctx, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;")
	err = tx.Exec(ctx, sql, params...)
	err = tx.Commit(ctx)
	return err
}
