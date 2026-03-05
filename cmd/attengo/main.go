package main

import (
	"log"

	"github.com/fernando8franco/attengo/internal/config"
	"github.com/fernando8franco/attengo/internal/db"
	"github.com/fernando8franco/attengo/internal/routes"
	"github.com/fernando8franco/attengo/internal/service"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg := config.Load()

	conn, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close()

	rhSvc := service.NewRequiredHourService(conn)

	router := routes.SetupRouter(rhSvc, cfg)

	log.Printf("Server starting on %s", cfg.Port)
	if err := router.Run(cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
