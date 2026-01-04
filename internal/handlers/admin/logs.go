package admin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetLogs handles GET /admin/logs
func GetLogs(c *fiber.Ctx) error {
	logType := c.Query("type", "app") // app or error

	var logFile string
	if logType == "error" {
		logFile = filepath.Join("logs", "error.log")
	} else {
		logFile = filepath.Join("logs", "app.log")
	}

	// Read log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		content = []byte("No logs available yet.")
	}

	// Split into lines and reverse (newest first)
	lines := strings.Split(string(content), "\n")

	// Reverse the lines
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	// Take last 500 lines
	if len(lines) > 500 {
		lines = lines[:500]
	}

	// Join back
	logContent := strings.Join(lines, "\n")

	return c.Render("admin/logs", fiber.Map{
		"title":      "System Logs",
		"logContent": logContent,
		"logType":    logType,
	})
}
