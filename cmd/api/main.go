package main

import (
	"action-tracker/internal/shared/config"
	"action-tracker/internal/shared/database"
	"action-tracker/internal/shared/server"
	"log"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	httpServer := server.NewServer(cfg, db)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := httpServer.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
