package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App struct {
		Env  string `env:"APP_ENV" env-default:"local"`
		Host string `env:"APP_HOST" env-default:"0.0.0.0"`
		Port string `env:"APP_PORT" env-default:"8080"`
	}

	Postgres struct {
		DSN string `env:"POSTGRES_DSN" env-required:"true"`
	}

	S3 struct {
		Endpoint     string `env:"S3_ENDPOINT" env-required:"true"`
		Region       string `env:"S3_REGION" env-required:"true"`
		AccessKey    string `env:"S3_ACCESS_KEY" env-required:"true"`
		SecretKey    string `env:"S3_SECRET_KEY" env-required:"true"`
		Bucket       string `env:"S3_BUCKET" env-required:"true"`
		UsePathStyle bool   `env:"S3_USE_PATH_STYLE" env-default:"true"`
	}

	Auth struct {
		JWTSecret      string        `env:"AUTH_JWT_SECRET" env-required:"true"`
		AccessTokenTTL time.Duration `env:"AUTH_ACCESS_TOKEN_TTL" env-default:"15m"`
	}
}

func MustLoad() Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}
	return cfg
}
