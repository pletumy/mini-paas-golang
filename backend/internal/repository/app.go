package repository

import (
	"context"
	"fmt"
	"mini-paas/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppFilter struct {
	OwnerID *uuid.UUID
	Status  string
	Search  string // name search (ILIKE)
}

type AppRepository interface {
	Create(ctx context.Context, app *models.Application) error
	Update(ctx context.Context, app *models.Application) error
	DeleteHard(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error)
	List(ctx context.Context, f AppFilter, page Page, sort Sort) (ListResult[models.Application], error)
	ExistsByNameForOwner(ctx context.Context, ownerID uuid.UUID, name string) (bool, error)
}

type appRepository struct {
	db *gorm.DB
}

func NewAppRepository(db *gorm.DB) AppRepository {
	return &appRepository{db: db}
}

func mapGormError(err error) error {
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func (r *appRepository) Create(ctx context.Context, app *models.Application) error {
	db := getDB(ctx, r.db)
	if err := db.Create(app).Error; err != nil {
		return mapGormError(err)
	}
	return nil
}

func (r *appRepository) Update(ctx context.Context, app *models.Application) error {
	db := getDB(ctx, r.db)

	if err := db.Model(&models.Application{}).
		Where("id = ?", app.ID).
		Updates(map[string]any{
			"name":        app.Name,
			"description": app.Description,
			"repo_url":    app.GitURL,
			"runtime":     app.Runtime,
			"status":      app.Status,
		}).Error; err != nil {
		return mapGormError(err)
	}
	return nil
}

func (r *appRepository) DeleteHard(ctx context.Context, id uuid.UUID) error {
	db := getDB(ctx, r.db)
	if err := db.Delete(&models.Application{}, "id = ?", id).Error; err != nil {
		return mapGormError(err)
	}
	return nil
}

func (r *appRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	db := getDB(ctx, r.db)

	var app models.Application
	if err := db.First(&app, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, mapGormError(err)
	}
	return &app, nil
}

func (r *appRepository) List(ctx context.Context, f AppFilter, page Page, sort Sort) (ListResult[models.Application], error) {
	db := getDB(ctx, r.db).Model(&models.Application{})

	//filters
	if f.OwnerID != nil {
		db = db.Where("owner_id = ?", *f.OwnerID)
	}

	if f.Status != "" {
		db = db.Where("status = ?", f.Status)
	}

	if f.Search != "" {
		like := "%" + f.Search + "%"
		db = db.Where("name ILIKE ?", like)
	}

	//count
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return ListResult[models.Application]{}, mapGormError(err)
	}

	//sorting
	order := "created_at DESC"
	if sort.Field != "" {
		dir := "ASC"
		if sort.Desc {
			dir = "DESC"
		}

		switch sort.Field {
		case "name", "status", "created_at", "updated_at":
			order = fmt.Sprintf("%s %s", sort.Field, dir)
		}
	}
	db = db.Order(order)

	// pagination
	p := page.Sanitize(100)
	var items []models.Application
	if err := db.Limit(p.Limit).Offset(p.Offset).Find(&items).Error; err != nil {
		return ListResult[models.Application]{}, mapGormError(err)
	}
	return ListResult[models.Application]{Items: items, Total: total}, nil
}

func (r *appRepository) ExistsByNameForOwner(ctx context.Context, ownerID uuid.UUID, name string) (bool, error) {
	db := getDB(ctx, r.db)

	var count int64
	if err := db.Model(&models.Application{}).
		Where("owner_id = ? AND name = ?", ownerID, name).
		Count(&count).Error; err != nil {
		return false, mapGormError(err)
	}
	return count > 0, nil

}
