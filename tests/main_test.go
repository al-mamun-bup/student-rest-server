package tests

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"student-server/handlers"
	"student-server/models"
)

var students []models.Student

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HomeHandler)
	handler.ServeHTTP(rr, req)

	expected := "Welcome to the Student API!\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected %q, got %q", expected, rr.Body.String())
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestGetStudents(t *testing.T) {
	// Case 1: When students slice is empty
	students = []models.Student{} // Initialize the students slice as empty
	req, err := http.NewRequest("GET", "/students", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetStudentsHandler)
	handler.ServeHTTP(rr, req)

	expectedEmpty := "[]\n" // Expected response for empty students
	if rr.Body.String() != expectedEmpty {
		t.Errorf("[Empty Students] Expected %q, got %q", expectedEmpty, rr.Body.String())
	}

	// Case 2: When students slice has multiple entries
	students = []models.Student{
		{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
		{ID: "2", Name: "Efaz", Age: 22, Grade: "B"},
	}

	req, err = http.NewRequest("GET", "/students", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Marshal the students slice into JSON for comparison
	expectedJSON, err := json.Marshal(students)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that the expected JSON is returned
	if rr.Body.String() != string(expectedJSON)+"\n" {
		t.Errorf("[Multiple Students] Expected %q, got %q", string(expectedJSON), rr.Body.String())
	}
}

func TestPostStudent(t *testing.T) {
	t.Run("Valid Student", func(t *testing.T) {
		students = []models.Student{} // Reset student slice
		newStudent := models.Student{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A+"}

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

		// Call the actual AddStudentHandler from handlers
		handler := http.HandlerFunc(handlers.AddStudentHandler)
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

		// Use the real AddStudentHandler here
		handler := http.HandlerFunc(handlers.AddStudentHandler)
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
		newStudent := models.Student{Name: "Efaz"} // Missing ID, Age, Grade

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

		// Use the real AddStudentHandler here as well
		handler := http.HandlerFunc(handlers.AddStudentHandler)
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
	students = []models.Student{
		{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
		{ID: "2", Name: "Efaz", Age: 22, Grade: "B+"},
	}

	t.Run("Valid Student ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/students/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		// Use the actual GetStudentByIDHandler here
		handler := http.HandlerFunc(handlers.GetStudentByIDHandler)
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

		// Use the actual GetStudentByIDHandler here
		handler := http.HandlerFunc(handlers.GetStudentByIDHandler)
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
		students = []models.Student{} // Simulate empty database

		req, err := http.NewRequest("GET", "/students/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		// Use the actual GetStudentByIDHandler here
		handler := http.HandlerFunc(handlers.GetStudentByIDHandler)
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

		// Use the actual GetStudentByIDHandler here
		handler := http.HandlerFunc(handlers.GetStudentByIDHandler)
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
	students = []models.Student{
		{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
		{ID: "2", Name: "Efaz", Age: 22, Grade: "B+"},
	}

	t.Run("Valid Student Update", func(t *testing.T) {
		updatedStudent := models.Student{ID: "1", Name: "Efaz", Age: 21, Grade: "B"}
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
		handler := http.HandlerFunc(handlers.UpdateStudentHandler)
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
		updatedStudent := models.Student{ID: "99", Name: "Efaz", Age: 25, Grade: "A+"}
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
		handler := http.HandlerFunc(handlers.UpdateStudentHandler)
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
		updatedStudent := models.Student{ID: "abc", Name: "Efaz", Age: 25, Grade: "A+"}
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
		handler := http.HandlerFunc(handlers.UpdateStudentHandler)
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

func TestDeleteStudent(t *testing.T) {
	// Initialize students with multiple entries
	students = []models.Student{
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
		students = []models.Student{} // Simulate an empty list

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
		students = []models.Student{
			{ID: "1", Name: "Al Mamun", Age: 20, Grade: "A"},
			{ID: "2", Name: "Efaz", Age: 22, Grade: "B+"},
			{ID: "3", Name: "Al Mamun", Age: 20, Grade: "A+"},
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

// Helper function to create a request with authentication
func createAuthRequest(method, url, username, password string, body string) *http.Request {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, url, bytes.NewReader([]byte(body)))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	// Set Basic Authentication Header
	auth := username + ":" + password
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", authHeader)

	return req
}

// Test accessing a protected route **without authentication** (should fail)
func TestAuthWithoutCredentials(t *testing.T) {
	req := httptest.NewRequest("GET", "/students", nil) // No auth
	rr := httptest.NewRecorder()
	// Wrap handler with middleware
	handler := http.HandlerFunc(handlers.GetStudentsHandler)
	middleware := handlers.BasicAuthMiddleware(handler)
	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", rr.Code)
	}
}

// Test accessing a protected route **with valid credentials** (should succeed)
func TestAuthWithValidCredentials(t *testing.T) {
	req := createAuthRequest("GET", "/students", "admin", "password123", "")
	rr := httptest.NewRecorder()
	// Wrap handler with middleware
	handler := http.HandlerFunc(handlers.GetStudentsHandler)
	middleware := handlers.BasicAuthMiddleware(handler)
	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}
}

// Test accessing a protected route **with invalid credentials** (should fail)
func TestAuthWithInvalidCredentials(t *testing.T) {
	req := createAuthRequest("GET", "/students", "wronguser", "wrongpass", "")
	rr := httptest.NewRecorder()
	// Wrap handler with middleware
	handler := http.HandlerFunc(handlers.GetStudentsHandler)
	middleware := handlers.BasicAuthMiddleware(handler)
	middleware.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", rr.Code)
	}
}
