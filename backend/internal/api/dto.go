package api

import "time"

// ===== Application DTOs =====
type CreateAppRequest struct {
	Name        string `json:"name" binding:"required"`
	GitURL      string `json:"git_url" binding:"required,url"`
	Description string `json:"description"`
}

type CreateAppResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type AppItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desciption"`
	GitURL      string `json:"git_url"`
	ImageURL    string `json:"image_url"`
	DeployURL   string `json:"deploy_url"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ListAppsRequest struct {
	OwnerID *string `form:"owner_id"`
	Status  *string `form:"status"`
	Search  *string `form:"search"`
	Limit   int     `form:"limit"`
	Offset  int     `form:"offset"`
	SortBy  string  `form:"sort"`
	Desc    bool    `form:"desc"`
}

// ===== Deployment DTOs =====
type CreateDeploymentRequest struct {
	AppID   string `json:"app_id" binding:"required"`
	Version string `json:"version" binding:"required"`
}

type DeploymentResponse struct {
	ID      string `json:"id"`
	AppID   string `json:"app_id"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

type ListDeploymentsRequest struct {
	AppID  *string `form:"owner_id"`
	Status *string `form:"status"`
	Limit  int     `form:"limit"`
	Offset int     `form:"offset"`
	SortBy string  `form:"sort"`
	Desc   bool    `form:"desc"`
}

// ===== User DTOs =====
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ===== Log DTOs =====
type CreateLogRequest struct {
	DeploymentID string `json:"deployment_id" binding:"required"`
	Message      string `json:"message" binding:"required"`
}

type LogResponse struct {
	ID           string    `json:"id"`
	DeploymentID string    `json:"deployment_id"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
}

type ListLogsRequest struct {
}
