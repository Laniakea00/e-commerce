package unit_test

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"inventory-service/domain"
	redisLocal "inventory-service/redis"
	"inventory-service/usecase"
	"testing"
	"time"
)

type mockProductRepo struct {
	createErr   error
	getByIDProd domain.Product
	getByIDErr  error
	updateErr   error
	deleteErr   error
}

func (m *mockProductRepo) Create(product domain.Product) error {
	return m.createErr
}

func (m *mockProductRepo) GetByID(id int) (domain.Product, error) {
	return m.getByIDProd, m.getByIDErr
}

func (m *mockProductRepo) Update(product domain.Product) error {
	return m.updateErr
}

func (m *mockProductRepo) DeleteByID(id int) error {
	return m.deleteErr
}

func (m *mockProductRepo) GetAll() ([]domain.Product, error) {
	return nil, nil
}

func (m *mockProductRepo) DecreaseStock(productId int, quantity int) error {
	return nil
}

func NewMockRedisClient() *redis.Client {
	opt := &redis.Options{
		Addr: "",
	}
	client := redis.NewClient(opt)
	return client
}

func TestProductUsecase_Create(t *testing.T) {
	repo := &mockProductRepo{}
	uc := usecase.NewProductUsecase(repo)

	err := uc.Create(domain.Product{ID: 1, Name: "Test Product"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProductUsecase_GetByID_CacheHit(t *testing.T) {
	repo := &mockProductRepo{}
	mockRedis := NewMockRedisClient()
	redisLocal.RedisClient = mockRedis
	redisLocal.Ctx = context.Background()

	product := domain.Product{ID: 1, Name: "Cached Product"}
	data, _ := json.Marshal(product)
	mockRedis.Set(redisLocal.Ctx, "product:1", string(data), 10*time.Minute)

	uc := usecase.NewProductUsecase(repo)
	got, err := uc.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != 1 {
		t.Errorf("expected product ID 1, got %d", got.ID)
	}
}

func TestProductUsecase_GetByID_DBHit(t *testing.T) {
	repo := &mockProductRepo{getByIDProd: domain.Product{ID: 1, Name: "DB Product"}}
	mockRedis := NewMockRedisClient()
	redisLocal.RedisClient = mockRedis
	redisLocal.Ctx = context.Background()

	uc := usecase.NewProductUsecase(repo)
	got, err := uc.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != 1 {
		t.Errorf("expected product ID 1, got %d", got.ID)
	}
}

func TestProductUsecase_Update(t *testing.T) {
	repo := &mockProductRepo{}
	mockRedis := NewMockRedisClient()
	redisLocal.RedisClient = mockRedis
	redisLocal.Ctx = context.Background()

	uc := usecase.NewProductUsecase(repo)
	err := uc.Update(domain.Product{ID: 1, Name: "Updated Product"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProductUsecase_DeleteByID(t *testing.T) {
	repo := &mockProductRepo{}
	mockRedis := NewMockRedisClient()
	redisLocal.RedisClient = mockRedis
	redisLocal.Ctx = context.Background()

	uc := usecase.NewProductUsecase(repo)
	err := uc.DeleteByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
