package repository

import (
	"context"
	"go_store/internal/model"
)

// ProductRepository определяет контракт для работы с товарами в хранилище данных.
type ProductRepository interface {
	// Create сохраняет новый товар в базе данных.
	// Возвращает ID созданного товара.
	Create(ctx context.Context, product *model.Product) (string, error)

	// GetByID возвращает товар по его идентификатору.
	// Возвращает ошибку, если товар не найден.
	GetByID(ctx context.Context, id string) (*model.Product, error)

	// Delete удаляет товар по его идентификатору.
	// Возвращает ошибку, если товар не найден или удаление не удалось.
	Delete(ctx context.Context, id string) error

	// List возвращает список товаров с поддержкой пагинации.
	// Параметры limit и offset управляют количеством возвращаемых результатов.
	List(ctx context.Context, limit, offset int32) ([]model.Product, error)
}

// OrderRepository определяет контракт для работы с заказами в хранилище данных.
type OrderRepository interface {
	// Create сохраняет новый заказ и связанные с ним позиции.
	// Возвращает ID созданного товара.
	Create(ctx context.Context, order *model.Order) (string, error)

	// GetByID возвращает заказ по его идентификатору, включая связанные позиции.
	// Возвращает ошибку, если заказ не найден.
	GetByID(ctx context.Context, id string) (*model.Order, error)

	// UpdateStatus обновляет статус заказа и поле updated_at.
	// Возвращает ошибку.
	UpdateStatus(ctx context.Context, order *model.Order) error

	// Delete удаляет заказ и все связанные с ним позиции по идентификатору.
	// Возвращает ошибку, если удаление не удалось.
	Delete(ctx context.Context, id string) error

	// List возвращает список заказов с поддержкой пагинации.
	// В каждый заказ включаются связанные позиции.
	List(ctx context.Context, limit, offset int32) ([]model.Order, error)
}
