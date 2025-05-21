package repository

import (
	"database/sql"
	_ "errors"
	"user-service/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByEmail(email string) (*domain.User, error)
	GetByID(id int) (*domain.User, error)
	ExistsByEmail(email string) (bool, error)
	GetAll() ([]*domain.User, error)
	DeleteByID(id int) error
	UpdateProfile(user *domain.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	_, err := r.db.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		user.Username, user.Email, user.Password,
	)
	return err
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", email)
	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	row := r.db.QueryRow("SELECT id, username, email, password FROM users WHERE id = ?", id)
	var u domain.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	row := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email)
	var count int
	err := row.Scan(&count)
	return count > 0, err
}

func (r *userRepository) GetAll() ([]*domain.User, error) {
	rows, err := r.db.Query("SELECT id, username, email, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password); err != nil {
			continue
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *userRepository) DeleteByID(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (r *userRepository) UpdateProfile(user *domain.User) error {
	_, err := r.db.Exec("UPDATE users SET username = ?, email = ? WHERE id = ?", user.Username, user.Email, user.ID)
	return err
}
