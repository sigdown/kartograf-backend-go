package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	apphttp "github.com/sigdown/kartograf-backend-go/internal/http"

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

	router := apphttp.NewRouter()

	addr := cfg.App.Host + ":" + cfg.App.Port
	logger.Info("http server starting", "addr", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}