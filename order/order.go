package order

import (
	"ecommerce/utilities"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateOrder(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"POST"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := r.FormValue("user_id")
	cartId := r.FormValue("cart_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	cartIdInt, err := strconv.Atoi(cartId)
	if err != nil {
		http.Error(w, "Invalid cart ID format", http.StatusBadRequest)
		return
	}

	error := RunOrderTableCreationQuery(conn)
	if error != nil {
		print(error)
		http.Error(w, "Internal Server Error When Creating Order Table"+error.Error(), http.StatusInternalServerError)
		return
	}

	_, err = CreateOrderInDB(conn, userIdInt, cartIdInt)
	if err != nil {
		print(err)
		http.Error(w, "Internal Server Error When Creating Order"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Order created successfully"))

}

func GetOrdersByUserID(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"GET"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userId := r.URL.Query().Get("user_id")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	error := RunOrderTableCreationQuery(conn)
	if error != nil {
		print(error)
		http.Error(w, "Internal Server Error When Creating Order Table"+error.Error(), http.StatusInternalServerError)
		return
	}

	data, err := GetOrdersByUserIDFromDB(conn, userIdInt)
	if err != nil {
		print(err)
		http.Error(w, "Internal Server Error When Fetching Orders"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
