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

func main() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		host,
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
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
