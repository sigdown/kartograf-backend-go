package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	apphttp "github.com/sigdown/kartograf-backend-go/internal/http"

	apiauth "github.com/sigdown/kartograf-backend-go/internal/auth"
	"github.com/sigdown/kartograf-backend-go/internal/config"
	"github.com/sigdown/kartograf-backend-go/internal/db"
	"github.com/sigdown/kartograf-backend-go/internal/repository"
	"github.com/sigdown/kartograf-backend-go/internal/storage"
	"github.com/sigdown/kartograf-backend-go/internal/usecase"
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

	s3Storage, err := storage.NewS3Storage(
		cfg.S3.Endpoint,
		cfg.S3.Region,
		cfg.S3.AccessKey,
		cfg.S3.SecretKey,
		cfg.S3.UsePathStyle,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := s3Storage.EnsureBucket(context.Background(), cfg.S3.Bucket); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewPostgresUserRepository(pool)
	pointRepo := repository.NewPostgresPointRepository(pool)
	mapRepo := repository.NewPostgresMapRepository(pool)

	tokenManager := apiauth.NewTokenManager(cfg.Auth.JWTSecret, cfg.Auth.AccessTokenTTL)

	router := apphttp.NewRouter(apphttp.Services{
		Auth:   usecase.NewAuthService(userRepo, tokenManager),
		Points: usecase.NewPointService(pointRepo),
		Maps: usecase.NewMapService(
			mapRepo,
			s3Storage,
			cfg.S3.Bucket,
			cfg.S3.PresignUploadTTL,
			cfg.S3.PresignDownloadTTL,
			cfg.S3.ProxyEnabled,
			cfg.S3.UploadProxyURL,
			cfg.S3.DownloadProxyURL,
		),
		Tokens: tokenManager,
	})

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
