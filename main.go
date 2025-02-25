package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"student-server/auth"
	"student-server/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Public route (No authentication required)
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Protected routes (Require authentication)
	protectedRoutes := router.PathPrefix("/students").Subrouter()
	protectedRoutes.Use(auth.BasicAuthMiddleware) // Apply auth middleware
	protectedRoutes.HandleFunc("", handlers.GetStudentsHandler).Methods("GET")
	protectedRoutes.HandleFunc("", handlers.AddStudentHandler).Methods("POST")
	protectedRoutes.HandleFunc("/{id}", handlers.GetStudentByIDHandler).Methods("GET")
	protectedRoutes.HandleFunc("/{id}", handlers.UpdateStudentHandler).Methods("PUT")
	protectedRoutes.HandleFunc("/{id}", handlers.DeleteStudentHandler).Methods("DELETE")

	// Create an HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM) // Capture Ctrl+C or kill signal

	// Run server in a goroutine
	go func() {
		log.Println("Server running on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Wait for termination signal
	<-stop
	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	log.Println("Server exited gracefully")
}
