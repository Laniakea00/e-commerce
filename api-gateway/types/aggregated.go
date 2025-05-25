// types/aggregated.go
package types

type OrderItem struct {
	ProductID   int32  `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int32  `json:"quantity"`
}

type AggregatedOrder struct {
	OrderID int32       `json:"order_id"`
	UserID  int32       `json:"user_id"`
	Status  string      `json:"status"`
	Items   []OrderItem `json:"items"`
}
