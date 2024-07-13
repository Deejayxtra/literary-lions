package main

import (
	"log"
	"net/http"

	"literary-lions/backend/src/internal/db"
	"literary-lions/backend/src/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	handlers.InitHandlers(database) // Pass the database instance to handlers

	err = handlers.InitAdminUser(database) // Pass the database instance to InitAdminUser
	if err != nil {
		log.Fatalf("Could not initialize admin user: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected endpoints
	r.HandleFunc("/posts", handlers.IsAuthorized(handlers.CreatePost, "user")).Methods("POST")
	r.HandleFunc("/posts/{id}", handlers.IsAuthorized(handlers.UpdatePost, "user")).Methods("PUT")
	r.HandleFunc("/posts/{id}", handlers.IsAuthorized(handlers.DeletePost, "user")).Methods("DELETE")
	r.HandleFunc("/admin/users", handlers.IsAuthorized(handlers.GetAllUsers, "admin")).Methods("GET")
	r.HandleFunc("/admin/users/{id}", handlers.IsAuthorized(handlers.DeleteUser, "admin")).Methods("DELETE")
	r.HandleFunc("/admin/users/{id}/role", handlers.IsAuthorized(handlers.UpdateUserRole, "admin")).Methods("PUT")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
