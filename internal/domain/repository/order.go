package repository

import (
	"context"
	"strings"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	selectAddAlbumProcedureSQL =
	/* sql */ `CALL add_album_to_user_order($1, $2);`
	selectDeleteAlbumProcedureSQL =
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
				WHERE u.id = $1
				`
	isPaidFilterSQl = " AND o.is_paid = FALSE;"
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
	_, err := o.db.Query(ctx, selectAddAlbumProcedureSQL, userID, albumID)
	return err
}

func (o *orderRepository) DeleteAlbumFromUserOrder(ctx context.Context, userID, albumID int) error {
	_, err := o.db.Query(ctx, selectDeleteAlbumProcedureSQL, userID, albumID)
	return err
}

func (o *orderRepository) GetUserOrders(ctx context.Context, userID int, unpaidOnly bool) ([]model.Order, error) {
	var sb strings.Builder
	_, err := sb.WriteString(selectUserOrdersSQL)
	if err != nil {
		return nil, err
	}

	switch unpaidOnly {
	case true:
		_, err := sb.WriteString(isPaidFilterSQl)
		if err != nil {
			return nil, err
		}
	case false:
		_, err := sb.WriteRune(';')
		if err != nil {
			return nil, err
		}
	}

	rows, err := o.db.Query(ctx, sb.String(), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[int]*model.Order)

	for rows.Next() {
		var (
			order model.Order
			album model.Album
		)

		err := rows.Scan(&order.ID, &order.Orderer.ID, &order.Orderer.Email, &order.Orderer.IsAdmin, &order.Orderer.Nickname, &order.Orderer.Balance, &order.Orderer.ImageURL, &order.Date, &order.TotalPrice, &order.IsPaid, &album.ID, &album.Name, &album.Author.ID, &album.Author.Name, &album.Author.Genre, &album.Author.ImageURL, &album.ImageURL, &album.Price)
		if err != nil {
			return nil, err
		}

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
