package resources

import (
	"fmt"
	"summarizer/backend/http/config"
	"summarizer/backend/http/models"

	"github.com/jinzhu/gorm"
)

func NewDBConnection(cfg *config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.AutoMigrate(&models.User{}, &models.Chat{}, &models.ChatEntry{})

	return db, nil
}
