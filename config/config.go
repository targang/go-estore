package config

import (
	"fmt"
	"net"
	"os"
)

type (
	Config struct {
		GRPC
		PG
		Admin
	}

	GRPC struct {
		Port string `env:"GRPC_PORT"`
	}

	PG struct {
		URL      string
		Host     string `env:"POSTGRES_HOST"`
		Port     string `env:"POSTGRES_PORT"`
		DB       string `env:"POSTGRES_DB"`
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
	}

	Admin struct {
		Username     string `env:"ADMIN_USERNAME"`
		PasswordHash string `env:"ADMIN_PASSWORD_HASH"`
		JWTSecret    string `env:"ADMIN_JWT_SECRET"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	cfg.GRPC.Port = os.Getenv("GRPC_PORT")

	cfg.PG.Host = os.Getenv("POSTGRES_HOST")
	cfg.PG.Port = os.Getenv("POSTGRES_PORT")
	cfg.PG.DB = os.Getenv("POSTGRES_DB")
	cfg.PG.User = os.Getenv("POSTGRES_USER")
	cfg.PG.Password = os.Getenv("POSTGRES_PASSWORD")

	cfg.Admin.Username = os.Getenv("ADMIN_USERNAME")
	cfg.Admin.PasswordHash = os.Getenv("ADMIN_PASSWORD_HASH")
	cfg.Admin.JWTSecret = os.Getenv("ADMIN_JWT_SECRET")

	cfg.PG.URL = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.PG.User,
		cfg.PG.Password,
		net.JoinHostPort(cfg.PG.Host, cfg.PG.Port),
		cfg.PG.DB,
	)

	return cfg, nil
}
