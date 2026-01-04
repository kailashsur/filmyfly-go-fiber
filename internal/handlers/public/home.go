package public

import (
	"filmyfly-go-fiber/internal/database"
	"filmyfly-go-fiber/internal/database/models"

	"github.com/gofiber/fiber/v2"
)

// GetHomePage handles GET /
func GetHomePage(c *fiber.Ctx) error {
	// Fetch trending movies
	var trendingMoviesData []models.TrendingMovie
	database.DB.
		Preload("Movie").
		Order("\"order\" ASC").
		Find(&trendingMoviesData)

	trendingMovies := make([]models.Movie, 0, len(trendingMoviesData))
	for _, tm := range trendingMoviesData {
		trendingMovies = append(trendingMovies, tm.Movie)
	}

	// Fetch recent movies
	var recentMovies []models.Movie
	database.DB.
		Select("id, title, slug, thumbnail, keywords, \"releaseYear\", genre").
		Order("\"createdAt\" DESC").
		Limit(50).
		Find(&recentMovies)

	// Fetch categories
	var categories []models.Category
	database.DB.Order("name ASC").Find(&categories)

	return c.Render("index", fiber.Map{
		"title":          "FilmyFly - Download Latest Movies",
		"trendingMovies": trendingMovies,
		"recentMovies":   recentMovies,
		"categories":     categories,
	})
}
