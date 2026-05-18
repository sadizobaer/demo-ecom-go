package middlewares

import (
	"context"
	"ecommerce/utilities"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const tokenExistenceQuery = `SELECT COUNT(*) FROM tokens WHERE token = $1;`

// ── Public routes that require NO authentication ──────────────
// These are readable by anyone (storefront product/category listing).
var publicPaths = []string{
	"/register",
	"/login",
	"/products",
	"/categories",
}

// isPublicPath returns true when the request path needs no token.
func isPublicPath(path string) bool {
	for _, p := range publicPaths {
		if p == path || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

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

// setCORSHeaders writes the CORS headers that allow the Next.js
// frontend (localhost:3000) to call the Go backend (localhost:8080).
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// AuthMiddleware handles CORS preflight, injects CORS headers on every
// response, and enforces JWT authentication on protected routes.
func AuthMiddleware(next http.Handler, conn *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Always set CORS headers first
		setCORSHeaders(w)

		// 2. Handle OPTIONS preflight immediately — no further processing needed
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// 3. Allow public paths without a token
		if isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// 4. Protected routes — extract Bearer token
		authHeader := r.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, "Bearer ", 2)
		if len(parts) < 2 || parts[1] == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		isValid, err := utilities.ValidateToken(token)
		if err != nil || !isValid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		isInDB, err := isTokenInDB(token, conn)
		if err != nil || !isInDB {
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
