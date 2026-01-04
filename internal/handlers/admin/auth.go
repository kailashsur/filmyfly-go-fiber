package admin

import (
	"context"

	"filmyfly-go-fiber/internal/config"
	"filmyfly-go-fiber/internal/middleware"
	"filmyfly-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// GetAdminLogin handles GET /admin/login
func GetAdminLogin(c *fiber.Ctx) error {
	return c.Render("admin/login", fiber.Map{
		"title": "Admin Login",
		"error": nil,
		"firebaseConfig": fiber.Map{
			"apiKey":            config.AppConfig.FirebaseAPIKey,
			"authDomain":        config.AppConfig.FirebaseAuthDomain,
			"projectId":         config.AppConfig.FirebaseProjectID,
			"storageBucket":     config.AppConfig.FirebaseStorageBucket,
			"messagingSenderId": config.AppConfig.FirebaseMessagingSenderID,
			"appId":             config.AppConfig.FirebaseAppID,
		},
	})
}

// PostAdminLogin handles POST /admin/login
func PostAdminLogin(c *fiber.Ctx) error {
	type LoginRequest struct {
		IDToken string `json:"idToken" form:"idToken"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Render("admin/login", fiber.Map{
			"title": "Admin Login",
			"error": "Invalid request",
			"firebaseConfig": fiber.Map{
				"apiKey":            config.AppConfig.FirebaseAPIKey,
				"authDomain":        config.AppConfig.FirebaseAuthDomain,
				"projectId":         config.AppConfig.FirebaseProjectID,
				"storageBucket":     config.AppConfig.FirebaseStorageBucket,
				"messagingSenderId": config.AppConfig.FirebaseMessagingSenderID,
				"appId":             config.AppConfig.FirebaseAppID,
			},
		})
	}

	if req.IDToken == "" {
		return c.Render("admin/login", fiber.Map{
			"title": "Admin Login",
			"error": "ID token is required",
			"firebaseConfig": fiber.Map{
				"apiKey":            config.AppConfig.FirebaseAPIKey,
				"authDomain":        config.AppConfig.FirebaseAuthDomain,
				"projectId":         config.AppConfig.FirebaseProjectID,
				"storageBucket":     config.AppConfig.FirebaseStorageBucket,
				"messagingSenderId": config.AppConfig.FirebaseMessagingSenderID,
				"appId":             config.AppConfig.FirebaseAppID,
			},
		})
	}

	// Verify the ID token with Firebase Admin
	decodedToken, err := config.FirebaseAuth.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		utils.Error("Login error: %v", err)
		return c.Render("admin/login", fiber.Map{
			"title": "Admin Login",
			"error": "Invalid credentials. Please try again.",
			"firebaseConfig": fiber.Map{
				"apiKey":            config.AppConfig.FirebaseAPIKey,
				"authDomain":        config.AppConfig.FirebaseAuthDomain,
				"projectId":         config.AppConfig.FirebaseProjectID,
				"storageBucket":     config.AppConfig.FirebaseStorageBucket,
				"messagingSenderId": config.AppConfig.FirebaseMessagingSenderID,
				"appId":             config.AppConfig.FirebaseAppID,
			},
		})
	}

	// Store user info in session
	sess, err := middleware.SessionStore.Get(c)
	if err != nil {
		return c.Status(500).SendString("Session error")
	}

	sess.Set("adminUser", decodedToken)
	sess.Save()

	// Set cookie for client-side
	c.Cookie(&fiber.Cookie{
		Name:     "adminToken",
		Value:    req.IDToken,
		HTTPOnly: true,
		Secure:   config.AppConfig.Environment == "production",
		MaxAge:   24 * 60 * 60, // 24 hours
	})

	return c.Redirect("/admin")
}

// PostAdminLogout handles POST /admin/logout
func PostAdminLogout(c *fiber.Ctx) error {
	sess, err := middleware.SessionStore.Get(c)
	if err == nil {
		sess.Destroy()
	}

	c.ClearCookie("adminToken")
	c.ClearCookie("connect.sid")

	return c.Redirect("/admin/login")
}

// GetAdminDashboard handles GET /admin
func GetAdminDashboard(c *fiber.Ctx) error {
	user := c.Locals("user")

	return c.Render("admin/dashboard", fiber.Map{
		"title": "Admin Dashboard",
		"user":  user,
	})
}
