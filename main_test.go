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
	// Case 1: When students slice is empty
	students = []Student{} // Reset students list

	req, err := http.NewRequest("GET", "/students", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(students)
	})

	handler.ServeHTTP(rr, req)

	expectedEmpty := "[]\n"
	if rr.Body.String() != expectedEmpty {
		t.Errorf("[Empty Students] Expected %q, got %q", expectedEmpty, rr.Body.String())
	}

	// Case 2: When students slice has multiple entries
	students = []Student{
		{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
		{ID: "2", Name: "Efaz", Age: 22, Grade: "B"},
	}

	req, err = http.NewRequest("GET", "/students", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expectedJSON, _ := json.Marshal(students)
	if rr.Body.String() != string(expectedJSON)+"\n" {
		t.Errorf("[Multiple Students] Expected %q, got %q", string(expectedJSON), rr.Body.String())
	}
}

func TestPostStudent(t *testing.T) {
	t.Run("Valid Student", func(t *testing.T) {
		students = []Student{} // Reset student slice
		newStudent := Student{ID: "1", Name: "Al Mamun", Age: 23, Grade: "A+"}

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

		// Verify that student is added
		if len(students) != 1 || students[0].Name != "Al Mamun" {
			t.Errorf("Student was not added correctly")
		}
	})

	t.Run("Invalid JSON Payload", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/students", bytes.NewBuffer([]byte("invalid-json")))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var student Student
			err := json.NewDecoder(r.Body).Decode(&student)
			if err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			students = append(students, student)

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Student added successfully"))
		})

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		expected := "Invalid input\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("Missing Fields", func(t *testing.T) {
		newStudent := Student{Name: "Efaz"} // Missing ID, Age, Grade

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
			err := json.NewDecoder(r.Body).Decode(&student)
			if err != nil || student.ID == "" || student.Age == 0 || student.Grade == "" {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}
			students = append(students, student)

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Student added successfully"))
		})

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		expected := "Invalid input\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})
}

func TestGetStudentByID(t *testing.T) {
	// Initialize test data
	students = []Student{
		{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
		{ID: "2", Name: "Efaz", Age: 22, Grade: "B+"},
	}

	t.Run("Valid Student ID", func(t *testing.T) {
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

		expected := `{"ID":"1","Name":"Al Mamun","Age":20,"Grade":"A"}` + "\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("Non-Existent Student ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/students/39", nil) // ID 39 doesn't exist
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

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
		}

		expected := "Student not found\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("Empty Student List", func(t *testing.T) {
		students = []Student{} // Simulate empty database

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

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
		}

		expected := "Student not found\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("Invalid ID Format", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/students/abc", nil) // Non-numeric ID
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.URL.Path[len("/students/"):]
			if id == "" || id == "abc" {
				http.Error(w, "Invalid student ID", http.StatusBadRequest)
				return
			}
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

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		expected := "Invalid student ID\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})
}

func TestUpdateStudent(t *testing.T) {
	// Initialize students with test data
	students = []Student{
		{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
		{ID: "2", Name: "Efaz", Age: 22, Grade: "B+"},
	}

	t.Run("Valid Student Update", func(t *testing.T) {
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
	})

	t.Run("Updating Non-Existent Student", func(t *testing.T) {
		updatedStudent := Student{ID: "99", Name: "Efaz", Age: 25, Grade: "A+"}
		updatedJSON, err := json.Marshal(updatedStudent)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("PUT", "/students/99", bytes.NewBuffer(updatedJSON))
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

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
		}

		expected := "Student not found\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("Empty Student List", func(t *testing.T) {
		students = []Student{} // Simulate empty database

		updatedStudent := Student{ID: "1", Name: "Efaz", Age: 23, Grade: "A+"}
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

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
		}

		expected := "Student not found\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("Invalid ID Format", func(t *testing.T) {
		updatedStudent := Student{ID: "abc", Name: "Efaz", Age: 25, Grade: "A+"}
		updatedJSON, err := json.Marshal(updatedStudent)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("PUT", "/students/abc", bytes.NewBuffer(updatedJSON)) // Non-numeric ID
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.URL.Path[len("/students/"):]
			if id == "" || id == "abc" {
				http.Error(w, "Invalid student ID", http.StatusBadRequest)
				return
			}
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

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		expected := "Invalid student ID\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("Malformed JSON Body", func(t *testing.T) {
		malformedJSON := `{"ID": "1", "Name": "Efaz", "Age": 23, "Grade": A+}` // Missing double quotes around "A+"

		req, err := http.NewRequest("PUT", "/students/1", bytes.NewBuffer([]byte(malformedJSON)))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.URL.Path[len("/students/"):]
			var updated Student
			err := json.NewDecoder(r.Body).Decode(&updated)
			if err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

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

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		expected := "Invalid request body\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})
}

func TestDeleteStudent(t *testing.T) {
	// Initialize students with multiple entries
	students = []Student{
		{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
		{ID: "2", Name: "Efaz", Age: 22, Grade: "B+"},
	}

	t.Run("Delete Existing Student", func(t *testing.T) {
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
		for _, student := range students {
			if student.ID == "1" {
				t.Errorf("Student was not deleted")
			}
		}
	})

	t.Run("Delete Non-Existent Student", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/students/99", nil)
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

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
		}

		expected := "Student not found\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("Delete from Empty Student List", func(t *testing.T) {
		students = []Student{} // Simulate an empty list

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

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
		}

		expected := "Student not found\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("Invalid ID Format", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/students/abc", nil) // Non-numeric ID
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.URL.Path[len("/students/"):]
			if id == "" || id == "abc" {
				http.Error(w, "Invalid student ID", http.StatusBadRequest)
				return
			}
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

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}

		expected := "Invalid student ID\n"
		if rr.Body.String() != expected {
			t.Errorf("Expected %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("Ensure Remaining Students Exist", func(t *testing.T) {
		students = []Student{
			{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
			{ID: "2", Name: "Efaz", Age: 22, Grade: "B+"},
			{ID: "3", Name: "Al Mamun", Age: 25, Grade: "A+"},
		}

		req, err := http.NewRequest("DELETE", "/students/2", nil) // Deleting "Efaz"
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

		// Verify that only student with ID "2" was deleted, and others remain
		if len(students) != 2 {
			t.Errorf("Expected 2 students remaining, got %d", len(students))
		}

		for _, student := range students {
			if student.ID == "2" {
				t.Errorf("Student with ID 2 was not deleted")
			}
		}
	})
}
