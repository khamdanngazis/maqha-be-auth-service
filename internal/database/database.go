// internal/database/database.go

package database

import (
	"fmt"
	"maqhaa/auth_service/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB creates a new database connection based on the provided configuration.
func NewDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port,
	)
	gormConfig := &gorm.Config{}
	if cfg.Debug {
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // Set logger level to Info
		}
	}
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return db, nil
}
