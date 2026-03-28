package favorite

import "ecommerce/product"

type FavoriteItemAdd struct {
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

type FavoriteItemView struct {
	ID      int                 `json:"favorite_id"`
	Product product.ProductView `json:"product"`
}
