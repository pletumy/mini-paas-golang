package api

import (
	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"
	"mini-paas/backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppHandler struct {
	appService services.AppService
}

func NewAppHandler(service services.AppService) *AppHandler {
	return &AppHandler{appService: service}
}

// POST /api/apps
func (h *AppHandler) CreateNewApp(c *gin.Context) {
	var req CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app := &models.Application{
		Name:        req.Name,
		GitURL:      req.GitURL,
		Description: req.Description,
	}

	newApp, err := h.appService.CreateApp(c.Request.Context(), app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateAppResponse{
		ID:          newApp.ID.String(),
		Name:        newApp.Name,
		Description: newApp.Description,
		Status:      newApp.Status,
	})
}

// GET /api/apps
func (h *AppHandler) ListAllApps(c *gin.Context) {
	var req ListAppsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var status string
	if req.Status != nil {
		status = *req.Status
	}
	var search string
	if req.Search != nil {
		search = *req.Search
	}

	f := repository.AppFilter{
		Status: status,
		Search: search,
	}

	page := repository.Page{Limit: req.Limit, Offset: req.Offset}
	sort := repository.Sort{Field: req.SortBy, Desc: req.Desc}

	apps, err := h.appService.ListApps(c.Request.Context(), f, page, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]CreateAppResponse, 0, len(apps.Items))
	for _, a := range apps.Items {
		resp = append(resp, CreateAppResponse{
			ID:          a.ID.String(),
			Name:        a.Name,
			Status:      a.Status,
			Description: a.Description,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  resp,
		"total":  apps.Total,
		"limit":  page.Limit,
		"offset": page.Offset,
	})
}

// GET /api/app/:id
func (h *AppHandler) GetApplicatonByID(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	app, err := h.appService.GetAppByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	c.JSON(http.StatusOK, CreateAppResponse{
		ID:          app.ID.String(),
		Name:        app.Name,
		Status:      app.Status,
		Description: app.Description,
	})
}

// DELETE /api/app/:id
func (h *AppHandler) DeleteApplication(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	err = h.appService.DeleteApp(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application deleted"})
}
