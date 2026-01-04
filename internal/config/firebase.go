package config

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

// InitFirebase initializes Firebase Admin SDK
func InitFirebase(cfg *Config) error {
	var opt option.ClientOption

	// Option 1: Using service account JSON
	if cfg.FirebaseServiceAccountKey != "" {
		opt = option.WithCredentialsJSON([]byte(cfg.FirebaseServiceAccountKey))
	} else if cfg.FirebaseProjectID != "" {
		// Option 2: Using individual environment variables
		privateKey := strings.ReplaceAll(cfg.FirebasePrivateKey, "\\n", "\n")

		serviceAccount := map[string]interface{}{
			"type":         "service_account",
			"project_id":   cfg.FirebaseProjectID,
			"private_key":  privateKey,
			"client_email": cfg.FirebaseClientEmail,
			"token_uri":    "https://oauth2.googleapis.com/token",
		}

		jsonData, err := json.Marshal(serviceAccount)
		if err != nil {
			return fmt.Errorf("failed to marshal service account: %w", err)
		}

		opt = option.WithCredentialsJSON(jsonData)
	} else {
		return fmt.Errorf("firebase credentials not configured")
	}

	// Initialize Firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	// Get Auth client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get firebase auth client: %w", err)
	}

	FirebaseAuth = authClient
	return nil
}
