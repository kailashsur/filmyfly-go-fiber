package admin

import (
	"filmyfly-go-fiber/internal/database"
	"filmyfly-go-fiber/internal/database/models"

	"github.com/gofiber/fiber/v2"
)

// GetSettings handles GET /admin/settings
func GetSettings(c *fiber.Ctx) error {
	var settings []models.Setting
	database.DB.Order("`key` ASC").Find(&settings)

	return c.Render("admin/settings", fiber.Map{
		"title":    "Settings",
		"settings": settings,
		"success":  c.Query("success"),
		"error":    c.Query("error"),
	})
}

// PostSettings handles POST /admin/settings
func PostSettings(c *fiber.Ctx) error {
	// Check if adding a new setting
	newKey := c.FormValue("newKey")
	newValue := c.FormValue("newValue")
	newDescription := c.FormValue("newDescription")

	if newKey != "" && newValue != "" {
		// Create new setting
		setting := models.Setting{
			Key:         newKey,
			Value:       newValue,
			Description: stringPtr(newDescription),
		}
		if err := database.DB.Create(&setting).Error; err != nil {
			return c.Redirect("/admin/settings?error=Failed to create setting")
		}
		return c.Redirect("/admin/settings?success=Setting added successfully")
	}

	// Update existing settings
	var settings []models.Setting
	database.DB.Find(&settings)

	for _, setting := range settings {
		value := c.FormValue(setting.Key)
		if value != "" && value != setting.Value {
			setting.Value = value
			database.DB.Save(&setting)
		}
	}

	return c.Redirect("/admin/settings?success=Settings updated successfully")
}

// GetAstroSettings handles GET /admin/astro-settings
func GetAstroSettings(c *fiber.Ctx) error {
	var settings []models.AstroSetting
	database.DB.Order("`key` ASC").Find(&settings)

	return c.Render("admin/astro-settings", fiber.Map{
		"title":    "Astro Settings",
		"settings": settings,
		"success":  c.Query("success"),
		"error":    c.Query("error"),
	})
}

// PostAstroSettings handles POST /admin/astro-settings
func PostAstroSettings(c *fiber.Ctx) error {
	// Check if adding a new setting
	newKey := c.FormValue("newKey")
	newValue := c.FormValue("newValue")
	newDescription := c.FormValue("newDescription")

	if newKey != "" && newValue != "" {
		// Create new setting
		setting := models.AstroSetting{
			Key:         newKey,
			Value:       newValue,
			Description: stringPtr(newDescription),
		}
		if err := database.DB.Create(&setting).Error; err != nil {
			return c.Redirect("/admin/astro-settings?error=Failed to create setting")
		}
		return c.Redirect("/admin/astro-settings?success=Setting added successfully")
	}

	// Update existing settings
	var settings []models.AstroSetting
	database.DB.Find(&settings)

	for _, setting := range settings {
		value := c.FormValue(setting.Key)
		if value != "" && value != setting.Value {
			setting.Value = value
			database.DB.Save(&setting)
		}
	}

	return c.Redirect("/admin/astro-settings?success=Astro settings updated successfully")
}

// GetSystemCheck handles GET /admin/system-check
func GetSystemCheck(c *fiber.Ctx) error {
	user := c.Locals("user")

	// Get database size
	type SizeResult struct {
		Size string
	}
	var result SizeResult
	database.DB.Raw("SELECT pg_size_pretty(pg_database_size(current_database())) as size").Scan(&result)

	databaseSize := result.Size
	if databaseSize == "" {
		databaseSize = "Unknown"
	}

	return c.Render("admin/system-check", fiber.Map{
		"title": "System Check",
		"user":  user,
		"databaseStatus": fiber.Map{
			"size":        databaseSize,
			"lastChecked": "Just now",
		},
	})
}
