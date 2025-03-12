
# 📚 **Student REST API Documentation**

This API allows you to manage student data, including CRUD operations: **Create, Read, Update, Delete**.

---

## **🔹 Base URL**

```
http://localhost:9090
```

---

## **📌 Endpoints**

### **1️⃣ Get All Students**
- **Endpoint:** `GET /students`
- **Description:** Fetches all student records.
- **Request:**
```bash
curl -X GET http://localhost:9090/students
```
- **Response:**
```json
[
    {
        "id": 1,
        "name": "Efaz",
        "age": 21,
        "grade": "A"
    },
    {
        "id": 2,
        "name": "Sadat",
        "age": 21,
        "grade": "A"
    },
    {
        "id": 3,
        "name": "Mamun",
        "age": 24,
        "grade": "A"
    }
]
```
- **Success Response:**
  - **Code:** 200
  - **Content:** List of all students.

- **Error Response:**
  - **Code:** 500
  - **Content:** Internal Server Error if something goes wrong.

---

### **2️⃣ Get Student by ID**
- **Endpoint:** `GET /students/{id}`
- **Description:** Fetches a single student by their ID.
- **Request:**
```bash
curl -X GET http://localhost:9090/students/1
```
- **Response:**
```json
{
    "id": 1,
    "name": "Efaz",
    "age": 21,
    "grade": "A"
}
```
- **Success Response:**
  - **Code:** 200
  - **Content:** The student's details.

- **Error Response:**
  - **Code:** 404
  - **Content:** If student with the given ID does not exist.

---

### **3️⃣ Create a New Student**
- **Endpoint:** `POST /students`
- **Description:** Creates a new student record.
- **Request:**
```bash
curl -X POST http://localhost:9090/students      -H "Content-Type: application/json"      -d '{
           "name": "Alice Smith",
           "age": 23,
           "grade": "A"
         }'
```
- **Response:**
```json
{
    "id": 4,
    "name": "Alice Smith",
    "age": 23,
    "grade": "A"
}
```
- **Success Response:**
  - **Code:** 201
  - **Content:** The newly created student’s details.

- **Error Response:**
  - **Code:** 400
  - **Content:** If required fields are missing or invalid data.

---

### **4️⃣ Update a Student**
- **Endpoint:** `PUT /students/{id}`
- **Description:** Updates an existing student's details by ID.
- **Request:**
```bash
curl -X PUT http://localhost:9090/students/1      -H "Content-Type: application/json"      -d '{
           "name": "Alice Brown",
           "age": 24,
           "grade": "A+"
         }'
```
- **Response:**
```json
{
    "id": 1,
    "name": "Alice Brown",
    "age": 24,
    "grade": "A+"
}
```
- **Success Response:**
  - **Code:** 200
  - **Content:** The updated student's details.

- **Error Response:**
  - **Code:** 404
  - **Content:** If student with the given ID does not exist.

---

### **5️⃣ Delete a Student**
- **Endpoint:** `DELETE /students/{id}`
- **Description:** Deletes a student by their ID.
- **Request:**
```bash
curl -X DELETE http://localhost:9090/students/2
```
- **Response:**
```json
{
    "message": "Student deleted successfully"
}
```
- **Success Response:**
  - **Code:** 200
  - **Content:** A confirmation message that the student was deleted.

- **Error Response:**
  - **Code:** 404
  - **Content:** If student with the given ID does not exist.

---

## **🔹 Notes**
- All requests should be in **JSON** format.
- The `ID` is required for fetching, updating, and deleting students.

---