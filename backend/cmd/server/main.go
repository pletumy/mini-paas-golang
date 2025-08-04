package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"mini-paas-golang/backend/configs"
	"mini-paas-golang/backend/internal/handlers"
	"mini-paas-golang/backend/internal/repository"
	"mini-paas-golang/backend/internal/services"
	"mini-paas-golang/backend/pkg/database"
	"mini-paas-golang/backend/pkg/docker"
	"mini-paas-golang/backend/pkg/kubernetes"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logging
	setupLogging()

	// Initialize database
	db, err := database.NewConnection(&config.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	appRepo := repository.NewApplicationRepository(db.DB)

	// Initialize services
	dockerSvc := docker.NewService(&config.Docker)
	k8sSvc := kubernetes.NewService(&config.K8s)
	appService := services.NewApplicationService(appRepo, dockerSvc, k8sSvc)

	// Initialize handlers
	appHandler := handlers.NewApplicationHandler(appService)

	// Setup Gin router
	router := setupRouter(appHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         config.Server.GetServerAddr(),
		Handler:      router,
		ReadTimeout:  time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(config.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logrus.Infof("Starting server on %s", config.Server.GetServerAddr())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited")
}

// setupLogging configures the logging
func setupLogging() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	if configs.IsDevelopment() {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

// setupRouter configures the Gin router with middleware and routes
func setupRouter(appHandler *handlers.ApplicationHandler) *gin.Engine {
	if !configs.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Mini PaaS API is running",
			"time":    time.Now().UTC(),
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Application routes
		applications := api.Group("/applications")
		{
			applications.POST("/", appHandler.CreateApplication)
			applications.GET("/", appHandler.GetUserApplications)
			applications.GET("/:id", appHandler.GetApplication)
			applications.PUT("/:id", appHandler.UpdateApplication)
			applications.DELETE("/:id", appHandler.DeleteApplication)
			applications.POST("/:id/deploy", appHandler.DeployApplication)
			applications.GET("/:id/logs", appHandler.GetApplicationLogs)
			applications.GET("/:id/metrics", appHandler.GetApplicationMetrics)
		}

		// TODO: Add authentication routes
		// TODO: Add user management routes
		// TODO: Add admin routes
	}

	// Swagger documentation
	// TODO: Add Swagger setup

	return router
} 