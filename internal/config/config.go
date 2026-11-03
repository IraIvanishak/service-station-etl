// Package config reads the server settings from the environment.
package config

import "os"

type Config struct {
	// DatabaseURL is a lib/pq connection string.
	DatabaseURL string
	// Addr is the address the HTTP server listens on.
	Addr string
	// CORSOrigins is a comma-separated list of origins allowed to call the API.
	CORSOrigins string
}

// Defaults match docker-compose.yml and the Vite dev server.
const (
	defaultDatabaseURL = "postgres://postgres:postgres@localhost:5452/service_station?sslmode=disable"
	defaultAddr        = ":3000"
	defaultCORSOrigins = "http://localhost:5173,http://127.0.0.1:5173"
)

func Load() Config {
	return Config{
		DatabaseURL: getenv("DATABASE_URL", defaultDatabaseURL),
		Addr:        getenv("HTTP_ADDR", defaultAddr),
		CORSOrigins: getenv("CORS_ORIGINS", defaultCORSOrigins),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
