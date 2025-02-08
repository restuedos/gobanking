package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   *slog.Logger
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Load() *Config {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Parse command line arguments
	host := flag.String("host", "0.0.0.0", "API server host")
	port := flag.Int("port", 3000, "API server port")
	flag.Parse()

	return &Config{
		Server: ServerConfig{
			Host: *host,
			Port: *port,
		},
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		Logger: logger,
	}
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.Name, c.Port,
	)
}
