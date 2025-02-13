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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Student API!")
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var student Student
			err := json.NewDecoder(r.Body).Decode(&student)
			if err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			students = append(students, student)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "Student added successfully")
		} else if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(students)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/students/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/students/"):]

		if r.Method == http.MethodGet {
			for _, student := range students {
				if student.ID == id {
					json.NewEncoder(w).Encode(student)
					return
				}
			}
			http.Error(w, "Student not found", http.StatusNotFound)
		} else if r.Method == http.MethodPut {
			var updatedStudent Student
			err := json.NewDecoder(r.Body).Decode(&updatedStudent)
			if err != nil {
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
		} else if r.Method == http.MethodDelete {
			for i, student := range students {
				if student.ID == id {
					students = append(students[:i], students[i+1:]...)
					w.WriteHeader(http.StatusOK)
					fmt.Fprintln(w, "Student deleted successfully")
					return
				}
			}
			http.Error(w, "Student not found", http.StatusNotFound)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
