package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khatibomar/gomania/internal/database"
)

type config struct {
	port int
	env  string
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	ctx    context.Context
	config config
	logger *slog.Logger
	db     database.Querier
}

func parseFlags(cfg *config) {
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	flag.Parse()
}

func main() {
	var cfg config

	parseFlags(&cfg)

	connString := os.Getenv("GOMANIA_CONNECTION_STRING")
	if connString == "" {
		log.Fatalf("Connection string is empty, please set env variable GOMANIA_CONNECTION_STRING")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	db := database.New(pool)

	app := &application{
		ctx:    ctx,
		config: cfg,
		logger: logger,
		db:     db,
	}

	if err = app.serve(); err != nil {
		log.Fatalf("failed to start listening on server: %v", err)
	}
}
