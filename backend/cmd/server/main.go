package main

import (
	"log"
	"mini-paas/backend/internal/api"
	"mini-paas/backend/internal/repository"
	"mini-paas/backend/internal/services"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	git := &services.FakeGitService{}
	build := &services.FakeBuildService{}
	registry := &services.FakeRegistryService{}
	deploy := &services.FakeDeployService{}
	repo := repository.NewMemoryRepo()

	appService := services.NewAppService(git, build, registry, deploy, repo)
	appHandler := api.NewAppHandler(appService, repo)

	// ====== HTTP server ======
	r := gin.Default()
	r.Use(cors.Default())

	apiGrp := r.Group("/api")
	{
		apiGrp.POST("/apps", appHandler.Create)
		apiGrp.GET("/apps", appHandler.List)
		apiGrp.GET("/apps/:id", appHandler.Get)
		apiGrp.DELETE("/apps/:id", appHandler.Delete)
		apiGrp.GET("/apps/:id/logs", appHandler.Logs)

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
