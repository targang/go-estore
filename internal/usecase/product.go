package usecase

import (
	"context"
	"go.uber.org/zap"
	"go_store/internal/model"
	"go_store/internal/repository"
)

var _ ProductUseCase = (*productUseCaseImpl)(nil)

type productUseCaseImpl struct {
	logger            *zap.Logger
	productRepository repository.ProductRepository
}

func NewProductUseCase(logger *zap.Logger, productRepository repository.ProductRepository) ProductUseCase {
	return &productUseCaseImpl{
		logger:            logger,
		productRepository: productRepository,
	}
}

func (p *productUseCaseImpl) Create(ctx context.Context, name string, description string, price int64) (string, error) {
	return p.productRepository.Create(ctx, &model.Product{Name: name, Description: description, Price: price})
}

func (p *productUseCaseImpl) Delete(ctx context.Context, id string) error {
	return p.productRepository.Delete(ctx, id)
}

func (p *productUseCaseImpl) Get(ctx context.Context, id string) (*model.Product, error) {
	return p.productRepository.GetByID(ctx, id)
}

func (p *productUseCaseImpl) List(ctx context.Context, limit, offset int32) ([]model.Product, error) {
	return p.productRepository.List(ctx, limit, offset)
}
