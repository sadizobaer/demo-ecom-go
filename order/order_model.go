package order

import "ecommerce/cart"

type OrderItemAdd struct {
	UserID int `json:"user_id"`
	CartID int `json:"cart_id"`
}

type OrderItemView struct {
	ID          int             `json:"order_id"`
	UserID      int             `json:"user_id"`
	Products    []cart.CartItem `json:"items"`
	TotalAmount float64         `json:"total_amount"`
	Status      string          `json:"status"`
}
