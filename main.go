package main

import (
	"log"
	"net/http"

	"student-server/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// Created a new Gorilla Mux router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	router.HandleFunc("/students", handlers.GetStudentsHandler).Methods("GET")
	router.HandleFunc("/students", handlers.AddStudentHandler).Methods("POST")
	router.HandleFunc("/students/{id}", handlers.GetStudentByIDHandler).Methods("GET")
	router.HandleFunc("/students/{id}", handlers.UpdateStudentHandler).Methods("PUT")
	router.HandleFunc("/students/{id}", handlers.DeleteStudentHandler).Methods("DELETE")

	// Starting the server
	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
