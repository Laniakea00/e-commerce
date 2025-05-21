package domain

type Order struct {
	ID     int         `json:"id"`
	UserID int         `json:"user_id"`
	Status string      `json:"status"`
	Items  []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
