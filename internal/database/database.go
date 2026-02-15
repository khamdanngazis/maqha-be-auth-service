// internal/database/database.go

package database

import (
	"fmt"
	"maqhaa/auth_service/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB creates a new database connection based on the provided configuration.
func NewDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	gormConfig := &gorm.Config{}
	if cfg.Debug {
		gormConfig = &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // Set logger level to Info
		}
	}
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return db, nil
}
