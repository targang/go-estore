package usecase

import (
	"context"
	"go_store/internal/model"
)

type ProductUseCase interface {
	Create(ctx context.Context, name string, description string, price int64) (string, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (*model.Product, error)
	List(ctx context.Context, limit, offset int32) ([]model.Product, error)
}

type OrderUseCase interface {
	Create(ctx context.Context, customerName string, customerEmail string, items []model.OrderItem) (string, error)
	Get(ctx context.Context, id string) (*model.Order, error)
	UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error
	List(ctx context.Context, limit, offset int32) ([]model.Order, error)
}
