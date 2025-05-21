package config

import (
	"database/sql"
	"fmt"
	userpb "github.com/Laniakea00/e-commerce/proto/user"
	"google.golang.org/grpc"
	"net"
	handler "user-service/handler/grpc"
	"user-service/repository"
	"user-service/usecase"
)

func SetupRouter(db *sql.DB) error {
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	server := grpc.NewServer()
	userpb.RegisterUserServiceServer(server, userHandler)

	fmt.Println("✅ User gRPC server running on :50051")
	return server.Serve(listener)
}
