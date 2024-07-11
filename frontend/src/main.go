package main

import (
	"literary-lions/frontend/src/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	log.Println("Frontend server started at :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
