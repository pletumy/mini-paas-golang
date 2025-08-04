package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"mini-paas-golang/backend/internal/models"
)

// ApplicationRepository handles database operations for applications
type ApplicationRepository struct {
	db *sqlx.DB
}

// NewApplicationRepository creates a new application repository
func NewApplicationRepository(db *sqlx.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

// Create creates a new application
func (r *ApplicationRepository) Create(ctx context.Context, app *models.Application) error {
	query := `
		INSERT INTO applications (
			id, name, description, git_url, branch, port, environment, 
			status, image_url, deployment_url, replicas, user_id, created_at, updated_at
		) VALUES (
			:id, :name, :description, :git_url, :branch, :port, :environment,
			:status, :image_url, :deployment_url, :replicas, :user_id, :created_at, :updated_at
		)
	`

	app.ID = uuid.New()
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	app.Status = "pending"

	_, err := r.db.NamedExecContext(ctx, query, app)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	return nil
}

// GetByID retrieves an application by ID
func (r *ApplicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	var app models.Application
	query := `SELECT * FROM applications WHERE id = $1`

	err := r.db.GetContext(ctx, &app, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return &app, nil
}

// GetByUserID retrieves all applications for a user
func (r *ApplicationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, page, limit int) (*models.ApplicationList, error) {
	var apps []models.Application
	var total int64

	// Get total count
	countQuery := `SELECT COUNT(*) FROM applications WHERE user_id = $1`
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count applications: %w", err)
	}

	// Get applications with pagination
	offset := (page - 1) * limit
	query := `
		SELECT * FROM applications 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	err = r.db.SelectContext(ctx, &apps, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	return &models.ApplicationList{
		Applications: apps,
		Total:        total,
		Page:         page,
		Limit:        limit,
	}, nil
}

// Update updates an application
func (r *ApplicationRepository) Update(ctx context.Context, id uuid.UUID, update *models.ApplicationUpdate) error {
	// Build dynamic query based on provided fields
	query := `UPDATE applications SET updated_at = $1`
	args := []interface{}{time.Now()}
	argIndex := 2

	if update.Name != nil {
		query += fmt.Sprintf(", name = $%d", argIndex)
		args = append(args, *update.Name)
		argIndex++
	}

	if update.Description != nil {
		query += fmt.Sprintf(", description = $%d", argIndex)
		args = append(args, *update.Description)
		argIndex++
	}

	if update.GitURL != nil {
		query += fmt.Sprintf(", git_url = $%d", argIndex)
		args = append(args, *update.GitURL)
		argIndex++
	}

	if update.Branch != nil {
		query += fmt.Sprintf(", branch = $%d", argIndex)
		args = append(args, *update.Branch)
		argIndex++
	}

	if update.Port != nil {
		query += fmt.Sprintf(", port = $%d", argIndex)
		args = append(args, *update.Port)
		argIndex++
	}

	if update.Environment != nil {
		query += fmt.Sprintf(", environment = $%d", argIndex)
		args = append(args, *update.Environment)
		argIndex++
	}

	if update.Replicas != nil {
		query += fmt.Sprintf(", replicas = $%d", argIndex)
		args = append(args, *update.Replicas)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("application not found")
	}

	return nil
}

// UpdateStatus updates application status
func (r *ApplicationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE applications SET status = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("application not found")
	}

	return nil
}

// Delete deletes an application
func (r *ApplicationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM applications WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("application not found")
	}

	return nil
}

// GetAll retrieves all applications with pagination
func (r *ApplicationRepository) GetAll(ctx context.Context, page, limit int) (*models.ApplicationList, error) {
	var apps []models.Application
	var total int64

	// Get total count
	countQuery := `SELECT COUNT(*) FROM applications`
	err := r.db.GetContext(ctx, &total, countQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to count applications: %w", err)
	}

	// Get applications with pagination
	offset := (page - 1) * limit
	query := `
		SELECT * FROM applications 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	err = r.db.SelectContext(ctx, &apps, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get applications: %w", err)
	}

	return &models.ApplicationList{
		Applications: apps,
		Total:        total,
		Page:         page,
		Limit:        limit,
	}, nil
} 