package wishlist

import "ecommerce/product"

type WishlistItemAdd struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

type WishlistItemView struct {
	ID      int                 `json:"wishlist_id"`
	Product product.ProductView `json:"product"`
}
