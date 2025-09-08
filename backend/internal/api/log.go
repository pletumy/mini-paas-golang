package api

import (
	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"
	"mini-paas/backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LogHandler struct {
	logService services.LogService
}

func NewLogHandler(s services.LogService) *LogHandler {
	return &LogHandler{logService: s}
}

// POST /api/logs
func (h *LogHandler) CreateLogHandler(c *gin.Context) {
	var req CreateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	depUUID, err := uuid.Parse(req.DeploymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	logEntry := &models.Log{
		DeploymentID: depUUID,
		Message:      req.Message,
	}
	newLog, err := h.logService.CreateLog(c, logEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, LogResponse{
		ID:           newLog.ID.String(),
		DeploymentID: newLog.DeploymentID.String(),
		Message:      newLog.Message,
	})

}

// GET /api/logs
func (h *LogHandler) ListAllLogsHandler(c *gin.Context) {
	var (
		filter repository.LogFilter
		limit  = 50 //default limit
	)

	depploymentID := c.Query("deployment_id")
	if depploymentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deployment_id is required"})
		return
	}

	parsedDepID, err := uuid.Parse(depploymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	filter.DeploymentID = parsedDepID

	//optional: limit
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 200 {
			limit = v
		}
	}

	logs, err := h.logService.ListAllLogs(c.Request.Context(), filter, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]LogResponse, 0, len(logs))
	for _, log := range logs {
		resp = append(resp, LogResponse{
			ID:           log.ID.String(),
			DeploymentID: log.DeploymentID.String(),
			Message:      log.Message,
			Timestamp:    log.Timestamp,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items": resp,
		"limit": limit,
	})
}

// // GET /api/logs/log/:id
// func ()  {

// }

// // DELETE /api/logs/log/:id
