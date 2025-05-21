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
	DecreaseStock(productId int, quantity int) error
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(product domain.Product) error {
	_, err := r.db.Exec(`
		INSERT INTO products 
		(name, description, price, stock, category_id, size, color, gender, material, season) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		product.Name, product.Description, product.Price, product.Stock, product.CategoryID,
		product.Size, product.Color, product.Gender, product.Material, product.Season)
	return err
}

func (r *productRepo) GetByID(id int) (domain.Product, error) {
	var p domain.Product
	row := r.db.QueryRow(`
		SELECT id, name, description, price, stock, category_id, size, color, gender, material, season 
		FROM products WHERE id = ?`, id)
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID,
		&p.Size, &p.Color, &p.Gender, &p.Material, &p.Season)
	return p, err
}

func (r *productRepo) Update(product domain.Product) error {
	_, err := r.db.Exec(`
		UPDATE products SET 
			name = ?, 
			description = ?, 
			price = ?, 
			stock = ?, 
			category_id = ?, 
			size = ?, 
			color = ?, 
			gender = ?, 
			material = ?, 
			season = ?
		WHERE id = ?`,
		product.Name, product.Description, product.Price, product.Stock, product.CategoryID,
		product.Size, product.Color, product.Gender, product.Material, product.Season, product.ID)
	return err
}

func (r *productRepo) DeleteByID(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	return err
}

func (r *productRepo) GetAll() ([]domain.Product, error) {
	rows, err := r.db.Query(`
		SELECT id, name, description, price, stock, category_id, size, color, gender, material, season 
		FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CategoryID,
			&p.Size, &p.Color, &p.Gender, &p.Material, &p.Season)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepo) DecreaseStock(productId int, quantity int) error {
	_, err := r.db.Exec("UPDATE products SET stock = stock - ? WHERE id = ?", quantity, productId)
	return err
}
