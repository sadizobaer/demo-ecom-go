package main

import (
	"ecommerce/db"
	"ecommerce/middlewares"
	"ecommerce/routes"
	"log"
	"net/http"
)

func main() {
	conn, err := db.ConnectDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	defer conn.Close()

	handler := http.NewServeMux()

	routes.AllPaths(handler, conn)

	authMiddleware := middlewares.AuthMiddleware(handler, conn)

	log.Fatal(http.ListenAndServe(":8080", authMiddleware))
}
