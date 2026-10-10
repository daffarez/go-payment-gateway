package main

import (
	"log"
	"net/http"
	"os"

	"go-payment-gateway/internal/database"
	"go-payment-gateway/internal/handler"
	customMiddleware "go-payment-gateway/internal/middleware"
	"go-payment-gateway/internal/repository"
	"go-payment-gateway/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := database.InitDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Dependency Injection
	trxRepo := repository.NewPostgresTransactionRepository(db)
	trxService := service.NewTransactionService(db, trxRepo)
	trxHandler := handler.NewTransactionHandler(trxService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Public Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Payment Gateway Service is Running"))
	})

	// Protected Routes (Membutuhkan X-Api-Key & X-Secret-Key)
	r.Group(func(protected chi.Router) {
		protected.Use(customMiddleware.AuthMiddleware(db))

		protected.Post("/v1/transactions", trxHandler.CreateTransaction)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
