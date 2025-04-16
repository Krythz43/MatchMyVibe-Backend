package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/matchmyvibe/backend/internal/auth"
	"github.com/matchmyvibe/backend/internal/db"
	"github.com/matchmyvibe/backend/internal/handlers"
	"github.com/matchmyvibe/backend/internal/middleware"
	"github.com/matchmyvibe/backend/internal/spotify"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Set up database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "matchmyvibe")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dbConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	database, err := db.New(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set up JWT service
	jwtSecret := getEnv("JWT_SECRET", "supersecret")
	jwtDuration := 24 * time.Hour // Token valid for 24 hours
	jwtService := auth.New(jwtSecret, jwtDuration)

	// Set up Spotify client
	spotifyClientID := getEnv("SPOTIFY_CLIENT_ID", "")
	spotifyClientSecret := getEnv("SPOTIFY_CLIENT_SECRET", "")
	spotifyRedirectURI := getEnv("SPOTIFY_REDIRECT_URI", "")
	spotifyClient := spotify.New(spotifyClientID, spotifyClientSecret, spotifyRedirectURI)

	// Set up handlers
	authHandler := &handlers.AuthHandler{
		DB:            database,
		JWTService:    jwtService,
		SpotifyClient: spotifyClient,
	}

	profileHandler := &handlers.ProfileHandler{
		DB:            database,
		SpotifyClient: spotifyClient,
	}

	// Set up router
	router := gin.Default()

	// Set up routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/spotify", authHandler.SpotifyAuth)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.AuthMiddleware(jwtService))
	{
		// Profile routes
		protectedRoutes.GET("/profile", profileHandler.GetProfile)
		protectedRoutes.PUT("/profile", profileHandler.UpdateProfile)
		protectedRoutes.PUT("/profile/currently-playing", profileHandler.UpdateCurrentlyPlaying)
	}

	// Start the server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
