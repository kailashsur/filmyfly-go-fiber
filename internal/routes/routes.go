package routes

import (
	"filmyfly-go-fiber/internal/handlers/admin"
	"filmyfly-go-fiber/internal/handlers/api"
	"filmyfly-go-fiber/internal/handlers/public"
	"filmyfly-go-fiber/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// Setup configures all application routes
func Setup(app *fiber.App) {
	// API Routes (for Astro frontend)
	apiGroup := app.Group("/api")
	{
		// Homepage data
		apiGroup.Get("/home", api.GetHomePageData)

		// Movies
		apiGroup.Get("/movies", api.GetMovies)
		apiGroup.Get("/movies/trending", api.GetTrendingMovies)
		apiGroup.Get("/movies/:slug", api.GetMovieBySlug)

		// Categories
		apiGroup.Get("/categories", api.GetCategories)
		apiGroup.Get("/categories/:slug", api.GetCategoryBySlug)

		// Search
		apiGroup.Get("/search", api.SearchMovies)

		// Static pages
		apiGroup.Get("/static-pages/:slug", api.GetStaticPageBySlug)

		// Astro Settings
		apiGroup.Get("/astro-settings", api.GetAstroSettings)
	}

	// Admin Routes
	adminGroup := app.Group("/admin")
	{
		// Public admin routes (login/logout)
		adminGroup.Get("/login", middleware.RedirectIfAuthenticated, admin.GetAdminLogin)
		adminGroup.Post("/login", middleware.RedirectIfAuthenticated, admin.PostAdminLogin)
		adminGroup.Post("/logout", admin.PostAdminLogout)

		// Protected admin routes
		adminGroup.Get("/", middleware.VerifyAdminToken, admin.GetAdminDashboard)
		adminGroup.Get("/system-check", middleware.VerifyAdminToken, admin.GetSystemCheck)
		adminGroup.Get("/logs", middleware.VerifyAdminToken, admin.GetLogs)

		// Settings
		adminGroup.Get("/settings", middleware.VerifyAdminToken, admin.GetSettings)
		adminGroup.Post("/settings", middleware.VerifyAdminToken, admin.PostSettings)

		// Astro Settings
		adminGroup.Get("/astro-settings", middleware.VerifyAdminToken, admin.GetAstroSettings)
		adminGroup.Post("/astro-settings", middleware.VerifyAdminToken, admin.PostAstroSettings)

		// Movie Management
		adminGroup.Get("/movies", middleware.VerifyAdminToken, admin.GetMovieList)
		adminGroup.Get("/movies/add", middleware.VerifyAdminToken, admin.GetAddMovie)
		adminGroup.Post("/movies/add", middleware.VerifyAdminToken, admin.PostAddMovie)
		adminGroup.Post("/movies/delete/:id", middleware.VerifyAdminToken, admin.DeleteMovie)

		// Trending Movies
		adminGroup.Post("/movies/trending/add/:id", middleware.VerifyAdminToken, admin.AddToTrending)
		adminGroup.Post("/movies/trending/remove/:id", middleware.VerifyAdminToken, admin.RemoveFromTrending)

		// Static Pages Management
		adminGroup.Get("/static-pages", middleware.VerifyAdminToken, admin.GetStaticPageList)
		adminGroup.Get("/static-pages/add", middleware.VerifyAdminToken, admin.GetAddStaticPage)
		adminGroup.Post("/static-pages/add", middleware.VerifyAdminToken, admin.PostAddStaticPage)
		adminGroup.Get("/static-pages/edit/:id", middleware.VerifyAdminToken, admin.GetEditStaticPage)
		adminGroup.Post("/static-pages/edit/:id", middleware.VerifyAdminToken, admin.PostEditStaticPage)
		adminGroup.Post("/static-pages/delete/:id", middleware.VerifyAdminToken, admin.DeleteStaticPage)
	}

	// Public Routes
	app.Get("/", public.GetHomePage)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("error", fiber.Map{
			"title":      "404 - Page Not Found",
			"message":    "The page you are looking for does not exist.",
			"statusCode": 404,
		})
	})
}
