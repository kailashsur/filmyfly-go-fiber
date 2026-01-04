package middleware

import (
	"context"
	"strings"

	"filmyfly-go-fiber/internal/config"
	"filmyfly-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var SessionStore *session.Store

// InitSession initializes the session store
func InitSession(cfg *config.Config) {
	SessionStore = session.New(session.Config{
		CookieName:     "connect.sid",
		CookieSecure:   cfg.Environment == "production",
		CookieHTTPOnly: true,
		Expiration:     24 * 60 * 60, // 24 hours in seconds
	})
}

// VerifyAdminToken middleware verifies Firebase token for admin routes
func VerifyAdminToken(c *fiber.Ctx) error {
	sess, err := SessionStore.Get(c)
	if err != nil {
		return c.Redirect("/admin/login")
	}

	// Check if user is already authenticated via session
	adminUser := sess.Get("adminUser")
	if adminUser != nil {
		c.Locals("user", adminUser)
		return c.Next()
	}

	// Get token from Authorization header or cookie
	token := ""
	authHeader := c.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		token = c.Cookies("adminToken")
	}

	if token == "" {
		return c.Redirect("/admin/login")
	}

	// Verify token with Firebase
	decodedToken, err := config.FirebaseAuth.VerifyIDToken(context.Background(), token)
	if err != nil {
		// Clear invalid token
		c.ClearCookie("adminToken")
		sess.Delete("adminUser")
		sess.Save()

		utils.Error("Token verification error: %v", err)
		return c.Redirect("/admin/login")
	}

	// Store in session for future requests
	sess.Set("adminUser", decodedToken)
	sess.Save()

	c.Locals("user", decodedToken)
	return c.Next()
}

// RedirectIfAuthenticated redirects to dashboard if already logged in
func RedirectIfAuthenticated(c *fiber.Ctx) error {
	sess, err := SessionStore.Get(c)
	if err == nil {
		adminUser := sess.Get("adminUser")
		if adminUser != nil {
			return c.Redirect("/admin")
		}
	}
	return c.Next()
}
