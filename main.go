package main

import (
	"fmt"
	"log"
	"net/http"

	"student-server/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.AddStudentHandler(w, r)
		} else if r.Method == http.MethodGet {
			handlers.GetStudentsHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/students/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetStudentByIDHandler(w, r)
		case http.MethodPut:
			handlers.UpdateStudentHandler(w, r)
		case http.MethodDelete:
			handlers.DeleteStudentHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
