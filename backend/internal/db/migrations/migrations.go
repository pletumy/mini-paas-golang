package db

import (
	"mini-paas/backend/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "20250828001_create_applications",
			Migrate: func(d *gorm.DB) error {
				return d.AutoMigrate(
					&models.Application{},
					&models.Deployment{},
					&models.Log{},
					&models.User{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable("applications", "logs", "deployments", "users")
			},
		},
	}
}
