package app

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"go_store/config"
	"go_store/db"
	"go_store/generated/proto/admin"
	"go_store/generated/proto/order"
	"go_store/generated/proto/product"
	controller "go_store/internal/controller/grpc"
	"go_store/internal/controller/interceptor"
	"go_store/internal/repository"
	"go_store/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const sleepDuration = 3

func Run(logger *zap.Logger, cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.PG.URL)

	if err != nil {
		logger.Error("can not create pgxpool", zap.Error(err))
		return
	}

	defer dbPool.Close()

	db.SetupPostgres(dbPool, logger)

	productRepository := repository.NewProductRepository(dbPool)
	orderRepository := repository.NewOrderRepository(dbPool)

	productUseCase := usecase.NewProductUseCase(logger, productRepository)
	orderUseCase := usecase.NewOrderUseCase(logger, orderRepository)
	adminUseCase := usecase.NewAdminUseCase(logger, &cfg.Admin)

	ctrl := controller.New(logger, productUseCase, orderUseCase, adminUseCase)
	go runGrpc(cfg, logger, ctrl)

	<-ctx.Done()
	time.Sleep(time.Second * sleepDuration)
}

func runGrpc(cfg *config.Config, logger *zap.Logger, server controller.Server) {
	port := ":" + cfg.GRPC.Port
	lis, err := net.Listen("tcp", port)

	if err != nil {
		logger.Error("can not open tcp socket", zap.Error(err))
		os.Exit(-1)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.AuthInterceptor(cfg.Admin.JWTSecret)))
	reflection.Register(s)

	product.RegisterProductServiceServer(s, server)
	order.RegisterOrderServiceServer(s, server)
	admin.RegisterAdminServiceServer(s, server)

	logger.Info("grpc server listening at port", zap.String("port", port))

	if err = s.Serve(lis); err != nil {
		logger.Error("grpc server listen error", zap.Error(err))
	}
}
