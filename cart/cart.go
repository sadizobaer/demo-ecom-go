package cart

import (
	"ecommerce/utilities"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetCartItems(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"GET"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := RunCartTableCreationQuery(conn)
	if err != nil {
		print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := r.URL.Query().Get("user_id") // Using Query for GET requests is cleaner
	userIdInt, _ := strconv.Atoi(userId)

	cartData, err := GetCartItemsByUserIDFromDB(conn, userIdInt)
	if err != nil {
		// If cart is empty, return empty structure rather than 500
		cartData = CartView{UserID: userIdInt, Items: []CartItem{}}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cartData)
}

func AddProductToCart(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"POST"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := RunCartTableCreationQuery(conn)
	if err != nil {
		print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userIdInt, _ := strconv.Atoi(r.FormValue("user_id"))
	productIdInt, _ := strconv.Atoi(r.FormValue("product_id"))
	quantityInt, _ := strconv.Atoi(r.FormValue("quantity"))

	cartID, err := AddProductToCartInDB(conn, CartItemAdd{
		UserID:    userIdInt,
		ProductID: productIdInt,
		Quantity:  quantityInt,
	})

	if err != nil {
		http.Error(w, "Failed to add to cart: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Product added to cart",
		"cart_id": cartID,
	})
}

func UpdateCartItemQuantity(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"PUT"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Parse Inputs
	cartId := r.FormValue("cart_id")
	productId := r.FormValue("product_id")
	quantity := r.FormValue("quantity")

	cartIdInt, err := strconv.Atoi(cartId)
	productIdInt, err2 := strconv.Atoi(productId)
	quantityInt, err3 := strconv.Atoi(quantity)

	if err != nil || err2 != nil || err3 != nil {
		http.Error(w, "Invalid input format. Ensure cart_id, product_id, and quantity are numbers.", http.StatusBadRequest)
		return
	}

	// 3. Update Database
	err = UpdateCartItemQuantityInDB(conn, cartIdInt, productIdInt, quantityInt)
	if err != nil {
		// Handle case where item wasn't in cart
		if err.Error() == "item not found in cart" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 4. Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Cart updated successfully",
	})
}

func RemoveProductFromCart(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"DELETE"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cartIdInt, _ := strconv.Atoi(r.FormValue("cart_id"))
	productIdInt, _ := strconv.Atoi(r.FormValue("product_id"))

	err := RemoveProductFromCartInDB(conn, cartIdInt, productIdInt)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item removed"})
}
