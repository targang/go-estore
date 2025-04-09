package repository

import (
	"context"
	"go_store/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ OrderRepository = (*orderRepositoryImpl)(nil)

type orderRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepositoryImpl{db: db}
}

func (o *orderRepositoryImpl) Create(ctx context.Context, order *model.Order) (string, error) {
	tx, err := o.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const orderInsert = `
INSERT INTO orders (customer_name, customer_email, status)	
VALUES ($1, $2, $3)
RETURNING id
`
	var createdID string

	err = tx.QueryRow(ctx, orderInsert, order.CustomerName, order.CustomerEmail, order.Status).
		Scan(&createdID)
	if err != nil {
		return "", err
	}

	const itemInsert = `
INSERT INTO order_item (order_id, product_id, quantity)
VALUES ($1, $2, $3)
`
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, itemInsert, createdID, item.ProductID, item.Quantity)
		if err != nil {
			return "", err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}
	return createdID, nil
}

func (o *orderRepositoryImpl) GetByID(ctx context.Context, id string) (*model.Order, error) {
	tx, err := o.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const orderQuery = `
SELECT customer_name, customer_email, status, created_at, updated_at
FROM orders 
WHERE id = $1
`
	var order model.Order
	order.ID = id
	err = tx.QueryRow(ctx, orderQuery, id).
		Scan(&order.CustomerName, &order.CustomerEmail, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	const itemsQuery = `
SELECT product_id, quantity
FROM order_item WHERE order_id = $1
`
	rows, err := tx.Query(ctx, itemsQuery, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		if err = rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *orderRepositoryImpl) UpdateStatus(ctx context.Context, order *model.Order) error {
	const query = `
UPDATE orders
SET status = $1 
WHERE id = $2
`
	_, err := o.db.Exec(ctx, query, order.Status, order.ID)

	return err
}

func (o *orderRepositoryImpl) Delete(ctx context.Context, id string) error {
	const query = `
DELETE FROM orders WHERE id = $1
`
	_, err := o.db.Exec(ctx, query, id)
	return err
}

func (o *orderRepositoryImpl) List(ctx context.Context, limit, offset int32) ([]model.Order, error) {
	tx, err := o.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
SELECT id, customer_name, customer_email, status, created_at, updated_at
FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

	rows, err := tx.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	orderMap := make(map[string]*model.Order)

	for rows.Next() {
		var order model.Order
		err = rows.Scan(
			&order.ID,
			&order.CustomerName,
			&order.CustomerEmail,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		order.Items = []model.OrderItem{}
		orders = append(orders, order)
		orderMap[order.ID] = &orders[len(orders)-1]
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(orderMap) > 0 {
		itemQuery := `
SELECT order_id, product_id, quantity
FROM order_item
WHERE order_id = ANY($1)
`
		orderIDs := make([]string, 0, len(orderMap))
		for id := range orderMap {
			orderIDs = append(orderIDs, id)
		}

		itemRows, err := tx.Query(ctx, itemQuery, orderIDs)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item model.OrderItem
			var orderID string
			err = itemRows.Scan(&orderID, &item.ProductID, &item.Quantity)
			if err != nil {
				return nil, err
			}
			if order, ok := orderMap[orderID]; ok {
				order.Items = append(order.Items, item)
			}
		}
		if err = itemRows.Err(); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}
	return orders, nil
}
