package wishlist

import (
	"ecommerce/utilities"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetWishlistProducts(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {

	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"GET"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := RunWishlistTableCreationQuery(conn)
	if err != nil {
		print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := r.FormValue("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	data, err := GetWishlistByUserIDFromDB(conn, userIdInt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if data == nil {
		data = []WishlistItemView{}
	}

	dataMap := map[string]interface{}{"wishlist": data}
	json.NewEncoder(w).Encode(dataMap)

}

func AddProductToWishlist(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"POST"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := r.FormValue("user_id")
	productId := r.FormValue("product_id")

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		http.Error(w, "Invalid product ID format", http.StatusBadRequest)
		return
	}

	error := RunWishlistTableCreationQuery(conn)
	if error != nil {
		print(error)
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	wishlistID, err := AddToWishlistInDB(conn, userIdInt, productIdInt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"message": "Product added to wishlist successfully", "wishlist_id": wishlistID}
	json.NewEncoder(w).Encode(response)

}

func RemoveProductFromWishlist(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"DELETE"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := r.FormValue("user_id")
	productId := r.FormValue("product_id")

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		http.Error(w, "Invalid product ID format", http.StatusBadRequest)
		return
	}

	error := RunWishlistTableCreationQuery(conn)
	if error != nil {
		print(error)
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}

	err = RemoveFromWishlistInDB(conn, userIdInt, productIdInt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Product removed from wishlist successfully"}
	json.NewEncoder(w).Encode(response)

}
