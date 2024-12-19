package config

import (
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DB     DB
	Server Server
}

type DB struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type Server struct {
	Host        string
	Port        string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

const (
	DefaultTimeout     = 10 * time.Second
	DefaultIdleTimeout = 60 * time.Second
)

func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		DB: DB{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
		},
		Server: Server{
			Host: os.Getenv("SRV_HOST"),
			Port: os.Getenv("SRV_PORT"),
		},
	}

	timeout := DefaultTimeout
	if val := os.Getenv("SRV_TIMEOUT"); val != "" {
		parsed, err := time.ParseDuration(val)
		if err != nil {
			return nil, fmt.Errorf("invalid SRV_TIMEOUT: %v", err)
		}
		timeout = parsed
	}

	idleTimeout := DefaultIdleTimeout
	if val := os.Getenv("SRV_IDLE_TIMEOUT"); val != "" {
		parsed, err := time.ParseDuration(val)
		if err != nil {
			return nil, fmt.Errorf("invalid SRV_IDLE_TIMEOUT: %v", err)
		}
		idleTimeout = parsed
	}

	cfg.Server.Timeout = timeout
	cfg.Server.IdleTimeout = idleTimeout

	return cfg, nil
}
