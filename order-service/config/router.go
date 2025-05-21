package config

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"order-service/infrastructure/messaging"

	"google.golang.org/grpc"

	handler "order-service/handler/grpc"
	"order-service/repository"
	"order-service/usecase"

	orderpb "github.com/Laniakea00/e-commerce/proto/order"
)

func SetupRouter(db *sql.DB) error {
	producer, err := messaging.NewNATSProducer("nats://localhost:4222")
	if err != nil {
		log.Fatalf("NATS error: %v", err)
	}

	orderRepo := repository.NewOrderRepository(db)
	orderUseCase := usecase.NewOrderUsecase(orderRepo, producer)
	orderHandler := handler.NewOrderServer(orderUseCase)

	// Запуск gRPC сервера
	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	server := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(server, orderHandler)

	fmt.Println("✅ Order gRPC server running on :50053")
	return server.Serve(listener)
}
