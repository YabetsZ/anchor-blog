package main

import (
	"log"

	"anchor-blog/api"
	"anchor-blog/config"
	"anchor-blog/pkg/db"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Config loaded")

	// Connect to MongoDB
	mongoClient, err := db.Connect(cfg.Mongo.URI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Disconnect(mongoClient)
	log.Println("MongoDB connected")

	// Start Server
	router := api.SetupRouter()
	log.Printf("🚀 Server is running on port %s\n", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
