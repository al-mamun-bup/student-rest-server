package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"student-server/auth"
	"student-server/database"
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
	// dbHost := "localhost"
	// dbUser := "mamun"
	// dbPassword := "1234"
	// dbName := "student_db"
	// dbPort := "5432"

	// PostgreSQL connection string
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	// 	dbHost, dbUser, dbPassword, dbName, dbPort)

	// // Connect to PostgreSQL using GORM
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatalf("❌ Failed to connect to database: %v", err)
	// }

	database.ConnectDB()
	db := database.DB

	// Drop and re-create the students table (if needed)
	// This will drop the existing table and create a new one based on the current model
	// if err := db.Migrator().DropTable(&models.Student{}); err != nil {
	// 	log.Fatalf("❌ Failed to drop table: %v", err)
	// }
	// if err := db.AutoMigrate(&models.Student{}); err != nil {
	// 	log.Fatalf("❌ Failed to migrate database: %v", err)
	// }
	log.Println("✅ Connected to PostgreSQL & migrated successfully!")

	// Set the database instance in handlers
	handlers.SetDB(db)

	// Read the PORT environment variable if set
	portStr := os.Getenv("PORT")
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("Invalid port value: %v", err)
		}
	}

	router := mux.NewRouter()

	// Public route
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Protected routes (require authentication)
	protectedRoutes := router.PathPrefix("/students").Subrouter()
	protectedRoutes.Use(auth.BasicAuthMiddleware)
	protectedRoutes.HandleFunc("", handlers.GetStudentsHandler).Methods("GET")
	protectedRoutes.HandleFunc("", handlers.AddStudentHandler).Methods("POST")
	protectedRoutes.HandleFunc("/{id}", handlers.GetStudentByIDHandler).Methods("GET")
	protectedRoutes.HandleFunc("/{id}", handlers.UpdateStudentHandler).Methods("PUT")
	protectedRoutes.HandleFunc("/{id}", handlers.DeleteStudentHandler).Methods("DELETE")

	address := fmt.Sprintf("0.0.0.0:%d", port)
	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🚀 Server running on port %d...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(" ListenAndServe error: %v", err)
		}
	}()

	<-stop
	log.Println("⚠️  Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}
	log.Println("✅ Server exited gracefully")
}
