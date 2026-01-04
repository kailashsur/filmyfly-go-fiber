package admin

import (
	"strconv"
	"strings"

	"filmyfly-go-fiber/internal/database"
	"filmyfly-go-fiber/internal/database/models"
	"filmyfly-go-fiber/internal/utils"

	"github.com/gofiber/fiber/v2"
)

const ItemsPerPage = 10

// GetMovieList handles GET /admin/movies
func GetMovieList(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	searchQuery := strings.TrimSpace(c.Query("q", ""))
	offset := utils.GetOffset(page, ItemsPerPage)

	// Build where clause for search
	var whereClause string
	var args []interface{}
	if searchQuery != "" {
		whereClause = "LOWER(title) LIKE ? OR LOWER(genre) LIKE ? OR LOWER(cast) LIKE ? OR LOWER(keywords) LIKE ? OR LOWER(slug) LIKE ?"
		searchPattern := "%" + strings.ToLower(searchQuery) + "%"
		args = []interface{}{searchPattern, searchPattern, searchPattern, searchPattern, searchPattern}
	}

	// Get total count
	var totalMovies int64
	query := database.DB.Model(&models.Movie{})
	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}
	query.Count(&totalMovies)

	// Get trending movie IDs
	var trendingMoviesData []models.TrendingMovie
	database.DB.Preload("Movie").Order("\"order\" ASC").Find(&trendingMoviesData)

	trendingMovieIds := make(map[int]bool)
	for _, tm := range trendingMoviesData {
		trendingMovieIds[tm.MovieID] = true
	}

	// Get movies with pagination
	var movies []models.Movie
	query = database.DB.Order("\"createdAt\" DESC").Offset(offset).Limit(ItemsPerPage)
	if whereClause != "" {
		query = query.Where(whereClause, args...)
	}
	query.Find(&movies)

	// Mark movies as trending
	type MovieWithTrending struct {
		models.Movie
		IsTrending bool
	}
	moviesWithTrending := make([]MovieWithTrending, len(movies))
	for i, movie := range movies {
		moviesWithTrending[i] = MovieWithTrending{
			Movie:      movie,
			IsTrending: trendingMovieIds[movie.ID],
		}
	}

	totalPages := int((totalMovies + int64(ItemsPerPage) - 1) / int64(ItemsPerPage))

	return c.Render("admin/movies/list", fiber.Map{
		"title":          "Manage Movies",
		"movies":         moviesWithTrending,
		"currentPage":    page,
		"totalPages":     totalPages,
		"totalMovies":    totalMovies,
		"hasNextPage":    page < totalPages,
		"hasPrevPage":    page > 1,
		"searchQuery":    searchQuery,
		"trendingMovies": trendingMoviesData,
	})
}

// GetAddMovie handles GET /admin/movies/add
func GetAddMovie(c *fiber.Ctx) error {
	var categories []models.Category
	database.DB.Order("name ASC").Find(&categories)

	return c.Render("admin/movies/add", fiber.Map{
		"title":      "Add Movie",
		"movie":      nil,
		"categories": categories,
		"error":      nil,
	})
}

// PostAddMovie handles POST /admin/movies/add
func PostAddMovie(c *fiber.Ctx) error {
	type MovieForm struct {
		Title       string `form:"title"`
		Slug        string `form:"slug"`
		Description string `form:"description"`
		Thumbnail   string `form:"thumbnail"`
		Genre       string `form:"genre"`
		Languages   string `form:"languages"`
		Duration    string `form:"duration"`
		ReleaseYear string `form:"releaseYear"`
		Cast        string `form:"cast"`
		Sizes       string `form:"sizes"`
		DownloadURL string `form:"downloadUrl"`
		Screenshot  string `form:"screenshot"`
		Keywords    string `form:"keywords"`
		CategoryID  string `form:"categoryId"`
	}

	var form MovieForm
	if err := c.BodyParser(&form); err != nil {
		return c.Status(400).SendString("Invalid form data")
	}

	// Validate required fields
	if form.Title == "" || form.Slug == "" {
		var categories []models.Category
		database.DB.Order("name ASC").Find(&categories)
		return c.Render("admin/movies/add", fiber.Map{
			"title":      "Add Movie",
			"movie":      form,
			"categories": categories,
			"error":      "Title and slug are required",
		})
	}

	// Check if slug already exists
	var existingMovie models.Movie
	if err := database.DB.Where("slug = ?", form.Slug).First(&existingMovie).Error; err == nil {
		var categories []models.Category
		database.DB.Order("name ASC").Find(&categories)
		return c.Render("admin/movies/add", fiber.Map{
			"title":      "Add Movie",
			"movie":      form,
			"categories": categories,
			"error":      "A movie with this slug already exists",
		})
	}

	// Create movie
	movie := models.Movie{
		Title:       form.Title,
		Slug:        form.Slug,
		Description: stringPtr(form.Description),
		Thumbnail:   stringPtr(form.Thumbnail),
		Genre:       stringPtr(form.Genre),
		Languages:   stringPtr(form.Languages),
		Duration:    stringPtr(form.Duration),
		Cast:        stringPtr(form.Cast),
		Sizes:       stringPtr(form.Sizes),
		DownloadURL: stringPtr(form.DownloadURL),
		Screenshot:  stringPtr(form.Screenshot),
		Keywords:    stringPtr(form.Keywords),
	}

	if form.ReleaseYear != "" {
		year, _ := strconv.Atoi(form.ReleaseYear)
		movie.ReleaseYear = &year
	}

	if form.CategoryID != "" {
		catID, _ := strconv.Atoi(form.CategoryID)
		movie.CategoryID = &catID
	}

	if err := database.DB.Create(&movie).Error; err != nil {
		utils.Error("Failed to create movie: %v", err)
		var categories []models.Category
		database.DB.Order("name ASC").Find(&categories)
		return c.Render("admin/movies/add", fiber.Map{
			"title":      "Add Movie",
			"movie":      form,
			"categories": categories,
			"error":      "Failed to add movie",
		})
	}

	return c.Redirect("/admin/movies?success=Movie added successfully")
}

// DeleteMovie handles POST /admin/movies/delete/:id
func DeleteMovie(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := database.DB.Delete(&models.Movie{}, id).Error; err != nil {
		utils.Error("Failed to delete movie: %v", err)
		return c.Redirect("/admin/movies?error=Failed to delete movie")
	}

	return c.Redirect("/admin/movies?success=Movie deleted successfully")
}

// AddToTrending handles POST /admin/movies/trending/add/:id
func AddToTrending(c *fiber.Ctx) error {
	movieID, _ := strconv.Atoi(c.Params("id"))

	// Check if movie exists
	var movie models.Movie
	if err := database.DB.First(&movie, movieID).Error; err != nil {
		return c.Redirect("/admin/movies?error=Movie not found")
	}

	// Check if already in trending
	var existing models.TrendingMovie
	if err := database.DB.Where("\"movieId\" = ?", movieID).First(&existing).Error; err == nil {
		return c.Redirect("/admin/movies?error=Movie is already in trending")
	}

	// Get current max order
	var maxOrder int
	database.DB.Model(&models.TrendingMovie{}).Select("COALESCE(MAX(\"order\"), -1)").Scan(&maxOrder)

	// Add to trending
	trending := models.TrendingMovie{
		MovieID: movieID,
		Order:   maxOrder + 1,
	}

	if err := database.DB.Create(&trending).Error; err != nil {
		utils.Error("Failed to add to trending: %v", err)
		return c.Redirect("/admin/movies?error=Failed to add movie to trending")
	}

	return c.Redirect("/admin/movies?success=Movie added to trending")
}

// RemoveFromTrending handles POST /admin/movies/trending/remove/:id
func RemoveFromTrending(c *fiber.Ctx) error {
	movieID, _ := strconv.Atoi(c.Params("id"))

	if err := database.DB.Where("\"movieId\" = ?", movieID).Delete(&models.TrendingMovie{}).Error; err != nil {
		utils.Error("Failed to remove from trending: %v", err)
		return c.Redirect("/admin/movies?error=Failed to remove movie from trending")
	}

	return c.Redirect("/admin/movies?success=Movie removed from trending")
}

// Helper function
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
