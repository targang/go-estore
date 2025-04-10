package grpc

import (
	"context"
	"go.uber.org/zap"
	"go_store/generated/proto/admin"
	"go_store/generated/proto/common"
	"go_store/generated/proto/order"
	"go_store/generated/proto/product"
	"go_store/internal/model"
	"go_store/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ Server = (*Implementation)(nil)

type Server interface {
	product.ProductServiceServer
	order.OrderServiceServer
	admin.AdminServiceServer
}

type Implementation struct {
	logger         *zap.Logger
	productUseCase usecase.ProductUseCase
	orderUseCase   usecase.OrderUseCase
	adminUseCase   usecase.AdminUseCase
}

func (i *Implementation) Login(ctx context.Context, request *admin.AdminLoginRequest) (*admin.AdminLoginResponse, error) {
	if err := request.ValidateAll(); err != nil {
		i.logger.Warn("validation error", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	token, err := i.adminUseCase.Login(ctx, request.Username, request.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &admin.AdminLoginResponse{Token: token}, nil
}

func (i *Implementation) GetProduct(ctx context.Context, request *product.GetProductRequest) (*product.GetProductResponse, error) {
	if err := request.Validate(); err != nil {
		i.logger.Warn("validation error", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := i.productUseCase.Get(ctx, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}
	return &product.GetProductResponse{Product: &common.Product{
		Id:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		Price:       result.Price,
	}}, nil
}

func (i *Implementation) ListProducts(ctx context.Context, request *product.ListProductsRequest) (*product.ListProductsResponse, error) {
	if err := request.ValidateAll(); err != nil {
		i.logger.Warn("validation error", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := i.productUseCase.List(ctx, request.Limit, request.Offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err)
	}
	products := make([]*common.Product, 0, len(result))
	for _, p := range result {
		products = append(products, &common.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return &product.ListProductsResponse{Products: products}, nil
}

func (i *Implementation) CreateProduct(ctx context.Context, request *admin.CreateProductRequest) (*admin.CreateProductResponse, error) {
	if err := request.ValidateAll(); err != nil {
		i.logger.Warn("validation error", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := i.productUseCase.Create(ctx, request.Name, request.Description, request.Price)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &admin.CreateProductResponse{Id: result}, nil
}

func (i *Implementation) DeleteProduct(ctx context.Context, request *admin.DeleteProductRequest) (*admin.DeleteProductResponse, error) {
	if err := request.ValidateAll(); err != nil {
		i.logger.Warn("validation error", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := i.productUseCase.Delete(ctx, request.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &admin.DeleteProductResponse{}, nil
}

func (i *Implementation) CreateOrder(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	if err := request.ValidateAll(); err != nil {
		i.logger.Warn("validation error", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	items := make([]model.OrderItem, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, model.OrderItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}
	result, err := i.orderUseCase.Create(ctx, request.CustomerName, request.CustomerEmail, items)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &order.CreateOrderResponse{Id: result}, nil
}

func (i *Implementation) GetOrder(ctx context.Context, request *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	if err := request.ValidateAll(); err != nil {
		i.logger.Warn("validation error", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	result, err := i.orderUseCase.Get(ctx, request.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &order.GetOrderResponse{Order: result.ConvertToMessage()}, nil
}

func (i *Implementation) ListOrders(ctx context.Context, request *admin.ListOrdersRequest) (*admin.ListOrdersResponse, error) {
	if err := request.ValidateAll(); err != nil {
		i.logger.Warn("validation error", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	modelOrders, err := i.orderUseCase.List(ctx, request.Limit, request.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	responseOrders := make([]*common.Order, 0, len(modelOrders))
	for _, modelOrder := range modelOrders {
		responseOrders = append(responseOrders, modelOrder.ConvertToMessage())
	}
	return &admin.ListOrdersResponse{Orders: responseOrders}, nil
}

func (i *Implementation) UpdateOrderStatus(ctx context.Context, request *admin.UpdateOrderStatusRequest) (*admin.UpdateOrderStatusResponse, error) {
	if err := request.ValidateAll(); err != nil {
		i.logger.Warn("validation error", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := i.orderUseCase.UpdateStatus(ctx, request.Id, model.OrderStatus(request.Status))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &admin.UpdateOrderStatusResponse{}, nil
}

func New(
	logger *zap.Logger,
	productUseCase usecase.ProductUseCase,
	orderUseCase usecase.OrderUseCase,
	adminUseCase usecase.AdminUseCase,
) *Implementation {
	return &Implementation{
		logger:         logger,
		productUseCase: productUseCase,
		orderUseCase:   orderUseCase,
		adminUseCase:   adminUseCase,
	}
}
