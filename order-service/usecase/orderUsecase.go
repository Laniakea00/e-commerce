package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"order-service/domain"
	"order-service/redis"
	"order-service/repository"
	"time"
)

type OrderUsecase interface {
	Create(order *domain.Order) error
	Get(id int) (*domain.Order, error)
	UpdateStatus(id int, status string) error
	ListByUser(userID int) ([]*domain.Order, error)
	UpdateOrder(order *domain.Order) error
	DeleteByID(id int) error
}

type orderUsecase struct {
	repo     repository.OrderRepository
	producer OrderProducer
}

type OrderProducer interface {
	PublishOrderCreated(order domain.Order) error
}

func NewOrderUsecase(r repository.OrderRepository, p OrderProducer) OrderUsecase {
	return &orderUsecase{repo: r, producer: p}
}

func (u *orderUsecase) Create(order *domain.Order) error {
	order.Status = "pending"
	if err := u.repo.Create(order); err != nil {
		return err
	}

	if err := u.producer.PublishOrderCreated(*order); err != nil {
		log.Printf("❌ Failed to publish order: %v", err)
	} else {
		log.Printf("✅ Published order %d to NATS", order.ID)
	}
	return nil
}

func (u *orderUsecase) Get(id int) (*domain.Order, error) {
	cacheKey := fmt.Sprintf("order:%d", id)

	// 🔍 Проверка кэша
	log.Println("🔍 Checking Redis for order", id)
	cached, err := redis.RedisClient.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		var order domain.Order
		if err := json.Unmarshal([]byte(cached), &order); err == nil {
			log.Println("💾 Found in Redis")
			return &order, nil
		}
	}

	// 📦 Получаем из базы
	order, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 💾 Кэшируем
	data, err := json.Marshal(order)
	if err == nil {
		err := redis.RedisClient.Set(redis.Ctx, cacheKey, data, 10*time.Minute).Err()
		if err == nil {
			log.Println("📦 Fetched from DB and cached")
		} else {
			log.Println("⚠️ Failed to cache product:", err)
		}
	}

	return order, nil
}

func (u *orderUsecase) ListByUser(userID int) ([]*domain.Order, error) {
	return u.repo.ListByUser(userID)
}

func (u *orderUsecase) UpdateOrder(order *domain.Order) error {
	log.Printf("✏️ Updating order ID %d", order.ID)

	err := u.repo.UpdateOrder(order)
	if err != nil {
		log.Printf("❌ Failed to update order in DB: %v", err)
		return err
	}

	// 🧹 Инвалидация кэша
	cacheKey := fmt.Sprintf("order:%d", order.ID)
	delErr := redis.RedisClient.Del(redis.Ctx, cacheKey).Err()
	if delErr != nil {
		log.Printf("⚠️ Failed to delete cache for order ID %d: %v", order.ID, delErr)
	} else {
		log.Printf("🧹 Cache invalidated for order ID %d", order.ID)
	}

	return nil
}

func (u *orderUsecase) UpdateStatus(id int, status string) error {
	log.Printf("✏️ Updating status of order ID %d to '%s'", id, status)

	err := u.repo.UpdateStatus(id, status)
	if err != nil {
		log.Printf("❌ Failed to update status in DB: %v", err)
		return err
	}

	// 🧹 Инвалидация кэша
	cacheKey := fmt.Sprintf("order:%d", id)
	delErr := redis.RedisClient.Del(redis.Ctx, cacheKey).Err()
	if delErr != nil {
		log.Printf("⚠️ Failed to delete cache for order ID %d: %v", id, delErr)
	} else {
		log.Printf("🧹 Cache invalidated for order ID %d", id)
	}

	return nil
}

func (u *orderUsecase) DeleteByID(id int) error {
	log.Printf("🗑️ Deleting order ID %d", id)

	err := u.repo.DeleteByID(id)
	if err != nil {
		log.Printf("❌ Failed to delete order from DB: %v", err)
		return err
	}

	// 🧹 Инвалидация кэша
	cacheKey := fmt.Sprintf("order:%d", id)
	delErr := redis.RedisClient.Del(redis.Ctx, cacheKey).Err()
	if delErr != nil {
		log.Printf("⚠️ Failed to delete cache for order ID %d: %v", id, delErr)
	} else {
		log.Printf("🧹 Cache invalidated for order ID %d", id)
	}

	return nil
}
