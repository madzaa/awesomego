package main

import (
	"awesomeProject/handlers"
	"awesomeProject/repository"
	"awesomeProject/service"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	host := getEnvWithDefault("DB_HOST", "localhost")
	user := getEnvWithDefault("DB_USER", "user")
	password := getEnvWithDefault("DB_PASSWORD", "password")
	port := getEnvWithDefault("DB_PORT", "5432")
	dbName := getEnvWithDefault("DB_NAME", "core")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	transactionRepo := repository.NewTransactionRepository(db)
	balanceRepo := repository.NewUserBalanceRepository(db)

	svc := service.NewServiceImpl(transactionRepo, balanceRepo)
	h := handlers.NewHandler(svc)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      http.HandlerFunc(h.ServeHTTP),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server error:", err)
	}
}
