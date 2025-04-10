package repository

import (
	"context"
	"go_store/internal/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) (string, error)

	GetByID(ctx context.Context, id string) (*model.Product, error)

	Delete(ctx context.Context, id string) error

	List(ctx context.Context, limit, offset int32) ([]model.Product, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (string, error)

	GetByID(ctx context.Context, id string) (*model.Order, error)

	UpdateStatus(ctx context.Context, order *model.Order) error

	Delete(ctx context.Context, id string) error

	List(ctx context.Context, limit, offset int32) ([]model.Order, error)
}
