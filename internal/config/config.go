package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port        string
	Environment string

	// Database
	DatabaseURL string

	// Session
	SessionSecret string

	// Frontend
	FrontendURL string

	// Firebase
	FirebaseProjectID         string
	FirebaseClientEmail       string
	FirebasePrivateKey        string
	FirebaseAPIKey            string
	FirebaseAuthDomain        string
	FirebaseStorageBucket     string
	FirebaseMessagingSenderID string
	FirebaseAppID             string
	FirebaseServiceAccountKey string

	// Logging
	PrismaLogQueries bool
}

var AppConfig *Config

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		Port:        getEnv("PORT", "3000"),
		Environment: getEnv("NODE_ENV", "development"),

		DatabaseURL: getEnv("DATABASE_URL", ""),

		SessionSecret: getEnv("SESSION_SECRET", "your-secret-key-change-in-production"),

		FrontendURL: getEnv("FRONTEND_URL", ""),

		FirebaseProjectID:         os.Getenv("FIREBASE_PROJECT_ID"),
		FirebaseClientEmail:       os.Getenv("FIREBASE_CLIENT_EMAIL"),
		FirebasePrivateKey:        os.Getenv("FIREBASE_PRIVATE_KEY"),
		FirebaseAPIKey:            os.Getenv("FIREBASE_API_KEY"),
		FirebaseAuthDomain:        os.Getenv("FIREBASE_AUTH_DOMAIN"),
		FirebaseStorageBucket:     os.Getenv("FIREBASE_STORAGE_BUCKET"),
		FirebaseMessagingSenderID: os.Getenv("FIREBASE_MESSAGING_SENDER_ID"),
		FirebaseAppID:             os.Getenv("FIREBASE_APP_ID"),
		FirebaseServiceAccountKey: os.Getenv("FIREBASE_SERVICE_ACCOUNT_KEY"),

		PrismaLogQueries: getEnv("PRISMA_LOG_QUERIES", "false") == "true",
	}

	// Validate required configuration
	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	AppConfig = config
	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
