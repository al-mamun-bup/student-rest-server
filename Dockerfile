# Use official Golang image as the base
FROM golang:1.21-alpine

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

# Expose port (you can expose the default 8080 port)
EXPOSE 8080

# Default to 8080, but can be overridden via environment variable
#ENV PORT=8080

# Command to run the executable with dynamic port
CMD ["./student-server", "serve"]
