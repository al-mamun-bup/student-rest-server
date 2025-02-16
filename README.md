# Student REST API

![Go](https://img.shields.io/badge/Go-1.20-blue)
![License](https://img.shields.io/badge/license-MIT-green)

A simple REST API built in Go to manage student records.

## Features
- Add a new student
- Retrieve all students
- Retrieve a student by ID
- Update student details
- Delete a student
- Unit testing included

## Repository
[GitHub Repository](https://github.com/al-mamun-bup/student-rest-server)

## Installation & Running
### Prerequisites
- Go installed (v1.20 or later)

### Steps
```sh
git clone https://github.com/al-mamun-bup/student-rest-server.git
cd student-rest-server
go run main.go
```

Server runs at: `http://localhost:8080`

## API Endpoints
| Method | Endpoint        | Description           |
|--------|---------------|----------------------|
| GET    | `/`           | Welcome message     |
| GET    | `/students`   | Get all students    |
| POST   | `/students`   | Add a new student   |
| GET    | `/students/{id}` | Get student by ID |
| PUT    | `/students/{id}` | Update student   |
| DELETE | `/students/{id}` | Delete student   |

## Example Requests
### Add a Student
```sh
curl -X POST http://localhost:8080/students \
     -H "Content-Type: application/json" \
     -d '{"id": "1", "name": "John Doe", "age": 20, "grade": "A"}'
```

### Get All Students
```sh
curl -X GET http://localhost:8080/students
```

## Running Tests
```sh
go test ./...
```

## License
This project is licensed under the MIT License.

