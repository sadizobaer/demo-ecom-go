package main

import (
	"ecommerce/routes"
	"log"
	"net/http"
)

func main() {

	handler := http.NewServeMux()

	routes.AllPaths(handler)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
