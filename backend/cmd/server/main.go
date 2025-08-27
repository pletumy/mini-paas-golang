package main

import (
	"log"
	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"
	"mini-paas/backend/internal/routes"
	"mini-paas/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postres password=postgres db_name=mini_paas port=5432 sslmode=disable"

	// connect db
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	// auto migrate
	if err := db.AutoMigrate(&models.Application{}); err != nil {
		log.Fatal("failed to migrate: ", err)
	}

	// create repo & services
	appRepo := repository.NewDBRepository(db)
	appService := services.NewMockAppService(appRepo)

	// setUp Gin Router
	r := gin.Default()
	api := r.Group("/api")
	routes.SetupRoutes(api, appService)

	// run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
