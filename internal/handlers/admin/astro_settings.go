package admin

import (
	"filmyfly-go-fiber/internal/database"
	"filmyfly-go-fiber/internal/database/models"

	"github.com/gofiber/fiber/v2"
)

// GetAstroSettings displays Astro settings page
func GetAstroSettings(c *fiber.Ctx) error {
	success := c.Query("success", "")
	errorMsg := c.Query("error", "")

	// Fetch all Astro settings
	var astroSettingsData []models.AstroSetting
	database.DB.Find(&astroSettingsData)

	// Convert to map for easier template access
	settings := make(map[string]string)
	for _, s := range astroSettingsData {
		settings[s.Key] = s.Value
	}

	return c.Render("admin/astro-settings", fiber.Map{
		"title":    "Astro Website Settings",
		"settings": settings,
		"success":  success,
		"error":    errorMsg,
		"user":     getUserFromSession(c),
	})
}

// PostAstroSettings handles Astro settings update
func PostAstroSettings(c *fiber.Ctx) error {
	// Get all form values
	formData := make(map[string]string)
	c.Request().PostArgs().VisitAll(func(key, value []byte) {
		formData[string(key)] = string(value)
	})

	// Update each setting in astro_settings table
	for key, value := range formData {
		var setting models.AstroSetting

		// Find or create setting
		result := database.DB.Where("`key` = ?", key).First(&setting)

		if result.Error != nil {
			// Create new setting
			setting = models.AstroSetting{
				Key:   key,
				Value: value,
			}
			database.DB.Create(&setting)
		} else {
			// Update existing setting
			setting.Value = value
			database.DB.Save(&setting)
		}
	}

	return c.Redirect("/admin/astro-settings?success=Astro settings updated successfully")
}
