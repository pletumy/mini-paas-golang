package api

import (
	"mini-paas/backend/internal/repository"
	"mini-paas/backend/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AppHandler struct {
	appSvc *services.AppService
	repo   repository.ApplicationRepository
}

func NewAppHandler(appSvc *services.AppService, repo repository.ApplicationRepository) *AppHandler {
	return &AppHandler{
		appSvc: appSvc,
		repo:   repo,
	}
}

func toItem(a *repository.Application) AppItem {
	return AppItem{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		GitURL:      a.GitURL,
		ImageURL:    a.ImageURL,
		DeployURL:   a.DeployURL,
		Status:      a.Status,
		CreatedAt:   a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   a.UpdatedAt.Format(time.RFC3339),
	}
}

// POST /api/apps
func (h *AppHandler) Create(c *gin.Context) {
	var req CreateAppRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.appSvc.Deploy(services.DeployRequest{
		Name:   req.Name,
		GitURL: req.GitURL,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	app, _ := h.repo.GetByID(id)
	c.JSON(http.StatusOK, CreateAppResponse{
		ID:     id,
		Status: app.Status,
	})
}

// GET /api/apps
func (h *AppHandler) List(c *gin.Context) {
	list, err := h.repo.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]AppItem, 0, len(list))
	for i := range list {
		a := list[i]
		out = append(out, toItem(&a))
	}
	c.JSON(http.StatusOK, out)
}

// GET /api/apps/:id
func (h *AppHandler) Get(c *gin.Context) {
	id := c.Param("id")
	a, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if a == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, toItem(a))
}

// DELETE /api/apps/:id
func (h *AppHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GET /api/apps/:id/logs (stub)
func (h *AppHandler) Logs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"lines": []string{
			"[stub] build started",
			"[stub] build done",
			"[stub] pushed image registry.local/app:latest",
			"[stub] deployed http://app.local",
		},
	})
}
