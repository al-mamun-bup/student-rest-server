package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"student-server/models"

	"github.com/gorilla/mux"
)

var students []models.Student

// HomeHandler handles the root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Student API!")
}

// GetStudentsHandler returns all students
func GetStudentsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
