package usecase

import (
	"errors"
	"user-service/domain"
	"user-service/repositiry"
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
	return u.repo.GetByID(id)
}

func (u *userUsecase) DeleteUser(id int) error {
	return u.repo.DeleteByID(id)
}

func (u *userUsecase) UpdateProfile(user *domain.User) error {
	return u.repo.UpdateProfile(user)
}

func (u *userUsecase) ListUsers() ([]*domain.User, error) {
	return u.repo.GetAll()
}
