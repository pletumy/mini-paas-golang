package api

import (
	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"
	"mini-paas/backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeploymentHandler struct {
	deploymentService services.DeploymentService
}

func NewDeploymentHandler(s services.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{deploymentService: s}
}

// POST /api/deployments
func (h *DeploymentHandler) CreateDeploymentHandler(c *gin.Context) {
	var req CreateDeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appUUID, err := uuid.Parse(req.AppID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid AppID"})
		return
	}
	dep := &models.Deployment{
		AppID:   appUUID,
		Version: req.Version,
	}

	newDep, err := h.deploymentService.CreateDeployment(c, dep)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, DeploymentResponse{
		ID:      newDep.ID.String(),
		AppID:   newDep.AppID.String(),
		Version: newDep.Version,
		Status:  newDep.Status,
	})
}

// GET api/deployments
func (h *DeploymentHandler) ListAllDeploymentsHandler(c *gin.Context) {
	var req ListDeploymentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var appUUID *uuid.UUID
	if req.AppID != nil {
		parsedUUID, err := uuid.Parse(*req.AppID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid AppID"})
			return
		}
		appUUID = &parsedUUID
	}

	f := repository.DeploymentFilter{
		AppID:  appUUID,
		Status: req.Status,
	}

	page := repository.Page{Limit: req.Limit, Offset: req.Offset}
	sort := repository.Sort{Field: req.SortBy, Desc: req.Desc}

	deps, err := h.deploymentService.ListAllDeployments(c.Request.Context(), f, page, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]DeploymentResponse, 0, len(deps.Items))
	for _, a := range deps.Items {
		resp = append(resp, DeploymentResponse{
			ID:      a.ID.String(),
			AppID:   a.AppID.String(),
			Version: a.Version,
			Status:  a.Status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  resp,
		"total":  deps.Total,
		"limit":  page.Limit,
		"offset": page.Offset,
	})
}

// GET /api/deployments/:id
func (h *DeploymentHandler) GetDeploymentByIDHandler(c *gin.Context) {
	idStr := c.Param("id")
	uid, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	dep, err := h.deploymentService.GetDeploymentByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "deployment not found"})
		return
	}

	c.JSON(http.StatusOK, DeploymentResponse{
		ID:      dep.ID.String(),
		AppID:   dep.AppID.String(),
		Version: dep.Version,
		Status:  dep.Status,
	})
}

// // DELETE /api/deployments/:id
// func (h *DeploymentHandler) DeleteApplication() {

//
