package repository

import (
	"database/sql"
	"inventory-service/domain"
)

type ProductRepository interface {
	Create(product domain.Product) error
	GetByID(id int) (domain.Product, error)
	Update(product domain.Product) error
	DeleteByID(id int) error
	GetAll() ([]domain.Product, error)
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(product domain.Product) error {
	_, err := r.db.Exec("INSERT INTO products (name, category, price, stock) VALUES (?, ?, ?, ?)",
		product.Name, product.Category, product.Price, product.Stock)
	return err
}

func (r *productRepo) GetByID(id int) (domain.Product, error) {
	var p domain.Product
	row := r.db.QueryRow("SELECT id, name, category, price, stock FROM products WHERE id = ?", id)
	err := row.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock)
	return p, err
}

func (r *productRepo) Update(product domain.Product) error {
	_, err := r.db.Exec("UPDATE products SET name=?, category=?, price=?, stock=? WHERE id=?",
		product.Name, product.Category, product.Price, product.Stock, product.ID)
	return err
}

func (r *productRepo) DeleteByID(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}

func (r *productRepo) GetAll() ([]domain.Product, error) {
	rows, err := r.db.Query("SELECT id, name, category, price, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
