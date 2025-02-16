package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Student struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade string `json:"grade"`
}

var students []Student

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Student API!")
}

func getStudentsHandler(w http.ResponseWriter, _ *http.Request) {
	json.NewEncoder(w).Encode(students)
}

func addStudentHandler(w http.ResponseWriter, r *http.Request) {
	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	students = append(students, student)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Student added successfully")
}

func getStudentByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/students/"):]
	for _, student := range students {
		if student.ID == id {
			json.NewEncoder(w).Encode(student)
			return
		}
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}

func updateStudentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/students/"):]
	var updatedStudent Student
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

func deleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/students/"):]
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

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			addStudentHandler(w, r)
		} else if r.Method == http.MethodGet {
			getStudentsHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/students/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getStudentByIDHandler(w, r)
		case http.MethodPut:
			updateStudentHandler(w, r)
		case http.MethodDelete:
			deleteStudentHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
