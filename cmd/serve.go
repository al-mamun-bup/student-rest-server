package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"student-server/auth"
	"student-server/handlers"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the student REST API server",
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
}

func startServer() {
	router := mux.NewRouter()

	// Public route
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Protected routes
	protectedRoutes := router.PathPrefix("/students").Subrouter()
	protectedRoutes.Use(auth.BasicAuthMiddleware)
	protectedRoutes.HandleFunc("", handlers.GetStudentsHandler).Methods("GET")
	protectedRoutes.HandleFunc("", handlers.AddStudentHandler).Methods("POST")
	protectedRoutes.HandleFunc("/{id}", handlers.GetStudentByIDHandler).Methods("GET")
	protectedRoutes.HandleFunc("/{id}", handlers.UpdateStudentHandler).Methods("PUT")
	protectedRoutes.HandleFunc("/{id}", handlers.DeleteStudentHandler).Methods("DELETE")

	address := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server running on port %d...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}
	log.Println("Server exited gracefully")
}
