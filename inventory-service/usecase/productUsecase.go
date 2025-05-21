package usecase

import (
	"encoding/json"
	"fmt"
	"inventory-service/domain"
	"inventory-service/redis"
	"inventory-service/repository"
	"log"
	"time"
)

type ProductUsecase interface {
	Create(product domain.Product) error
	GetByID(id int) (domain.Product, error)
	Update(product domain.Product) error
	DeleteByID(id int) error
	GetAll() ([]domain.Product, error)
	DecreaseStock(productId int, quantity int) error
}

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (u *productUsecase) Create(product domain.Product) error {
	return u.repo.Create(product)
}

func (u *productUsecase) GetByID(id int) (domain.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", id)

	// 🔍 Проверка кэша
	log.Println("🔍 Checking Redis for product", id)
	cached, err := redis.RedisClient.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		var product domain.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			log.Println("💾 Found in Redis")
			return product, nil
		}
	}

	// 📦 Получение из базы
	product, err := u.repo.GetByID(id)
	if err != nil {
		return domain.Product{}, err
	}

	// 💾 Кэшируем
	data, err := json.Marshal(product)
	if err == nil {
		err := redis.RedisClient.Set(redis.Ctx, cacheKey, data, 10*time.Minute).Err()
		if err == nil {
			log.Println("📦 Fetched from DB and cached")
		} else {
			log.Println("⚠️ Failed to cache product:", err)
		}
	}

	return product, nil
}

func (u *productUsecase) Update(product domain.Product) error {
	log.Printf("✏️ Updating product ID %d", product.ID)

	err := u.repo.Update(product)
	if err != nil {
		log.Printf("❌ Failed to update product in DB: %v", err)
		return err
	}

	// 🧹 Инвалидация кэша
	cacheKey := fmt.Sprintf("product:%d", product.ID)
	delErr := redis.RedisClient.Del(redis.Ctx, cacheKey).Err()
	if delErr != nil {
		log.Printf("⚠️ Failed to delete cache for product ID %d: %v", product.ID, delErr)
	} else {
		log.Printf("🧹 Cache invalidated for product ID %d", product.ID)
	}

	return nil
}

func (u *productUsecase) DeleteByID(id int) error {
	log.Printf("🗑️ Deleting product ID %d", id)

	err := u.repo.DeleteByID(id)
	if err != nil {
		log.Printf("❌ Failed to delete product from DB: %v", err)
		return err
	}

	// 🧹 Инвалидация кэша
	cacheKey := fmt.Sprintf("product:%d", id)
	delErr := redis.RedisClient.Del(redis.Ctx, cacheKey).Err()
	if delErr != nil {
		log.Printf("⚠️ Failed to delete cache for product ID %d: %v", id, delErr)
	} else {
		log.Printf("🧹 Cache invalidated for product ID %d", id)
	}

	return nil
}

func (u *productUsecase) GetAll() ([]domain.Product, error) {
	return u.repo.GetAll()
}

func (u *productUsecase) DecreaseStock(productId int, quantity int) error {
	return u.repo.DecreaseStock(productId, quantity)
}
