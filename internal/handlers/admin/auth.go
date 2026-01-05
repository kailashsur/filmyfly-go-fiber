package admin

import (
	"filmyfly-go-fiber/internal/config"

	"github.com/gofiber/fiber/v2"
)

// GetAdminLogin renders the login page
func GetAdminLogin(c *fiber.Ctx) error {
	cfg := config.Load()

	return c.Render("admin/login", fiber.Map{
		"title": "Admin Login",
		"error": c.Query("error"),
		"firebaseConfig": fiber.Map{
			"apiKey":            cfg.FirebaseAPIKey,
			"authDomain":        cfg.FirebaseAuthDomain,
			"projectId":         cfg.FirebaseProjectID,
			"storageBucket":     cfg.FirebaseStorageBucket,
			"messagingSenderId": cfg.FirebaseMessagingSenderID,
			"appId":             cfg.FirebaseAppID,
		},
	})
}

// PostAdminLogin handles login (temporarily simplified - no Firebase verification)
func PostAdminLogin(c *fiber.Ctx) error {
	// For now, just redirect to dashboard
	// TODO: Implement proper Firebase authentication
	return c.Redirect("/admin")
}

// PostAdminLogout handles logout
func PostAdminLogout(c *fiber.Ctx) error {
	return c.Redirect("/admin/login")
}

// GetAdminDashboard renders the admin dashboard
func GetAdminDashboard(c *fiber.Ctx) error {
	// Mock user data for now
	user := map[string]interface{}{
		"email": "admin@filmyfly.work",
		"name":  "Admin User",
	}

	return c.Render("admin/dashboard", fiber.Map{
		"title": "Admin Dashboard",
		"user":  user,
	})
}
