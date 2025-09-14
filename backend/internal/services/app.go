package services

import (
	"context"
	"errors"

	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"

	"github.com/google/uuid"
)

type appService struct {
	repo repository.AppRepository
}

func NewAppService(repo repository.AppRepository) AppService {
	return &appService{repo: repo}
}

func (s *appService) CreateApp(ctx context.Context, app *models.Application) (*models.Application, error) {
	if app.Name == "" {
		return nil, errors.New("Application Namm is required")
	}
	if err := s.repo.Create(ctx, app); err != nil {
		return nil, err
	}
	return app, nil
}

func (s *appService) GetAppByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *appService) ListApps(ctx context.Context, f repository.AppFilter, page repository.Page, sort repository.Sort) (repository.ListResult[models.Application], error) {
	return s.repo.List(ctx, f, page, sort)
}

func (s *appService) DeleteApp(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteHard(ctx, id)
}
