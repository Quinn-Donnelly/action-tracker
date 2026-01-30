package main

import (
	"log"
	"net/http"

	"action-tracker/internal/config"
	"action-tracker/internal/db"
	"action-tracker/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	if err := db.Migrate(database); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	actionHandler := handlers.NewActionHandler(database)

	r := gin.Default()

	r.POST("/api/actions", actionHandler.CreateAction)
	r.GET("/api/actions", actionHandler.GetActions)
	r.GET("/api/actions/:id", actionHandler.GetAction)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
