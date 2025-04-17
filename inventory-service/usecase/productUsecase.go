package usecase

import (
	"inventory-service/domain"
	"inventory-service/repository"
)

type ProductUsecase interface {
	Create(product domain.Product) error
	GetByID(id int) (domain.Product, error)
	Update(product domain.Product) error
	DeleteByID(id int) error
	GetAll() ([]domain.Product, error)
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
	return u.repo.GetByID(id)
}

func (u *productUsecase) Update(product domain.Product) error {
	return u.repo.Update(product)
}

func (u *productUsecase) DeleteByID(id int) error {
	return u.repo.DeleteByID(id)
}

func (u *productUsecase) GetAll() ([]domain.Product, error) {
	return u.repo.GetAll()
}
