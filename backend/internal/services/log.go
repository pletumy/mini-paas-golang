package services

import (
	"context"
	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"
)

type logService struct {
	repo repository.LogRepository
}

func NewLogService(repo repository.LogRepository) LogService {
	return &logService{repo: repo}
}

func (s *logService) CreateLog(ctx context.Context, log *models.Log) (*models.Log, error) {
	if err := s.repo.Append(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *logService) ListAllLogs(ctx context.Context, f repository.LogFilter, limit int) ([]models.Log, error) {
	return s.repo.List(ctx, f, limit)
}
