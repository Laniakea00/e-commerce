package usecase

import (
	"consumer-service/domain"
	"context"
	"log"

	invpb "github.com/Laniakea00/e-commerce/proto/inventory"
	"google.golang.org/grpc"
)

type StockReducer interface {
	ProcessOrder(order domain.Order) error
}

type stockReducer struct {
	inventoryClient invpb.InventoryServiceClient
}

func NewStockReducer(conn *grpc.ClientConn) StockReducer {
	client := invpb.NewInventoryServiceClient(conn)
	return &stockReducer{inventoryClient: client}
}

func (s *stockReducer) ProcessOrder(order domain.Order) error {
	for _, item := range order.Items {
		req := &invpb.DecreaseStockRequest{
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
		}

		_, err := s.inventoryClient.DecreaseStock(context.Background(), req)
		if err != nil {
			log.Printf("❌ Failed to decrease stock for product %d: %v", item.ProductID, err)
			continue
		}
		log.Printf("🛒 Decreased stock for product %d by %d", item.ProductID, item.Quantity)
	}
	return nil
}
