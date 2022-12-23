package db

import (
	"hasher/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db, nil
	// Migrate the schema
	// db.AutoMigrate(&Product{})
}
