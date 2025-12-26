package database

import (
	"fmt"
	"os"

	"laundry-go/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	var dsn string

	// Check if DATABASE_URL is provided (for Supabase, Heroku, etc.)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		// Use DATABASE_URL directly
		dsn = databaseURL
	} else {
		// Build DSN from individual config values
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.Port,
			cfg.Database.SSLMode,
		)
	}

	var logLevel logger.LogLevel
	if cfg.Server.Env == "production" {
		logLevel = logger.Silent
	} else {
		logLevel = logger.Info
	}

	// Configure for Supabase pgbouncer (disable prepared statement cache)
	// Use PreferSimpleProtocol to avoid prepared statement issues with pgbouncer
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disable prepared statements for pgbouncer compatibility
	}), &gorm.Config{
		Logger:      logger.Default.LogMode(logLevel),
		PrepareStmt: false, // Disable prepared statements for pgbouncer
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

