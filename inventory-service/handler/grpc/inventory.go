package grpc

import (
	"context"
	invpb "github.com/Laniakea00/e-commerce/proto/inventory"
	"inventory-service/domain"
	"inventory-service/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryHandler struct {
	invpb.UnimplementedInventoryServiceServer
	uc usecase.ProductUsecase
}

func NewInventoryHandler(uc usecase.ProductUsecase) *InventoryHandler {
	return &InventoryHandler{uc: uc}
}

func (h *InventoryHandler) CreateProduct(ctx context.Context, req *invpb.Product) (*invpb.ProductResponse, error) {
	product := domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		Stock:       int(req.Stock),
		CategoryID:  int(req.CategoryId),
		Size:        req.Size,
		Color:       req.Color,
		Gender:      req.Gender,
		Material:    req.Material,
		Season:      req.Season,
	}

	err := h.uc.Create(product)
	if err != nil {
		return &invpb.ProductResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// В ответ добавим объект product
	return &invpb.ProductResponse{
		Success: true,
		Message: "Product created successfully",
		Product: &invpb.Product{
			Id:          int32(product.ID), // если ID доступен
			Name:        product.Name,
			Description: product.Description,
			Price:       float32(product.Price),
			Stock:       int32(product.Stock),
			CategoryId:  int32(product.CategoryID),
			Size:        product.Size,
			Color:       product.Color,
			Gender:      product.Gender,
			Material:    product.Material,
			Season:      product.Season,
		},
	}, nil
}

func (h *InventoryHandler) GetProductByID(ctx context.Context, req *invpb.ProductID) (*invpb.Product, error) {
	product, err := h.uc.GetByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Product not found: %v", err)
	}

	return &invpb.Product{
		Id:          int32(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Price:       float32(product.Price),
		Stock:       int32(product.Stock),
		CategoryId:  int32(product.CategoryID),
		Size:        product.Size,
		Color:       product.Color,
		Gender:      product.Gender,
		Material:    product.Material,
		Season:      product.Season,
	}, nil
}

func (h *InventoryHandler) UpdateProduct(ctx context.Context, req *invpb.Product) (*invpb.ProductResponse, error) {
	product := domain.Product{
		ID:          int(req.Id),
		Name:        req.Name,
		Description: req.Description,
		Price:       float64(req.Price),
		Stock:       int(req.Stock),
		CategoryID:  int(req.CategoryId),
		Size:        req.Size,
		Color:       req.Color,
		Gender:      req.Gender,
		Material:    req.Material,
		Season:      req.Season,
	}

	err := h.uc.Update(product)
	if err != nil {
		return &invpb.ProductResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &invpb.ProductResponse{
		Success: true,
		Message: "Product updated successfully",
		Product: req,
	}, nil
}

func (h *InventoryHandler) DeleteProductByID(ctx context.Context, req *invpb.ProductID) (*invpb.DeleteResponse, error) {
	err := h.uc.DeleteByID(int(req.Id))
	if err != nil {
		return &invpb.DeleteResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &invpb.DeleteResponse{
		Success: true,
		Message: "Product deleted successfully",
	}, nil
}

func (h *InventoryHandler) ListProducts(ctx context.Context, _ *invpb.Empty) (*invpb.ProductList, error) {
	products, err := h.uc.GetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list products: %v", err)
	}

	var pbProducts []*invpb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &invpb.Product{
			Id:          int32(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Price:       float32(p.Price),
			Stock:       int32(p.Stock),
			CategoryId:  int32(p.CategoryID),
			Size:        p.Size,
			Color:       p.Color,
			Gender:      p.Gender,
			Material:    p.Material,
			Season:      p.Season,
		})
	}

	return &invpb.ProductList{Products: pbProducts}, nil
}

func (h *InventoryHandler) DecreaseStock(ctx context.Context, req *invpb.DecreaseStockRequest) (*invpb.ProductResponse, error) {
	err := h.uc.DecreaseStock(int(req.ProductId), int(req.Quantity))
	if err != nil {
		return &invpb.ProductResponse{
			Success: false,
			Message: "Failed to decrease stock: " + err.Error(),
		}, nil
	}

	return &invpb.ProductResponse{
		Success: true,
		Message: "Stock decreased successfully",
	}, nil
}
