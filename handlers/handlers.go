package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"student-server/models"

	"github.com/gorilla/mux"
)

var students []models.Student

var Students []models.Student

// BasicAuthMiddleware will check for valid basic authentication.
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")

		// Check if the Authorization header is missing or invalid
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse the "Basic" part of the Authorization header
		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decode the base64-encoded credentials
		decoded, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Split into username and password
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 || credentials[0] != "admin" || credentials[1] != "password123" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Pass control to the next handler if authentication is successful
		next.ServeHTTP(w, r)
	})
}

// HomeHandler handles the root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Student API!")
}

// GetStudentsHandler returns all students
// GetStudentsHandler returns all students as JSON
// func GetStudentsHandler(w http.ResponseWriter, _ *http.Request) {
// 	// If students slice is empty, return an empty array `[]` and not `null`
// 	w.Header().Set("Content-Type", "application/json")
// 	if len(students) == 0 {
// 		// Explicitly returning an empty array
// 		json.NewEncoder(w).Encode([]models.Student{})
// 		return
// 	}

//		// Otherwise, encode the students list into JSON
//		json.NewEncoder(w).Encode(students)
//	}
func GetStudentsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Simulate a long-running request
	log.Println("Processing request: Fetching students...") // Log request start
	time.Sleep(2 * time.Second)                             // Simulate slow query or processing
	log.Println("Request completed: Students fetched")      // Log request end

	// If students slice is empty, return an empty array `[]` instead of `null`
	if len(students) == 0 {
		json.NewEncoder(w).Encode([]models.Student{})
		return
	}

	// Otherwise, return the list of students
	json.NewEncoder(w).Encode(students)
}

// AddStudentHandler adds a student
func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	students = append(students, student)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Student added successfully")
}

// GetStudentByIDHandler fetches a student by ID
func GetStudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, student := range students {
		if student.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(student)
			return
		}
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}

// UpdateStudentHandler updates a student's details
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedStudent models.Student
	if err := json.NewDecoder(r.Body).Decode(&updatedStudent); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	for i, student := range students {
		if student.ID == id {
			students[i] = updatedStudent
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Student updated successfully")
			return
		}
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}

// DeleteStudentHandler deletes a student by ID
func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for i, student := range students {
		if student.ID == id {
			students = append(students[:i], students[i+1:]...)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Student deleted successfully")
			return
		}
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}
