# 🎓 Student REST API

![Go](https://img.shields.io/badge/Go-1.23-blue) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue) ![Docker](https://img.shields.io/badge/Docker-Supported-blue) ![License](https://img.shields.io/badge/license-MIT-green)

A simple, data-persistent REST API built in Go to manage student records, powered by PostgreSQL and Docker. 🚀

## ✨ Features
- 📝 Add, retrieve, update, and delete students
- 🗄️ Persistent storage with PostgreSQL
- 🐳 Dockerized setup with `docker-compose`
- 🔒 Basic authentication support
- ⚙️ CLI support for specifying port number
- ✅ Graceful shutdown
- 🧪 Unit testing included

## 📂 Repository
[🔗 GitHub Repository](https://github.com/al-mamun-bup/student-rest-server)

## ⚡ Installation & Running
### 🔧 Prerequisites
- Go (v1.23 or later)
- Docker & Docker Compose
- PostgreSQL (optional if using Docker)

### 🚀 Running the Server
#### Using Makefile
To start the server with default settings (port `8080`):
```sh
make serve
```
To run on a custom port, e.g., `9090`:
```sh
make serve PORT=9090
```

#### 🐳 Using Docker Compose
Run the server and database in a Docker environment:
```sh
docker-compose up --build
```
This will start the API and PostgreSQL database inside the same Docker network.

#### 💻 Manually Running
If running locally with PostgreSQL installed, set up a `.env` file with database credentials:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=mamun
DB_PASSWORD=1234
DB_NAME=student_db
```
Then run:
```sh
go run main.go
```
📡 Server runs at: `http://localhost:8080`

## 🔐 Authentication
This API supports basic authentication. To access protected endpoints, include the `Authorization` header:
```sh
curl -u username:password http://localhost:8080/students
```

## 🔗 API Endpoints
| 🛠️ Method | 🌍 Endpoint        | 📌 Description           |
|--------|---------------|----------------------|
| GET    | `/`           | Welcome message     |
| GET    | `/students`   | Get all students    |
| POST   | `/students`   | Add a new student   |
| GET    | `/students/{id}` | Get student by ID |
| PUT    | `/students/{id}` | Update student   |
| DELETE | `/students/{id}` | Delete student   |

## 📤 Example Requests
### ➕ Add a Student
```sh
curl -X POST http://localhost:8080/students \
     -H "Content-Type: application/json" \
     -d '{"name": "John Doe", "age": 20, "grade": "A"}'
```

### 📋 Get All Students
```sh
curl -X GET http://localhost:8080/students
```

## 🧪 Running Tests
```sh
go test ./...
```

## 📜 License
This project is licensed under the MIT License.

