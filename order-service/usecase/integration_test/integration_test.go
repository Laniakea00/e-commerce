package integration_test

import (
	"database/sql"
	"encoding/json"
	"order-service/domain"
	"order-service/redis"
	"order-service/repository"
	"order-service/usecase"
	"strconv"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

var (
	db        *sql.DB
	orderRepo repository.OrderRepository
	orderUC   usecase.OrderUsecase
)

type mockProducer struct{}

func (m *mockProducer) PublishOrderCreated(order domain.Order) error {
	return nil
}

func setupTestDB(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite", "orders_test.db")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create the orders table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS orders (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            user_id INTEGER NOT NULL,
            status TEXT NOT NULL,
            items TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		t.Fatalf("Failed to create orders table: %v", err)
	}

	// Create the order_items table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS order_items (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            order_id INTEGER NOT NULL,
            product_id INTEGER NOT NULL,
            quantity INTEGER NOT NULL,
            FOREIGN KEY(order_id) REFERENCES orders(id)
        )
    `)
	if err != nil {
		t.Fatalf("Failed to create order_items table: %v", err)
	}

	// Clear existing data
	_, err = db.Exec("DELETE FROM order_items")
	if err != nil {
		t.Fatalf("Failed to clear order_items table: %v", err)
	}
	_, err = db.Exec("DELETE FROM orders")
	if err != nil {
		t.Fatalf("Failed to clear orders table: %v", err)
	}

	orderRepo = repository.NewOrderRepository(db)
	producer := &mockProducer{}
	orderUC = usecase.NewOrderUsecase(orderRepo, producer)

	// Initialize Redis
	redis.InitRedis()
}

func cleanup(t *testing.T) {
	if db != nil {
		db.Close()
	}
	if redis.RedisClient != nil {
		redis.RedisClient.FlushAll(redis.Ctx)
		redis.RedisClient.Close()
	}
}

func TestOrderIntegration_CreateAndGet(t *testing.T) {
	setupTestDB(t)
	defer cleanup(t)

	// Create test order
	order := &domain.Order{
		UserID: 1,
		Items: []domain.OrderItem{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	}

	// Test Create
	err := orderUC.Create(order)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	if order.ID == 0 {
		t.Error("Expected order ID to be set after creation")
	}

	// Test Get
	fetchedOrder, err := orderUC.Get(order.ID)
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}

	if fetchedOrder.ID != order.ID {
		t.Errorf("Expected order ID %d, got %d", order.ID, fetchedOrder.ID)
	}
	if fetchedOrder.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", fetchedOrder.Status)
	}

	// Verify Redis cache
	time.Sleep(100 * time.Millisecond) // Wait for async cache operations
	cacheKey := "order:" + strconv.Itoa(order.ID)
	cached, err := redis.RedisClient.Get(redis.Ctx, cacheKey).Result()
	if err != nil {
		t.Fatalf("Failed to get from cache: %v", err)
	}

	var cachedOrder domain.Order
	err = json.Unmarshal([]byte(cached), &cachedOrder)
	if err != nil {
		t.Fatalf("Failed to unmarshal cached order: %v", err)
	}
	if cachedOrder.ID != order.ID {
		t.Errorf("Cached order ID mismatch: expected %d, got %d", order.ID, cachedOrder.ID)
	}
}

func TestOrderIntegration_UpdateAndCache(t *testing.T) {
	setupTestDB(t)
	defer cleanup(t)

	// Create initial order
	order := &domain.Order{
		UserID: 1,
		Items: []domain.OrderItem{
			{ProductID: 1, Quantity: 2},
		},
	}
	err := orderUC.Create(order)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	// Update status
	err = orderUC.UpdateStatus(order.ID, "shipped")
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	// Verify cache invalidation
	time.Sleep(100 * time.Millisecond)
	cacheKey := "order:" + strconv.Itoa(order.ID)
	_, err = redis.RedisClient.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		t.Error("Expected cache to be invalidated after update")
	}

	// Fetch updated order
	updated, err := orderUC.Get(order.ID)
	if err != nil {
		t.Fatalf("Failed to get updated order: %v", err)
	}
	if updated.Status != "shipped" {
		t.Errorf("Expected status 'shipped', got '%s'", updated.Status)
	}
}

func TestOrderIntegration_ListByUser(t *testing.T) {
	setupTestDB(t)
	defer cleanup(t)

	// Create multiple orders for the same user
	userID := 1
	for i := 0; i < 3; i++ {
		order := &domain.Order{
			UserID: userID,
			Items: []domain.OrderItem{
				{ProductID: 1, Quantity: 1},
			},
		}
		err := orderUC.Create(order)
		if err != nil {
			t.Fatalf("Failed to create order %d: %v", i, err)
		}
	}

	// Test ListByUser
	orders, err := orderUC.ListByUser(userID)
	if err != nil {
		t.Fatalf("Failed to list orders: %v", err)
	}
	if len(orders) != 3 {
		t.Errorf("Expected 3 orders, got %d", len(orders))
	}

	for _, o := range orders {
		if o.UserID != userID {
			t.Errorf("Expected userID %d, got %d", userID, o.UserID)
		}
	}
}
