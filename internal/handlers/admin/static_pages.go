package admin

import (
	"strconv"

	"filmyfly-go-fiber/internal/database"
	"filmyfly-go-fiber/internal/database/models"
	"filmyfly-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// GetStaticPageList handles GET /admin/static-pages
func GetStaticPageList(c *fiber.Ctx) error {
	var pages []models.StaticPage
	database.DB.Order("\"createdAt\" DESC").Find(&pages)

	return c.Render("admin/static-pages/list", fiber.Map{
		"title":   "Manage Static Pages",
		"pages":   pages,
		"success": c.Query("success"),
		"error":   c.Query("error"),
	})
}

// GetAddStaticPage handles GET /admin/static-pages/add
func GetAddStaticPage(c *fiber.Ctx) error {
	return c.Render("admin/static-pages/add", fiber.Map{
		"title": "Add Static Page",
		"page":  nil,
		"error": nil,
	})
}

// PostAddStaticPage handles POST /admin/static-pages/add
func PostAddStaticPage(c *fiber.Ctx) error {
	type PageForm struct {
		Title           string `form:"title"`
		Slug            string `form:"slug"`
		Content         string `form:"content"`
		MetaTitle       string `form:"metaTitle"`
		MetaDescription string `form:"metaDescription"`
		MetaKeywords    string `form:"metaKeywords"`
		IsPublished     string `form:"isPublished"`
	}

	var form PageForm
	if err := c.BodyParser(&form); err != nil {
		return c.Status(400).SendString("Invalid form data")
	}

	// Validate required fields
	if form.Title == "" || form.Slug == "" || form.Content == "" {
		return c.Render("admin/static-pages/add", fiber.Map{
			"title": "Add Static Page",
			"page":  form,
			"error": "Title, slug, and content are required",
		})
	}

	// Check if slug already exists
	var existingPage models.StaticPage
	if err := database.DB.Where("slug = ?", form.Slug).First(&existingPage).Error; err == nil {
		return c.Render("admin/static-pages/add", fiber.Map{
			"title": "Add Static Page",
			"page":  form,
			"error": "A page with this slug already exists",
		})
	}

	// Create page
	page := models.StaticPage{
		Title:           form.Title,
		Slug:            form.Slug,
		Content:         form.Content,
		MetaTitle:       stringPtr(form.MetaTitle),
		MetaDescription: stringPtr(form.MetaDescription),
		MetaKeywords:    stringPtr(form.MetaKeywords),
		IsPublished:     form.IsPublished == "on" || form.IsPublished == "true",
	}

	if err := database.DB.Create(&page).Error; err != nil {
		utils.Error("Failed to create static page: %v", err)
		return c.Render("admin/static-pages/add", fiber.Map{
			"title": "Add Static Page",
			"page":  form,
			"error": "Failed to add page",
		})
	}

	return c.Redirect("/admin/static-pages?success=Page added successfully")
}

// GetEditStaticPage handles GET /admin/static-pages/edit/:id
func GetEditStaticPage(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var page models.StaticPage
	if err := database.DB.First(&page, id).Error; err != nil {
		return c.Status(404).Render("error", fiber.Map{
			"title":      "404 - Not Found",
			"message":    "Page not found",
			"statusCode": 404,
		})
	}

	return c.Render("admin/static-pages/edit", fiber.Map{
		"title": "Edit Static Page",
		"page":  page,
		"error": nil,
	})
}

// PostEditStaticPage handles POST /admin/static-pages/edit/:id
func PostEditStaticPage(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	type PageForm struct {
		Title           string `form:"title"`
		Slug            string `form:"slug"`
		Content         string `form:"content"`
		MetaTitle       string `form:"metaTitle"`
		MetaDescription string `form:"metaDescription"`
		MetaKeywords    string `form:"metaKeywords"`
		IsPublished     string `form:"isPublished"`
	}

	var form PageForm
	if err := c.BodyParser(&form); err != nil {
		return c.Status(400).SendString("Invalid form data")
	}

	// Validate required fields
	if form.Title == "" || form.Slug == "" || form.Content == "" {
		var page models.StaticPage
		database.DB.First(&page, id)
		return c.Render("admin/static-pages/edit", fiber.Map{
			"title": "Edit Static Page",
			"page":  page,
			"error": "Title, slug, and content are required",
		})
	}

	// Check if slug is taken by another page
	var existingPage models.StaticPage
	if err := database.DB.Where("slug = ? AND id != ?", form.Slug, id).First(&existingPage).Error; err == nil {
		var page models.StaticPage
		database.DB.First(&page, id)
		return c.Render("admin/static-pages/edit", fiber.Map{
			"title": "Edit Static Page",
			"page":  page,
			"error": "A page with this slug already exists",
		})
	}

	// Update page
	var page models.StaticPage
	if err := database.DB.First(&page, id).Error; err != nil {
		return c.Redirect("/admin/static-pages?error=Page not found")
	}

	page.Title = form.Title
	page.Slug = form.Slug
	page.Content = form.Content
	page.MetaTitle = stringPtr(form.MetaTitle)
	page.MetaDescription = stringPtr(form.MetaDescription)
	page.MetaKeywords = stringPtr(form.MetaKeywords)
	page.IsPublished = form.IsPublished == "on" || form.IsPublished == "true"

	if err := database.DB.Save(&page).Error; err != nil {
		utils.Error("Failed to update static page: %v", err)
		return c.Render("admin/static-pages/edit", fiber.Map{
			"title": "Edit Static Page",
			"page":  page,
			"error": "Failed to update page",
		})
	}

	return c.Redirect("/admin/static-pages?success=Page updated successfully")
}

// DeleteStaticPage handles POST /admin/static-pages/delete/:id
func DeleteStaticPage(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := database.DB.Delete(&models.StaticPage{}, id).Error; err != nil {
		utils.Error("Failed to delete static page: %v", err)
		return c.Redirect("/admin/static-pages?error=Failed to delete page")
	}

	return c.Redirect("/admin/static-pages?success=Page deleted successfully")
}
