package middleware

import (
	"filmyfly-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler handles errors and returns appropriate responses
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	utils.Error("Error %d: %s", code, message)

	// Return JSON for API routes
	if len(c.Path()) > 4 && c.Path()[:4] == "/api" {
		return c.Status(code).JSON(fiber.Map{
			"success": false,
			"error":   message,
		})
	}

	// Return HTML for other routes
	return c.Status(code).Render("error", fiber.Map{
		"title":      code,
		"message":    message,
		"statusCode": code,
	})
}
