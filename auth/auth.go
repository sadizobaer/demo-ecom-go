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
	registerData.Username = r.FormValue("username")
	registerData.Email = r.FormValue("email")
	registerData.Password = r.FormValue("password")

	if registerData.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	} else if !utilities.IsValidEmail(registerData.Email) {
		http.Error(w, "Valid Email is required", http.StatusBadRequest)
		return
	} else if registerData.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	error := RunUserTableCreationQuery(conn)
	if error != nil {
		http.Error(w, "Internal Server Error When Creating User Table", http.StatusInternalServerError)
		return
	}

	err := CreateUserInDB(conn, registerData)
	if err != nil {
		http.Error(w, "Internal Server Error When Creating User", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "User registered successfully"}
	json.NewEncoder(w).Encode(response)

}

func LoginUser(w http.ResponseWriter, r *http.Request, conn *pgxpool.Pool) {
	defer r.Body.Close()

	isValidMethod := utilities.IsMethodValid(r.Method, []string{"POST"})
	if !isValidMethod {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginData User
	loginData.Email = r.FormValue("email")
	loginData.Password = r.FormValue("password")

	if !utilities.IsValidEmail(loginData.Email) {
		http.Error(w, "Valid Email is required", http.StatusBadRequest)
		return
	} else if loginData.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	userFromDB, err := GetUserByEmail(conn, loginData.Email)

	print(userFromDB.UserId)
	print(userFromDB.Email)

	if err != nil {
		http.Error(w, "Invalid email", http.StatusUnauthorized)
		return
	}

	isValidPassword := utilities.CheckPasswordHash(loginData.Password, userFromDB.Password)
	if !isValidPassword {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	err = RunTokenTableCreationQuery(conn)
	if err != nil {
		http.Error(w, "Internal Server Error When Creating Token Table", http.StatusInternalServerError)
		return
	}

	authToken, refreshToken, err := utilities.GenerateTokens(userFromDB.Username, userFromDB.Email)

	if err != nil {
		http.Error(w, "Error generating tokens", http.StatusInternalServerError)
		return
	}

	err = StoreTokenInDB(conn, userFromDB.UserId, authToken, refreshToken)

	if err != nil {
		http.Error(w, "Error storing tokens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Login successful", "token": authToken, "refresh": refreshToken}
	json.NewEncoder(w).Encode(response)
}
