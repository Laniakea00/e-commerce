package domain

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CategoryID  int     `json:"category_id"`
	Size        string  `json:"size"`
	Color       string  `json:"color"`
	Gender      string  `json:"gender"`
	Material    string  `json:"material"`
	Season      string  `json:"season"`
}
