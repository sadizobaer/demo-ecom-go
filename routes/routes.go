package routes

import (
	"fmt"
	"net/http"
	"time"
)

func timeHandler(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	fmt.Fprintf(w, "The time is: %s", tm)
}

func AllPaths(handler *http.ServeMux) {
	handler.HandleFunc("/", timeHandler)
}
