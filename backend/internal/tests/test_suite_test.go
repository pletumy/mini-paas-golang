package tests

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"mini-paas/backend/internal/api"
	"mini-paas/backend/internal/db"
	"mini-paas/backend/internal/repository"
	"mini-paas/backend/internal/services"

	"github.com/gin-gonic/gin"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	dsn := "host=localhost user=postgres password=123 dbname=mini_paas_test port=5432 sslmode=disable"
	database, err := db.ConnectDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect testDB: %v", err)
	}

	// run migrations
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("failed to migrate DB: %v", err)
	}

	// init repo
	appRepo := repository.NewAppRepository(database)
	depRepo := repository.NewDeploymentRepository(database)
	userRepo := repository.NewUserRepository(database)
	logRepo := repository.NewLogRepository(database)

	// init services
	appSvc := services.NewAppService(appRepo)
	depSvc := services.NewDeploymentService(depRepo)
	userSvc := services.NewUserService(userRepo)
	logSvc := services.NewLogService(logRepo)

	// set up gin + routes
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	api.SetUpRoutes(r, appSvc, depSvc, userSvc, logSvc)

	// start server
	testServer = httptest.NewServer(r)

	//	run test
	code := m.Run()
	testServer.Close()

	os.Exit(code)
}
