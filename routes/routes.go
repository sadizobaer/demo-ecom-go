package routes

import (
	"ecommerce/auth"
	"ecommerce/category"
	"ecommerce/product"
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
}
