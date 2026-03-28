package cart

import "ecommerce/product"

type CartView struct {
	ID     int        `json:"cart_id"`
	UserID int        `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartItemAdd struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CartItem struct {
	Product  product.ProductView `json:"product"`
	Quantity int                 `json:"quantity"`
}
