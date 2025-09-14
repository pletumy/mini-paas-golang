package db

import (
	"fmt"

	"mini-paas/backend/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202309040001_create_users",
			Migrate: func(d *gorm.DB) error {
				return d.AutoMigrate(&models.User{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable("users")
			},
		},
		{
			ID: "202309040002_create_apps",
			Migrate: func(d *gorm.DB) error {
				return d.AutoMigrate(&models.Application{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable("applications")
			},
		},
		{
			ID: "202309040003_create_deployments",
			Migrate: func(d *gorm.DB) error {
				return d.AutoMigrate(&models.Deployment{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable("deployments")
			},
		},
		{
			ID: "202309040004_create_logs",
			Migrate: func(d *gorm.DB) error {
				return d.AutoMigrate(&models.Log{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable("logs")
			},
		},
	})

	if err := m.Migrate(); err != nil {
		return err
	}

	fmt.Println("migration ran successfully")
	return nil
}
