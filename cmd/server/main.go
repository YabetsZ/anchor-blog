package main

import (
	"context"
	"log"
	"time"

	"anchor-blog/api"
	"anchor-blog/api/handler"
	"anchor-blog/api/handler/content"
	"anchor-blog/api/handler/oauth"
	"anchor-blog/api/handler/post"
	"anchor-blog/api/handler/user"
	"anchor-blog/config"
	"anchor-blog/internal/repository/gemini"
	postrepo "anchor-blog/internal/repository/post"
	tokenrepo "anchor-blog/internal/repository/token"
	userrepo "anchor-blog/internal/repository/user"
	contentsvc "anchor-blog/internal/service/content"
	postsvc "anchor-blog/internal/service/post"
	usersvc "anchor-blog/internal/service/user"
	viewsvc "anchor-blog/internal/service/view"
	"anchor-blog/pkg/db"
	redisclient "anchor-blog/pkg/redis"
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

	// Initialize Redis client
	redisClient := redisclient.NewRedisClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx); err != nil {
		log.Printf("⚠️  Redis connection failed: %v (continuing without Redis)", err)
		redisClient = nil // Set to nil so handlers can handle gracefully
	} else {
		log.Println("✅ Redis connected")
	}

	// Initialize repositories
	userRepository := userrepo.NewUserRepository(userCollection)
	tokenRepository := tokenrepo.NewMongoTokenRepository(tokenCollection)
	postRepository := postrepo.NewMongoPostRepository(postCollection)
	activationTokenRepo := tokenrepo.NewActivationTokenRepository(activationTokenCollection)
	passwordResetTokenRepo := tokenrepo.NewPasswordResetTokenRepository(passwordResetTokenCollection)

	// Initialize services
	activationService := usersvc.NewActivationService(userRepository, activationTokenRepo)
	passwordResetService := usersvc.NewPasswordResetService(userRepository, passwordResetTokenRepo)

	// Initialize view tracking service (with Redis if available)
	var viewTrackingService *viewsvc.ViewTrackingService
	if redisClient != nil {
		viewTrackingService = viewsvc.NewViewTrackingService(redisClient, postRepository, cfg.Redis.ViewTrackingTTL)
		log.Println("✅ View tracking service initialized with Redis")
	} else {
		log.Println("⚠️  View tracking service disabled (Redis unavailable)")
	}

	// Initialize handlers
	userHandler := user.NewUserHandler(usersvc.NewUserServices(userRepository, tokenRepository, cfg))
	postHandler := post.NewPostHandler(postsvc.NewPostService(postRepository), viewTrackingService)
	activationHandler := handler.NewActivationHandler(activationService)
	passwordResetHandler := handler.NewPasswordResetHandler(passwordResetService)
	contentHandler := content.NewContentHandler(contentsvc.NewContentUsecase(gemini.NewGeminiRepo(cfg.GenAI.GeminiAPIKey, cfg.GenAI.GeminiModel)))

	oauth.InitializeGoogleOAuthConfig(cfg)
	oauthHandler := oauth.NewOAuthHandler(usersvc.NewUserServices(userRepository, tokenRepository, cfg))

	// Start Server
	router := api.SetupRouter(cfg, userHandler, postHandler, activationHandler, passwordResetHandler, contentHandler, oauthHandler)
	log.Printf("🚀 Server is running on port %s\n", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
