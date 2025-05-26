package usecase_test

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"order-service/domain"
	redisLocal "order-service/redis"
	"order-service/usecase"
	"testing"
	"time"
)

// --- Extend mockRepo for GetByID and ListByUser ---
type mockRepo struct {
	createErr       error
	getByIDOrder    *domain.Order
	getByIDErr      error
	listByUser      []*domain.Order
	listByUserErr   error
	updateStatusErr error
	updateOrderErr  error
	deleteByIDErr   error
}

func (m *mockRepo) Create(order *domain.Order) error {
	return m.createErr
}
func (m *mockRepo) GetByID(id int) (*domain.Order, error) {
	return m.getByIDOrder, m.getByIDErr
}
func (m *mockRepo) ListByUser(userID int) ([]*domain.Order, error) {
	return m.listByUser, m.listByUserErr
}
func (m *mockRepo) UpdateStatus(id int, status string) error {
	return m.updateStatusErr
}
func (m *mockRepo) UpdateOrder(order *domain.Order) error {
	return m.updateOrderErr
}
func (m *mockRepo) DeleteByID(id int) error {
	return m.deleteByIDErr
}

type mockProducer struct{ publishErr error }

func (m *mockProducer) PublishOrderCreated(order domain.Order) error { return m.publishErr }

// --- Mock Redis for Get, UpdateOrder, UpdateStatus, DeleteByID ---
var redisCache = make(map[string]string)

type MockRedisClient struct {
	*redis.Client // Embed the real client to inherit methods
}

func NewMockRedisClient() *redis.Client {
	// Create mock options that won't actually connect
	opt := &redis.Options{
		Addr: "",
	}
	client := redis.NewClient(opt)
	return client
}

// Override only the methods we need for testing
func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	val, ok := redisCache[key]
	if !ok {
		return redis.NewStringResult("", redis.Nil)
	}
	return redis.NewStringResult(val, nil)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	redisCache[key] = value.(string)
	return redis.NewStatusResult("OK", nil)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	count := 0
	for _, key := range keys {
		if _, ok := redisCache[key]; ok {
			delete(redisCache, key)
			count++
		}
	}
	return redis.NewIntResult(int64(count), nil)
}

func init() {
	redisLocal.RedisClient = NewMockRedisClient()
	redisLocal.Ctx = context.Background()
}

// --- Tests ---

func TestOrderUsecase_Get_Success(t *testing.T) {
	order := &domain.Order{ID: 42, Status: "pending"}
	repo := &mockRepo{getByIDOrder: order}
	producer := &mockProducer{}
	uc := usecase.NewOrderUsecase(repo, producer)

	got, err := uc.Get(42)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != 42 {
		t.Errorf("expected ID 42, got %d", got.ID)
	}
}

func TestOrderUsecase_Get_RepoError(t *testing.T) {
	repo := &mockRepo{getByIDErr: errors.New("not found")}
	producer := &mockProducer{}
	uc := usecase.NewOrderUsecase(repo, producer)

	_, err := uc.Get(1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderUsecase_ListByUser_Success(t *testing.T) {
	orders := []*domain.Order{{ID: 1}, {ID: 2}}
	repo := &mockRepo{listByUser: orders}
	producer := &mockProducer{}
	uc := usecase.NewOrderUsecase(repo, producer)

	got, err := uc.ListByUser(123)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 orders, got %d", len(got))
	}
}

func TestOrderUsecase_ListByUser_Error(t *testing.T) {
	repo := &mockRepo{listByUserErr: errors.New("db error")}
	producer := &mockProducer{}
	uc := usecase.NewOrderUsecase(repo, producer)

	_, err := uc.ListByUser(123)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderUsecase_UpdateStatus_Success(t *testing.T) {
	repo := &mockRepo{}
	producer := &mockProducer{}
	uc := usecase.NewOrderUsecase(repo, producer)

	err := uc.UpdateStatus(1, "shipped")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestOrderUsecase_UpdateStatus_Error(t *testing.T) {
	repo := &mockRepo{updateStatusErr: errors.New("update error")}
	producer := &mockProducer{}
	uc := usecase.NewOrderUsecase(repo, producer)

	err := uc.UpdateStatus(1, "shipped")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderUsecase_UpdateOrder_Success(t *testing.T) {
	repo := &mockRepo{}
	producer := &mockProducer{}
	uc := usecase.NewOrderUsecase(repo, producer)

	order := &domain.Order{ID: 1}
	err := uc.UpdateOrder(order)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestOrderUsecase_UpdateOrder_Error(t *testing.T) {
	repo := &mockRepo{updateOrderErr: errors.New("update error")}
	producer := &mockProducer{}
	uc := usecase.NewOrderUsecase(repo, producer)

	order := &domain.Order{ID: 1}
	err := uc.UpdateOrder(order)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOrderUsecase_DeleteByID_Success(t *testing.T) {
	repo := &mockRepo{}
	producer := &mockProducer{}
	uc := usecase.NewOrderUsecase(repo, producer)

	err := uc.DeleteByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestOrderUsecase_DeleteByID_Error(t *testing.T) {
	repo := &mockRepo{deleteByIDErr: errors.New("delete error")}
	producer := &mockProducer{}
	uc := usecase.NewOrderUsecase(repo, producer)

	err := uc.DeleteByID(1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
