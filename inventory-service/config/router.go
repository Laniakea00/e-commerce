package config

import (
	"database/sql"
	"fmt"
	invpb "github.com/Laniakea00/e-commerce/proto/inventory"
	handler "inventory-service/handler/grpc"
	"inventory-service/repository"
	"inventory-service/usecase"

	"google.golang.org/grpc"
	"net"
)

func SetupRouter(db *sql.DB) error {
	productRepo := repository.NewProductRepository(db)
	productUC := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewInventoryHandler(productUC)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	server := grpc.NewServer()
	invpb.RegisterInventoryServiceServer(server, productHandler)

	fmt.Println("✅ Inventory gRPC server running on :50052")
	return server.Serve(listener)
}
