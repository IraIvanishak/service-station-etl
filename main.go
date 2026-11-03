package main

import (
	"database/sql"
	"log"

	"github.com/IraIvanishak/service-station-etl/internal/api"
	"github.com/IraIvanishak/service-station-etl/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("database is not reachable (DATABASE_URL): %v", err)
	}

	app := api.New(db, cfg.CORSOrigins)
	log.Fatal(app.Listen(cfg.Addr))
}
