package usecase

import (
	"context"
	"go.uber.org/zap"
	"go_store/internal/model"
	"go_store/internal/repository"
)

var _ OrderUseCase = (*orderUseCaseImpl)(nil)

type orderUseCaseImpl struct {
	logger          *zap.Logger
	orderRepository repository.OrderRepository
}

func NewOrderUseCase(logger *zap.Logger, orderRepository repository.OrderRepository) OrderUseCase {
	return &orderUseCaseImpl{
		logger:          logger,
		orderRepository: orderRepository,
	}
}

func (o *orderUseCaseImpl) Create(ctx context.Context, customerName string, customerEmail string, items []model.OrderItem) (string, error) {
	return o.orderRepository.Create(ctx, &model.Order{
		CustomerName:  customerName,
		CustomerEmail: customerEmail,
		Items:         items,
		Status:        model.UNSPECIFIED,
	})
}

func (o *orderUseCaseImpl) Get(ctx context.Context, id string) (*model.Order, error) {
	return o.orderRepository.GetByID(ctx, id)
}

func (o *orderUseCaseImpl) UpdateStatus(ctx context.Context, id string, status model.OrderStatus) error {
	return o.orderRepository.UpdateStatus(ctx, &model.Order{
		ID:     id,
		Status: status,
	})
}

func (o *orderUseCaseImpl) List(ctx context.Context, limit, offset int32) ([]model.Order, error) {
	return o.orderRepository.List(ctx, limit, offset)
}
