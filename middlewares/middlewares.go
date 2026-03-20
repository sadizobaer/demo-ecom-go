package middlewares

import (
	"context"
	"ecommerce/utilities"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const tokenExistenceQuery = `SELECT COUNT(*) FROM tokens WHERE token = $1;`

func executeCountQuery(query string, token string, conn *pgxpool.Pool) (int, error) {

	var count int
	err := conn.QueryRow(context.Background(), query, token).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil

}

func isTokenInDB(token string, conn *pgxpool.Pool) (bool, error) {

	count, err := executeCountQuery(tokenExistenceQuery, token, conn)
	if err != nil {
		return false, err
	}

	return count > 0, nil

}

func isPathInList(path string, paths []string) bool {
	for _, p := range paths {
		if p == path {
			return true
		}
	}
	return false
}

func AuthMiddleware(next http.Handler, conn *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		nonAuthorizationPaths := []string{"/register", "/login"}
		if isPathInList(r.URL.Path, nonAuthorizationPaths) {
			next.ServeHTTP(w, r)
			return
		}

		token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
		if token == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		isValid, err := utilities.ValidateToken(token)
		if err != nil || !isValid {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		isTokenInDB, err := isTokenInDB(token, conn)
		if err != nil || !isTokenInDB {
			http.Error(w, "Token not found in database", http.StatusUnauthorized)
			return
		}

		isAdminEndpoint := utilities.IsAdminEndpoint(r.URL.Path)
		if isAdminEndpoint {
			isAdmin, err := utilities.IsAdminToken(token)
			if err != nil || !isAdmin {
				http.Error(w, "Admin privileges required", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
