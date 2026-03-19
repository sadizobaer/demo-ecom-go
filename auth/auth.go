package auth

import (
	"ecommerce/utilities"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUser(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"POST"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var registerData User
	err := json.NewDecoder(r.Body).Decode(&registerData)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

}

func LoginUser(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"POST"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginData User
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}
