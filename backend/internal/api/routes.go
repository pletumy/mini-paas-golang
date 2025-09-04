package api

import (
	"mini-paas/backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, appSvc *services.AppService) {
	api := r.Group("/api")

	// Deploy new app
	api.POST("/apps", func(c *gin.Context) {
		var req struct {
			Name   string `json:"name" binding:"required"`
			GitURL string `json:"git_url" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := appSvc.Deploy(services.DeployRequest{
			Name:   req.Name,
			GitURL: req.GitURL,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})

	})

	// Get app detail
	api.GET("/apps", func(c *gin.Context) {
		apps, err := appSvc.GetAllApps()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, apps)
	})

	api.GET("/apps/:id", func(c *gin.Context) {
		id := c.Param("id")
		app, err := appSvc.GetAppByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "app not found"})
			return
		}
		c.JSON(http.StatusOK, app)
	})

	api.DELETE("/apps/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := appSvc.DeleteApp(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	})
}
