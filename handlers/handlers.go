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
	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB sets the database instance to be used in the handlers
func SetDB(database *gorm.DB) {
	db = database
	log.Println("Database instance set in handlers")
}

// BasicAuthMiddleware checks for valid basic authentication.
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 || credentials[0] != "admin" || credentials[1] != "password123" {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// HomeHandler handles the root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Student API! Made my Mamun ;)")
}

// GetStudentsHandler retrieves all students from the database
func GetStudentsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var students []models.Student
	if err := db.Find(&students).Error; err != nil {
		http.Error(w, "Failed to fetch students", http.StatusInternalServerError)
		return
	}

	// Log and simulate slow request (optional)
	log.Println("Processing request: Fetching students...")
	time.Sleep(2 * time.Second)
	log.Println("Request completed: Students fetched")

	json.NewEncoder(w).Encode(students)
}

// AddStudentHandler adds a new student to the database
func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := db.Create(&student).Error; err != nil {
		http.Error(w, "Failed to add student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Student added successfully")
}

// GetStudentByIDHandler retrieves a student by ID
func GetStudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var student models.Student
	if err := db.First(&student, id).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// UpdateStudentHandler updates an existing student's details
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var student models.Student
	if err := db.First(&student, id).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := db.Save(&student).Error; err != nil {
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Student updated successfully")
}

// DeleteStudentHandler deletes a student by ID
func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var student models.Student
	if err := db.First(&student, id).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	if err := db.Delete(&student).Error; err != nil {
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Student deleted successfully")
}
