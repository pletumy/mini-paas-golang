package api

import (
	"mini-paas/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(
	r *gin.Engine,
	appService services.AppService,
	deployService services.DeploymentService,
	userService services.UserService,
	logService services.LogService,
) {
	api := r.Group("/api")

	// app routes
	appHandler := NewAppHandler(appService)
	api.POST("/apps", appHandler.CreateNewApp)
	api.GET("/apps", appHandler.ListAllApps)
	api.GET("/apps/app/:id", appHandler.GetApplicatonByID)
	api.DELETE("/apps/app/:id", appHandler.DeleteApplication)

	// deployment
	depHandler := NewDeploymentHandler(deployService)
	api.POST("/deployments", depHandler.CreateDeploymentHandler)
	api.GET("/deployments", depHandler.ListAllDeploymentsHandler)
	api.GET("/deployments/:id", depHandler.GetDeploymentByIDHandler)
	api.POST("/deployments/deploy", depHandler.DeployAppHandler)
	api.GET("/deployments/:id/status", depHandler.GetDeploymentStatusHandler)

	// user
	userHandler := NewUserHandler(userService)
	api.POST("/users", userHandler.CreateUserHandler)
	api.GET("/users/user/:id", userHandler.GetUserByIDHandler)
	api.GET("/users/user/email", userHandler.GetUserByEmailHandler)

	// log
	logHandler := NewLogHandler(logService)
	api.POST("/logs", logHandler.CreateLogHandler)
	api.GET("/logs", logHandler.ListAllLogsHandler)
	api.GET("/deployments/:id", logHandler.StreamLogsHandler)
}
