package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainFunctionality(t *testing.T) {

}

func TestRoothandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the Student API!"))
	})
	handler.ServeHTTP(rr, req)
	expected := "Welcome to the Student API!"
	if rr.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, rr.Body.String())
	}
}

func TestGetStudents(t *testing.T) {
	req, err := http.NewRequest("GET", "/students", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]")) // Expecting an empty list initially
	})

	handler.ServeHTTP(rr, req)

	expected := "[]"
	if rr.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, rr.Body.String())
	}
}

func TestPostStudent(t *testing.T) {
	newStudent := Student{
		ID:    "1",
		Name:  "Mamun",
		Age:   23,
		Grade: "A+",
	}

	studentJSON, err := json.Marshal(newStudent)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/students", bytes.NewBuffer(studentJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var student Student
		json.NewDecoder(r.Body).Decode(&student)
		students = append(students, student)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Student added successfully"))
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	expected := "Student added successfully"
	if rr.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, rr.Body.String())
	}
}
func TestGetStudentByID(t *testing.T) {
	// First, add a student to the in-memory list
	students = []Student{
		{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
	}

	req, err := http.NewRequest("GET", "/students/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/students/"):]
		for _, student := range students {
			if student.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(student)
				return
			}
		}
		http.Error(w, "Student not found", http.StatusNotFound)
	})

	handler.ServeHTTP(rr, req)

	expected := `{"id":"1","name":"Al Mamun","age":20,"grade":"A"}` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, rr.Body.String())
	}
}
func TestUpdateStudent(t *testing.T) {
	// Initialize students with a student
	students = []Student{
		{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
	}

	updatedStudent := Student{ID: "1", Name: "Efaz", Age: 21, Grade: "B"}
	updatedJSON, err := json.Marshal(updatedStudent)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/students/1", bytes.NewBuffer(updatedJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/students/"):]
		var updated Student
		json.NewDecoder(r.Body).Decode(&updated)

		for i, student := range students {
			if student.ID == id {
				students[i] = updated
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Student updated successfully"))
				return
			}
		}
		http.Error(w, "Student not found", http.StatusNotFound)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expected := "Student updated successfully"
	if rr.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, rr.Body.String())
	}

	// Verify that the student was actually updated
	if students[0].Name != "Efaz" || students[0].Age != 21 || students[0].Grade != "B" {
		t.Errorf("Student was not updated correctly")
	}
}

func TestDeleteStudent(t *testing.T) {
	// Initialize students with a student
	students = []Student{
		{ID: "1", Name: "John Doe", Age: 20, Grade: "A"},
	}

	req, err := http.NewRequest("DELETE", "/students/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/students/"):]
		for i, student := range students {
			if student.ID == id {
				students = append(students[:i], students[i+1:]...)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Student deleted successfully"))
				return
			}
		}
		http.Error(w, "Student not found", http.StatusNotFound)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expected := "Student deleted successfully"
	if rr.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, rr.Body.String())
	}

	// Verify that the student was actually deleted
	if len(students) != 0 {
		t.Errorf("Student was not deleted")
	}
}
