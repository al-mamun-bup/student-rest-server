# Use official Golang image as the base (updated to Go 1.23)
FROM golang:1.23-alpine

# Set the Current Working Directory inside the container
WORKDIR /student-server

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o student-server .

# Expose port 9090 (as you're running the server on port 9090)
EXPOSE 9090

# Command to run the executable with dynamic port
CMD ["./student-server", "serve"]
