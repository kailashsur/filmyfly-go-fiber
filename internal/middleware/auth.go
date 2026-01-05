package middleware

import (
	"context"
	"strings"

	"filmyfly-go-fiber/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

// InitSession initializes the session store
func InitSession() {
	cfg := config.Load()
	Store = session.New(session.Config{
		KeyLookup:      "cookie:session_id",
		CookieSecure:   cfg.Environment == "production",
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})
}

// VerifyAdminToken middleware verifies Firebase ID token and manages session
func VerifyAdminToken(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return c.Redirect("/admin/login")
	}

	// Check if user is already authenticated via session
	adminUser := sess.Get("adminUser")
	if adminUser != nil {
		c.Locals("user", adminUser)
		return c.Next()
	}

	// Get token from cookie or Authorization header
	token := c.Cookies("adminToken")
	if token == "" {
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		return c.Redirect("/admin/login")
	}

	// Verify token with Firebase
	ctx := context.Background()
	decodedToken, err := config.FirebaseAuth.VerifyIDToken(ctx, token)
	if err != nil {
		// Clear invalid token
		c.ClearCookie("adminToken")
		sess.Delete("adminUser")
		sess.Save()
		return c.Redirect("/admin/login")
	}

	// Store user info in session
	userInfo := map[string]interface{}{
		"uid":   decodedToken.UID,
		"email": decodedToken.Claims["email"],
		"name":  decodedToken.Claims["name"],
	}
	sess.Set("adminUser", userInfo)
	sess.Save()

	c.Locals("user", userInfo)
	return c.Next()
}

// RedirectIfAuthenticated redirects to dashboard if already logged in
func RedirectIfAuthenticated(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return c.Next()
	}

	adminUser := sess.Get("adminUser")
	if adminUser != nil {
		return c.Redirect("/admin")
	}

	return c.Next()
}
