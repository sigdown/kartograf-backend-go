package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/sigdown/kartograf-backend-go/internal/config"
	"github.com/sigdown/kartograf-backend-go/internal/db"
)

func main() {
	_ = godotenv.Load()

	cfg := config.MustLoad()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pool, err := db.NewPostgresPool(context.Background(), cfg.Postgres.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	logger.Info("postgres connected")
}