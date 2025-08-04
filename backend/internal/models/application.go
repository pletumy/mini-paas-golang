package models

import (
	"time"

	"github.com/google/uuid"
)

// Application represents a deployed application
type Application struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=1,max=100"`
	Description string    `json:"description" db:"description"`
	GitURL      string    `json:"git_url" db:"git_url" validate:"required,url"`
	Branch      string    `json:"branch" db:"branch" validate:"required"`
	Port        int       `json:"port" db:"port" validate:"required,min=1,max=65535"`
	Environment string    `json:"environment" db:"environment" validate:"required,oneof=development staging production"`
	Status      string    `json:"status" db:"status" validate:"required,oneof=pending building deploying running failed stopped"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	DeploymentURL string  `json:"deployment_url" db:"deployment_url"`
	Replicas    int       `json:"replicas" db:"replicas" validate:"required,min=1,max=10"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ApplicationCreate represents the data needed to create a new application
type ApplicationCreate struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description"`
	GitURL      string `json:"git_url" validate:"required,url"`
	Branch      string `json:"branch" validate:"required"`
	Port        int    `json:"port" validate:"required,min=1,max=65535"`
	Environment string `json:"environment" validate:"required,oneof=development staging production"`
	Replicas    int    `json:"replicas" validate:"required,min=1,max=10"`
}

// ApplicationUpdate represents the data needed to update an application
type ApplicationUpdate struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description"`
	GitURL      *string `json:"git_url" validate:"omitempty,url"`
	Branch      *string `json:"branch"`
	Port        *int    `json:"port" validate:"omitempty,min=1,max=65535"`
	Environment *string `json:"environment" validate:"omitempty,oneof=development staging production"`
	Replicas    *int    `json:"replicas" validate:"omitempty,min=1,max=10"`
}

// ApplicationList represents a paginated list of applications
type ApplicationList struct {
	Applications []Application `json:"applications"`
	Total        int64         `json:"total"`
	Page         int           `json:"page"`
	Limit        int           `json:"limit"`
}

// ApplicationLog represents application logs
type ApplicationLog struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ApplicationID uuid.UUID `json:"application_id" db:"application_id"`
	Level         string    `json:"level" db:"level"`
	Message       string    `json:"message" db:"message"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
}

// ApplicationMetrics represents application metrics
type ApplicationMetrics struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ApplicationID uuid.UUID `json:"application_id" db:"application_id"`
	CPUUsage      float64   `json:"cpu_usage" db:"cpu_usage"`
	MemoryUsage   float64   `json:"memory_usage" db:"memory_usage"`
	RequestCount  int64     `json:"request_count" db:"request_count"`
	ErrorCount    int64     `json:"error_count" db:"error_count"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
} 