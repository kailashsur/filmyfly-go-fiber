package admin

import (
	"encoding/json"
	"strconv"

	"filmyfly-go-fiber/internal/database"
	"filmyfly-go-fiber/internal/database/models"

	"github.com/gofiber/fiber/v2"
)

// GetMovieList displays paginated list of movies
func GetMovieList(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit := 10
	offset := (page - 1) * limit

	search := c.Query("search", "")
	success := c.Query("success", "")
	errorMsg := c.Query("error", "")

	// Get trending movies with movie details
	var trendingMoviesData []models.TrendingMovie
	database.DB.Preload("Movie").Order("\"order\" ASC").Find(&trendingMoviesData)

	trendingMovies := make([]models.Movie, 0)
	trendingMovieIDs := make(map[int]bool)
	for _, tm := range trendingMoviesData {
		trendingMovies = append(trendingMovies, tm.Movie)
		trendingMovieIDs[tm.Movie.ID] = true
	}

	// Get all movies with pagination
	var movies []models.Movie
	var total int64

	query := database.DB.Model(&models.Movie{})

	if search != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+search+"%")
	}

	query.Count(&total)
	query.Order("\"createdAt\" DESC").
		Limit(limit).
		Offset(offset).
		Find(&movies)

	// Add isTrending flag to each movie
	type MovieWithTrending struct {
		models.Movie
		IsTrending bool
	}

	moviesWithTrending := make([]MovieWithTrending, len(movies))
	for i, movie := range movies {
		moviesWithTrending[i] = MovieWithTrending{
			Movie:      movie,
			IsTrending: trendingMovieIDs[movie.ID],
		}
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.Render("admin/movies/list", fiber.Map{
		"title":          "Movies Management",
		"movies":         moviesWithTrending,
		"trendingMovies": trendingMovies,
		"currentPage":    page,
		"totalPages":     totalPages,
		"search":         search,
		"success":        success,
		"error":          errorMsg,
		"user": map[string]interface{}{
			"email": "admin@filmyfly.work",
		},
	})
}

// GetAddMovie renders add movie form
func GetAddMovie(c *fiber.Ctx) error {
	var categories []models.Category
	database.DB.Order("name ASC").Find(&categories)

	return c.Render("admin/movies/add", fiber.Map{
		"title":      "Add Movie",
		"categories": categories,
		"user": map[string]interface{}{
			"email": "admin@filmyfly.work",
		},
	})
}

// PostAddMovie handles movie creation
func PostAddMovie(c *fiber.Ctx) error {
	movie := new(models.Movie)

	if err := c.BodyParser(movie); err != nil {
		return c.Redirect("/admin/movies/add?error=Invalid data")
	}

	if movie.Title == "" || movie.Slug == "" {
		return c.Redirect("/admin/movies/add?error=Title and slug are required")
	}

	if err := database.DB.Create(movie).Error; err != nil {
		return c.Redirect("/admin/movies/add?error=Failed to create movie")
	}

	return c.Redirect("/admin/movies?success=Movie added successfully")
}

// GetBulkAddMovies renders bulk import page
func GetBulkAddMovies(c *fiber.Ctx) error {
	return c.Render("admin/movies/bulk-add", fiber.Map{
		"title": "Bulk Add Movies",
		"user": map[string]interface{}{
			"email": "admin@filmyfly.work",
		},
	})
}

// PostBulkAddMovies handles bulk movie import from JSON
func PostBulkAddMovies(c *fiber.Ctx) error {
	type BulkImportRequest struct {
		MoviesJSON string `form:"moviesJson"`
	}

	var req BulkImportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Redirect("/admin/movies/bulk-add?error=Invalid request")
	}

	if req.MoviesJSON == "" {
		return c.Redirect("/admin/movies/bulk-add?error=JSON data is required")
	}

	// Parse JSON
	var movies []models.Movie
	if err := json.Unmarshal([]byte(req.MoviesJSON), &movies); err != nil {
		return c.Redirect("/admin/movies/bulk-add?error=Invalid JSON format: " + err.Error())
	}

	if len(movies) == 0 {
		return c.Redirect("/admin/movies/bulk-add?error=No movies found in JSON")
	}

	// Import movies
	successCount := 0
	failCount := 0

	for _, movie := range movies {
		if movie.Title == "" || movie.Slug == "" {
			failCount++
			continue
		}

		if err := database.DB.Create(&movie).Error; err != nil {
			failCount++
			continue
		}
		successCount++
	}

	return c.Redirect("/admin/movies?success=Imported " + strconv.Itoa(successCount) + " movies (" + strconv.Itoa(failCount) + " failed)")
}

// GetEditMovie renders edit movie form
func GetEditMovie(c *fiber.Ctx) error {
	id := c.Params("id")

	var movie models.Movie
	if err := database.DB.First(&movie, id).Error; err != nil {
		return c.Redirect("/admin/movies?error=Movie not found")
	}

	var categories []models.Category
	database.DB.Order("name ASC").Find(&categories)

	return c.Render("admin/movies/edit", fiber.Map{
		"title":      "Edit Movie",
		"movie":      movie,
		"categories": categories,
		"user": map[string]interface{}{
			"email": "admin@filmyfly.work",
		},
	})
}

// PostEditMovie handles movie update
func PostEditMovie(c *fiber.Ctx) error {
	id := c.Params("id")

	var movie models.Movie
	if err := database.DB.First(&movie, id).Error; err != nil {
		return c.Redirect("/admin/movies?error=Movie not found")
	}

	if err := c.BodyParser(&movie); err != nil {
		return c.Redirect("/admin/movies/edit/" + id + "?error=Invalid data")
	}

	if err := database.DB.Save(&movie).Error; err != nil {
		return c.Redirect("/admin/movies/edit/" + id + "?error=Failed to update movie")
	}

	return c.Redirect("/admin/movies?success=Movie updated successfully")
}

// DeleteMovie handles movie deletion
func DeleteMovie(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := database.DB.Delete(&models.Movie{}, id).Error; err != nil {
		return c.Redirect("/admin/movies?error=Failed to delete movie")
	}

	return c.Redirect("/admin/movies?success=Movie deleted successfully")
}

// AddToTrending adds movie to trending
func AddToTrending(c *fiber.Ctx) error {
	id := c.Params("id")
	movieID, _ := strconv.Atoi(id)

	// Get max order
	var maxOrder int
	database.DB.Model(&models.TrendingMovie{}).Select("COALESCE(MAX(\"order\"), 0)").Scan(&maxOrder)

	trending := models.TrendingMovie{
		MovieID: movieID,
		Order:   maxOrder + 1,
	}

	if err := database.DB.Create(&trending).Error; err != nil {
		return c.Redirect("/admin/movies?error=Failed to add to trending")
	}

	return c.Redirect("/admin/movies?success=Added to trending")
}

// RemoveFromTrending removes movie from trending
func RemoveFromTrending(c *fiber.Ctx) error {
	id := c.Params("id")
	movieID, _ := strconv.Atoi(id)

	if err := database.DB.Where("\"movieId\" = ?", movieID).Delete(&models.TrendingMovie{}).Error; err != nil {
		return c.Redirect("/admin/movies?error=Failed to remove from trending")
	}

	return c.Redirect("/admin/movies?success=Removed from trending")
}
