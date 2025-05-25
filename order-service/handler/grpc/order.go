package grpc

import (
	"context"
	orderpb "github.com/Laniakea00/e-commerce/proto/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"order-service/domain"
	"order-service/usecase"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	Usecase usecase.OrderUsecase
}

func NewOrderServer(uc usecase.OrderUsecase) *OrderServer {
	return &OrderServer{
		Usecase: uc,
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderpb.OrderRequest) (*orderpb.OrderResponse, error) {
	order := &domain.Order{
		UserID: int(req.UserId),
		Items:  fromProtoItems(req.Items),
	}

	err := s.Usecase.Create(order)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create order: %v", err)
	}

	return &orderpb.OrderResponse{
		Success: true,
		Message: "Order created",
		Order:   toProtoOrder(order),
	}, nil
}

func (s *OrderServer) GetOrderByID(ctx context.Context, req *orderpb.OrderID) (*orderpb.Order, error) {
	order, err := s.Usecase.Get(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Order not found: %v", err)
	}
	return toProtoOrder(order), nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *orderpb.StatusUpdateRequest) (*orderpb.OrderResponse, error) {
	err := s.Usecase.UpdateStatus(int(req.OrderId), req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update status: %v", err)
	}
	order, _ := s.Usecase.Get(int(req.OrderId))
	return &orderpb.OrderResponse{
		Success: true,
		Message: "Status updated",
		Order:   toProtoOrder(order),
	}, nil
}

func (s *OrderServer) ListOrdersByUser(ctx context.Context, req *orderpb.UserID) (*orderpb.OrderList, error) {
	orders, err := s.Usecase.ListByUser(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list orders: %v", err)
	}

	var protoOrders []*orderpb.Order
	for _, o := range orders {
		protoOrders = append(protoOrders, toProtoOrder(o))
	}

	return &orderpb.OrderList{Orders: protoOrders}, nil
}

func (s *OrderServer) UpdateOrder(ctx context.Context, req *orderpb.Order) (*orderpb.OrderResponse, error) {
	order := &domain.Order{
		ID:     int(req.Id),
		UserID: int(req.UserId),
		Status: req.Status,
		Items:  fromProtoItems(req.Items),
	}

	err := s.Usecase.UpdateOrder(order)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update order: %v", err)
	}

	return &orderpb.OrderResponse{
		Success: true,
		Message: "Order updated",
		Order:   toProtoOrder(order),
	}, nil
}

func (s *OrderServer) DeleteOrderByID(ctx context.Context, req *orderpb.OrderID) (*orderpb.DeleteResponse, error) {
	err := s.Usecase.DeleteByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete order: %v", err)
	}
	return &orderpb.DeleteResponse{
		Success: true,
		Message: "Order deleted",
	}, nil
}

// Helpers
func fromProtoItems(protoItems []*orderpb.OrderItem) []domain.OrderItem {
	var items []domain.OrderItem
	for _, item := range protoItems {
		items = append(items, domain.OrderItem{
			ProductID: int(item.ProductId),
			Quantity:  int(item.Quantity),
		})
	}
	return items
}

func toProtoOrder(o *domain.Order) *orderpb.Order {
	var items []*orderpb.OrderItem
	for _, item := range o.Items {
		items = append(items, &orderpb.OrderItem{
			ProductId: int32(item.ProductID),
			Quantity:  int32(item.Quantity),
		})
	}
	return &orderpb.Order{
		Id:     int32(o.ID),
		UserId: int32(o.UserID),
		Status: o.Status,
		Items:  items,
	}
}
