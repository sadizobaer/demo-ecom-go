package routes

import (
	"ecommerce/auth"
	"ecommerce/cart"
	"ecommerce/category"
	"ecommerce/favorite"
	"ecommerce/order"
	"ecommerce/product"
	"ecommerce/wishlist"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AllPaths(handler *http.ServeMux, conn *pgxpool.Pool) {

	// ======== AUTH ROUTES ========
	handler.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		auth.RegisterUser(w, r, conn)
	})
	handler.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginUser(w, r, conn)
	})

	// ========== CATEGORY ROUTES ==========
	handler.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		category.GetAllCategories(w, r, conn)
	})
	handler.HandleFunc("/categories/create", func(w http.ResponseWriter, r *http.Request) {
		category.CreateCategory(w, r, conn)
	})
	handler.HandleFunc("/categories/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		category.UpdateCategory(w, r, conn)
	})
	handler.HandleFunc("/categories/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		category.DeleteCategory(w, r, conn)
	})

	// ======== PRODUCT ROUTES ========
	handler.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		product.GetAllProducts(w, r, conn)
	})
	handler.HandleFunc("/products/{id}", func(w http.ResponseWriter, r *http.Request) {
		product.GetProductByID(w, r, conn)
	})
	handler.HandleFunc("/products/create", func(w http.ResponseWriter, r *http.Request) {
		product.CreateProduct(w, r, conn)
	})
	handler.HandleFunc("/products/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		product.UpdateProduct(w, r, conn)
	})
	handler.HandleFunc("/products/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		product.DeleteProduct(w, r, conn)
	})

	// ======== WISHLIST ROUTES ========
	handler.HandleFunc("/wishlist", func(w http.ResponseWriter, r *http.Request) {
		wishlist.GetWishlistProducts(w, r, conn)
	})
	handler.HandleFunc("/wishlist/add", func(w http.ResponseWriter, r *http.Request) {
		wishlist.AddProductToWishlist(w, r, conn)
	})
	handler.HandleFunc("/wishlist/remove", func(w http.ResponseWriter, r *http.Request) {
		wishlist.RemoveProductFromWishlist(w, r, conn)
	})

	// ======== FAVORITE ROUTES ========
	handler.HandleFunc("/favorites", func(w http.ResponseWriter, r *http.Request) {
		favorite.GetFavoriteProducts(w, r, conn)
	})
	handler.HandleFunc("/favorites/add", func(w http.ResponseWriter, r *http.Request) {
		favorite.AddProductToFavorite(w, r, conn)
	})
	handler.HandleFunc("/favorites/remove", func(w http.ResponseWriter, r *http.Request) {
		favorite.RemoveProductFromFavorite(w, r, conn)
	})

	// ======== CART ROUTES ========
	handler.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		cart.GetCartItems(w, r, conn)
	})
	handler.HandleFunc("/cart/add", func(w http.ResponseWriter, r *http.Request) {
		cart.AddProductToCart(w, r, conn)
	})
	handler.HandleFunc("/cart/update", func(w http.ResponseWriter, r *http.Request) {
		print("Update Cart Route Hit")
		cart.UpdateCartItemQuantity(w, r, conn)
	})
	handler.HandleFunc("/cart/remove", func(w http.ResponseWriter, r *http.Request) {
		cart.RemoveProductFromCart(w, r, conn)
	})

	// ======== ORDER ROUTES ========

	handler.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		order.GetOrdersByUserID(w, r, conn)
	})
	handler.HandleFunc("/orders/create", func(w http.ResponseWriter, r *http.Request) {
		order.CreateOrder(w, r, conn)
	})

}
