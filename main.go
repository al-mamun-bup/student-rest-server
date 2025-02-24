package main

import (
	"log"
	"net/http"

	"student-server/auth"
	"student-server/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Public route (No authentication required)
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Protected routes (Require authentication)
	protectedRoutes := router.PathPrefix("/students").Subrouter()
	protectedRoutes.Use(auth.BasicAuthMiddleware) // Apply auth middleware
	protectedRoutes.HandleFunc("", handlers.GetStudentsHandler).Methods("GET")
	protectedRoutes.HandleFunc("", handlers.AddStudentHandler).Methods("POST")
	protectedRoutes.HandleFunc("/{id}", handlers.GetStudentByIDHandler).Methods("GET")
	protectedRoutes.HandleFunc("/{id}", handlers.UpdateStudentHandler).Methods("PUT")
	protectedRoutes.HandleFunc("/{id}", handlers.DeleteStudentHandler).Methods("DELETE")

	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
