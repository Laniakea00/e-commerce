package usecase

import (
	"order-service/domain"
	"order-service/repository"
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
	repo repository.OrderRepository
}

func NewOrderUsecase(r repository.OrderRepository) OrderUsecase {
	return &orderUsecase{repo: r}
}

func (u *orderUsecase) Create(order *domain.Order) error {
	order.Status = "pending"
	return u.repo.Create(order)
}

func (u *orderUsecase) Get(id int) (*domain.Order, error) {
	return u.repo.GetByID(id)
}

func (u *orderUsecase) UpdateStatus(id int, status string) error {
	return u.repo.UpdateStatus(id, status)
}

func (u *orderUsecase) ListByUser(userID int) ([]*domain.Order, error) {
	return u.repo.ListByUser(userID)
}

func (u *orderUsecase) UpdateOrder(order *domain.Order) error {
	return u.repo.UpdateOrder(order)
}

func (u *orderUsecase) DeleteByID(id int) error {
	return u.repo.DeleteByID(id)
}
