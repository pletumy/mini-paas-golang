package repository

import (
	"context"
	"mini-paas/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeploymentFilter struct {
	AppID  *uuid.UUID
	Status *string
}

type DeploymentRepository interface {
	Create(ctx context.Context, d *models.Deployment) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Deployment, error)
	List(ctx context.Context, f DeploymentFilter, page Page, sort Sort) (ListResult[models.Deployment], error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}
type deploymentRepository struct{ db *gorm.DB }

func NewDeploymentRepository(db *gorm.DB) DeploymentRepository {
	return &deploymentRepository{db: db}
}

func (r *deploymentRepository) Create(ctx context.Context, d *models.Deployment) error {
	return getDB(ctx, r.db).Create(d).Error
}

func (r *deploymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Deployment, error) {
	var d models.Deployment
	if err := getDB(ctx, r.db).First(d, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *deploymentRepository) List(ctx context.Context, f DeploymentFilter, page Page, sort Sort) (ListResult[models.Deployment], error) {
	db := getDB(ctx, r.db).Model(&models.Deployment{})
	if f.Status != nil {
		db = db.Where("status = ?", *f.Status)
	}

	if f.AppID != nil {
		db = db.Where("app_id = ?", *f.AppID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return ListResult[models.Deployment]{}, err
	}

	order := "created_at DESC"
	if sort.Field != "" {
		order = sort.Field
		if sort.Desc {
			order += " DESC"
		}
	}

	p := page.Sanitize(100)

	var items []models.Deployment
	if err := db.Order(order).Limit(p.Limit).Offset(p.Offset).Find(&items).Error; err != nil {
		return ListResult[models.Deployment]{}, err
	}
	return ListResult[models.Deployment]{Items: items, Total: total}, nil
}

func (r *deploymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return getDB(ctx, r.db).Model(&models.Deployment{}).Where("id = ?", id).Update("status", status).Error
}
