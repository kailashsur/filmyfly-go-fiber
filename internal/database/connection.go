package database

import (
	"fmt"
	"log"
	"time"

	"filmyfly-go-fiber/internal/config"
	"filmyfly-go-fiber/internal/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *config.Config) error {
	var err error

	// Configure GORM logger
	gormLogger := logger.Default
	if cfg.Environment == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else if cfg.PrismaLogQueries {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Open database connection
	DB, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connection established")

	return nil
}

// AutoMigrate runs database migrations for all models
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Movie{},
		&models.TrendingMovie{},
		&models.Setting{},
		&models.StaticPage{},
		&models.AstroSetting{},
	)
}
