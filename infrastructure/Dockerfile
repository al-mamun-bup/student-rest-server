# Use official Golang image as the base (updated to Go 1.23)
FROM golang:1.23-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY ../../go.mod ../../go.sum ./

# Download all dependencies (this will cache dependencies if go.mod and go.sum are not changed)
RUN go mod tidy

# Copy the rest of the source code
COPY ../../. .


# Build the Go app
RUN go build -o student-server .

# Expose the port the app will run on
EXPOSE 9090

# Command to run the application
CMD ["./student-server", "serve"]
