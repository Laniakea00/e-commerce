package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"user-service/domain"
	"user-service/redis"
	"user-service/repository"
)

type UserUsecase interface {
	Register(user *domain.User) error
	Authenticate(email, password string) (*domain.User, error)
	GetProfile(id int) (*domain.User, error)
	UpdateProfile(user *domain.User) error
	DeleteUser(id int) error
	ListUsers() ([]*domain.User, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecase{repo: r}
}

func (u *userUsecase) Register(user *domain.User) error {
	exists, err := u.repo.ExistsByEmail(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}
	return u.repo.Create(user)
}

func (u *userUsecase) Authenticate(email, password string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil || user.Password != password {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (u *userUsecase) GetProfile(id int) (*domain.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// 🔍 Проверка кэша
	log.Println("🔍 Checking Redis for order", id)
	cached, err := redis.RedisClient.Get(redis.Ctx, cacheKey).Result()
	if err == nil {
		var user domain.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			log.Println("💾 Found in Redis")
			return &user, nil
		}
	}

	// 📦 Получаем из базы
	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 💾 Кэшируем
	data, err := json.Marshal(user)
	if err == nil {
		err := redis.RedisClient.Set(redis.Ctx, cacheKey, data, 10*time.Minute).Err()
		if err == nil {
			log.Println("📦 Fetched from DB and cached")
		} else {
			log.Println("⚠️ Failed to cache user:", err)
		}
	}

	return user, nil
}

func (u *userUsecase) UpdateProfile(user *domain.User) error {
	log.Printf("✏️ Updating profile for user ID %d", user.ID)

	err := u.repo.UpdateProfile(user)
	if err != nil {
		log.Printf("❌ Failed to update user in DB: %v", err)
		return err
	}

	// 🧹 Инвалидация кэша
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	delErr := redis.RedisClient.Del(redis.Ctx, cacheKey).Err()
	if delErr != nil {
		log.Printf("⚠️ Failed to delete cache for user ID %d: %v", user.ID, delErr)
	} else {
		log.Printf("🧹 Cache invalidated for user ID %d", user.ID)
	}

	return nil
}

func (u *userUsecase) DeleteUser(id int) error {
	log.Printf("🗑️ Deleting user ID %d", id)

	err := u.repo.DeleteByID(id)
	if err != nil {
		log.Printf("❌ Failed to delete user from DB: %v", err)
		return err
	}

	// 🧹 Инвалидация кэша
	cacheKey := fmt.Sprintf("user:%d", id)
	delErr := redis.RedisClient.Del(redis.Ctx, cacheKey).Err()
	if delErr != nil {
		log.Printf("⚠️ Failed to delete cache for user ID %d: %v", id, delErr)
	} else {
		log.Printf("🧹 Cache invalidated for user ID %d", id)
	}

	return nil
}

func (u *userUsecase) ListUsers() ([]*domain.User, error) {
	return u.repo.GetAll()
}
