package main

import (
	"log"
	"os"

	"mini-paas/backend/internal/api"
	"mini-paas/backend/internal/db"
	"mini-paas/backend/internal/repository"
	"mini-paas/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// dsn = "host=localhost user=postgres password=123 dbname=mini_paas port=5432 sslmode=disable"      //dev
		dsn = "host=localhost user=postgres password=123 dbname=mini_paas_test port=5432 sslmode=disable" // test
	}

	gormDB, err := db.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := db.RunMigrations(gormDB); err != nil {
			log.Fatalf("Migration failed: %w", err)
		}
		return
	}

	// repo layers
	appRepo := repository.NewAppRepository(gormDB)
	depRepo := repository.NewDeploymentRepository(gormDB)
	userRepo := repository.NewUserRepository(gormDB)
	logRepo := repository.NewLogRepository(gormDB)

	// service layers
	appService := services.NewAppService(appRepo)
	depService := services.NewDeploymentService(depRepo)
	userService := services.NewUserService(userRepo)
	logService := services.NewLogService(logRepo)

	// api router
	r := gin.Default()
	api.SetUpRoutes(r, appService, depService, userService, logService)

	// start server
	log.Println("server running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
