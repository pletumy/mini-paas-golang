package main

import (
	"log"
	"mini-paas/backend/internal/db"
	"mini-paas/backend/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=123 dbname=mini_paas port=5432 sslmode=disable"

	// connect db
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// auto migrate
	m := gormigrate.New(dbConn, gormigrate.DefaultOptions, db.Migrations())
	if err := m.Migrate(); err != nil {
		log.Fatalf("cound not migrate: %v", err)
	}
	log.Printf("Migrating tables: %T %T %T %T",
		models.Application{},
		models.Deployment{},
		models.Log{},
		models.User{},
	)

	log.Printf("Migration ran successfully :3")

	// // create repo & services
	// appRepo := repository.NewDBRepository(dbConn)
	// appService := services.NewMockAppService(appRepo)

	// // setUp Gin Router
	// r := gin.Default()
	// api.SetupRoutes(r, appService)

	// // run server
	// if err := r.Run(":8080"); err != nil {
	// 	log.Fatal("Failed to start server: ", err)
	// }
}
