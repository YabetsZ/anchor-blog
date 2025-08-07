package main

import (
	"log"

	"anchor-blog/api"
	"anchor-blog/api/handler/post"
	"anchor-blog/api/handler/user"
	"anchor-blog/config"
	postrepo "anchor-blog/internal/repository/post"
	tokenrepo "anchor-blog/internal/repository/token"
	userrepo "anchor-blog/internal/repository/user"
	"anchor-blog/internal/service"
	usersvc "anchor-blog/internal/service/user"
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

	// Initialize collections
	userCollection := mongoClient.Database(cfg.Mongo.Database).Collection(cfg.Mongo.UserCollection)
	tokenCollection := mongoClient.Database(cfg.Mongo.Database).Collection(cfg.Mongo.TokenCollection)
	postCollection := mongoClient.Database(cfg.Mongo.Database).Collection(cfg.Mongo.PostCollection)
	activationTokenCollection := mongoClient.Database(cfg.Mongo.Database).Collection("activation_tokens")
	passwordResetTokenCollection := mongoClient.Database(cfg.Mongo.Database).Collection("password_reset_tokens")

	// Initialize repositories
	userRepository := userrepo.NewUserRepository(userCollection)
	activationTokenRepo := tokenrepo.NewActivationTokenRepository(activationTokenCollection)
	passwordResetTokenRepo := tokenrepo.NewPasswordResetTokenRepository(passwordResetTokenCollection)

	// Initialize services
	activationService := usersvc.NewActivationService(userRepository, activationTokenRepo)
	passwordResetService := usersvc.NewPasswordResetService(userRepository, passwordResetTokenRepo)

	// Initialize handlers
	userHandler := user.NewUserHandler(usersvc.NewUserServices(userrepo.NewUserRepository(userCollection), tokenrepo.NewMongoTokenRepository(tokenCollection), cfg))
	postHandler := post.NewPostHandler(service.NewPostService(postrepo.NewMongoPostRepository(postCollection)))
	activationHandler := handler.NewActivationHandler(activationService)
	passwordResetHandler := handler.NewPasswordResetHandler(passwordResetService)

	// Start Server
	router := api.SetupRouter(cfg, userHandler, postHandler, activationHandler, passwordResetHandler)
	log.Printf("🚀 Server is running on port %s\n", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
