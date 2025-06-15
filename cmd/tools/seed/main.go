package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed seed.sql
var seedSQL string

func main() {
	connString := os.Getenv("GOMANIA_CONNECTION_STRING")
	if connString == "" {
		log.Fatalf("Connection string is empty, please set env variable GOMANIA_CONNECTION_STRING")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	fmt.Println("Seeding database...")

	// Execute the embedded SQL file
	_, err = pool.Exec(ctx, seedSQL)
	if err != nil {
		log.Fatalf("Failed to execute seed SQL: %v", err)
	}

	fmt.Println("Database seeding completed successfully")
}
