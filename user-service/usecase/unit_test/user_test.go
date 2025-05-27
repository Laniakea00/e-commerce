package unit_test

import (
	"testing"
	"user-service/domain"
	"user-service/usecase"
)

type mockUserRepo struct {
	getByEmailUser *domain.User
	getByEmailErr  error
}

func (m *mockUserRepo) ExistsByEmail(email string) (bool, error) {
	return false, nil
}

func (m *mockUserRepo) Create(user *domain.User) error {
	return nil
}

func (m *mockUserRepo) GetByEmail(email string) (*domain.User, error) {
	return m.getByEmailUser, m.getByEmailErr
}

func (m *mockUserRepo) GetByID(id int) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserRepo) UpdateProfile(user *domain.User) error {
	return nil
}

func (m *mockUserRepo) DeleteByID(id int) error {
	return nil
}

func (m *mockUserRepo) GetAll() ([]*domain.User, error) {
	return nil, nil
}

func TestUserUsecase_Authenticate_Success(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailUser: &domain.User{Email: "test@example.com", Password: "password123"},
	}
	uc := usecase.NewUserUsecase(repo)

	user, err := uc.Authenticate("test@example.com", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got %s", user.Email)
	}
}

func TestUserUsecase_Authenticate_InvalidCredentials(t *testing.T) {
	repo := &mockUserRepo{
		getByEmailUser: &domain.User{Email: "test@example.com", Password: "password123"},
	}
	uc := usecase.NewUserUsecase(repo)

	_, err := uc.Authenticate("test@example.com", "wrongpassword")
	if err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected 'invalid credentials' error, got %v", err)
	}
}
