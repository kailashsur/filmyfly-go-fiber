package config

import (
	"context"
	"log"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

// InitializeFirebase initializes Firebase Admin SDK
func InitializeFirebase() error {
	ctx := context.Background()
	cfg := Load()

	// Try to initialize with service account key file if it's a valid file path
	if cfg.FirebaseServiceAccountKey != "" {
		// Check if it's a file path (not JSON content)
		if len(cfg.FirebaseServiceAccountKey) < 500 && !strings.Contains(cfg.FirebaseServiceAccountKey, "{") {
			opt := option.WithCredentialsFile(cfg.FirebaseServiceAccountKey)
			app, err := firebase.NewApp(ctx, nil, opt)
			if err != nil {
				log.Printf("⚠️  Firebase initialization with file failed: %v", err)
				return err
			}

			client, err := app.Auth(ctx)
			if err != nil {
				log.Printf("⚠️  Firebase Auth client creation failed: %v", err)
				return err
			}

			FirebaseAuth = client
			log.Println("✅ Firebase initialized with service account key file")
			return nil
		}

		// If it looks like JSON content, try to parse it
		log.Println("⚠️  Firebase service account key appears to be JSON content, not a file path")
		log.Println("⚠️  Please set FIREBASE_SERVICE_ACCOUNT_KEY to a file path, not JSON content")
	}

	// Try with default credentials or environment variables
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Printf("⚠️  Firebase initialization failed: %v", err)
		log.Println("⚠️  Admin authentication will not work without Firebase")
		return err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("⚠️  Firebase Auth client creation failed: %v", err)
		return err
	}

	FirebaseAuth = client
	log.Println("✅ Firebase initialized with default credentials")
	return nil
}
