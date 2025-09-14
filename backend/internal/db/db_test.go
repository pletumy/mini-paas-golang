package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectTestDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sqlDB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	log.Println("database connected successfully")
	return db, nil
}

func TruncateAll(db *gorm.DB) error {
	tables := []string{"applications", "users", "deployments", "logs"}
	for _, t := range tables {
		if err := db.Exec("TRUNCATE TABLE " + t + " RESTART IDENTITY CASCADE").Error; err != nil {
			return err
		}
	}
	return nil
}
