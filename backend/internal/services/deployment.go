package services

import (
	"context"
	"errors"
	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"

	"github.com/google/uuid"
)

type deploymentService struct {
	repo repository.DeploymentRepository
}

func NewDeploymentRepository(repo repository.DeploymentRepository) DeploymentService {
	return &deploymentService{repo: repo}
}

func (s *deploymentService) CreateDeployment(ctx context.Context, dep *models.Deployment) (*models.Deployment, error) {
	if dep.AppID == uuid.Nil {
		return nil, errors.New("Deployment AppID is required")
	}
	if err := s.repo.Create(ctx, dep); err != nil {
		return nil, err
	}
	return dep, nil
}

func (s *deploymentService) GetDeploymentByID(ctx context.Context, id uuid.UUID) (*models.Deployment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *deploymentService) ListAllDeployments(ctx context.Context, f repository.DeploymentFilter, page repository.Page, sort repository.Sort) (repository.ListResult[models.Deployment], error) {
	return s.repo.List(ctx, f, page, sort)
}
