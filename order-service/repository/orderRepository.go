package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"order-service/domain"
)

type OrderRepository interface {
	Create(order *domain.Order) error
	GetByID(id int) (*domain.Order, error)
	UpdateStatus(id int, status string) error
	ListByUser(userID int) ([]*domain.Order, error)
	UpdateOrder(order *domain.Order) error
	DeleteByID(id int) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *domain.Order) error {
	itemsJSON, _ := json.Marshal(order.Items)
	_, err := r.db.Exec(
		"INSERT INTO orders (user_id, status, items) VALUES (?, ?, ?)",
		order.UserID, order.Status, string(itemsJSON),
	)
	return err
}

func (r *orderRepository) GetByID(id int) (*domain.Order, error) {
	row := r.db.QueryRow("SELECT id, user_id, status, items FROM orders WHERE id = ?", id)
	var o domain.Order
	var itemsJSON string
	err := row.Scan(&o.ID, &o.UserID, &o.Status, &itemsJSON)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(itemsJSON), &o.Items)
	return &o, nil
}

func (r *orderRepository) UpdateStatus(id int, status string) error {
	res, err := r.db.Exec("UPDATE orders SET status = ? WHERE id = ?", status, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *orderRepository) ListByUser(userID int) ([]*domain.Order, error) {
	rows, err := r.db.Query("SELECT id, user_id, status, items FROM orders WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var o domain.Order
		var itemsJSON string
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &itemsJSON); err != nil {
			continue
		}
		json.Unmarshal([]byte(itemsJSON), &o.Items)
		orders = append(orders, &o)
	}
	return orders, nil
}

func (r *orderRepository) UpdateOrder(order *domain.Order) error {
	itemsJSON, _ := json.Marshal(order.Items)
	res, err := r.db.Exec("UPDATE orders SET user_id = ?, status = ?, items = ? WHERE id = ?",
		order.UserID, order.Status, string(itemsJSON), order.ID,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (r *orderRepository) DeleteByID(id int) error {
	_, err := r.db.Exec("DELETE FROM orders WHERE id = ?", id)
	return err
}
