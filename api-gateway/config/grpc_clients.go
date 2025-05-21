package config

import (
	"log"

	"google.golang.org/grpc"

	inventorypb "github.com/Laniakea00/e-commerce/proto/inventory"
	orderpb "github.com/Laniakea00/e-commerce/proto/order"
	userpb "github.com/Laniakea00/e-commerce/proto/user"
)

type Clients struct {
	UserClient      userpb.UserServiceClient
	InventoryClient inventorypb.InventoryServiceClient
	OrderClient     orderpb.OrderServiceClient
}

func InitClients() *Clients {
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to User service: %v", err)
	}

	inventoryConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Inventory service: %v", err)
	}

	orderConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Order service: %v", err)
	}

	return &Clients{
		UserClient:      userpb.NewUserServiceClient(userConn),
		InventoryClient: inventorypb.NewInventoryServiceClient(inventoryConn),
		OrderClient:     orderpb.NewOrderServiceClient(orderConn),
	}
}
