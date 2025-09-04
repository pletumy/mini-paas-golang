package repository

import (
	"context"
	"mini-paas/backend/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LogFilter struct {
	DeploymentID uuid.UUID
	Level        *string
	SinceTime    *time.Time
	AfterID      *uint // phan trang theo id tang dan
}

type LogRepository interface {
	Append(ctx context.Context, l *models.Log) error
	List(ctx context.Context, f LogFilter, limit int) ([]models.Log, error)
}

type logRepository struct{ db *gorm.DB }

func NewLogRepository(db *gorm.DB) LogRepository { return &logRepository{db: db} }

func (r *logRepository) Append(ctx context.Context, l *models.Log) error {
	return getDB(ctx, r.db).Create(l).Error
}

func (r *logRepository) List(ctx context.Context, f LogFilter, limit int) ([]models.Log, error) {
	db := getDB(ctx, r.db).Model(&models.Log{}).Where("deployment_id = ?", f.DeploymentID)

	if f.Level != nil && *f.Level != "" {
		db = db.Where("level = ?", *f.Level)
	}

	if f.SinceTime != nil {
		db = db.Where("timestamp >= ?", *f.SinceTime)
	}

	if f.AfterID != nil {
		db = db.Where("id > ?", *f.AfterID)
	}

	if limit <= 0 || limit > 500 {
		limit = 100
	}

	var items []models.Log
	if err := db.Order("id ASC").Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
