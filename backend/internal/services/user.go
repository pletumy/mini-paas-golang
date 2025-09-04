package services

import (
	"context"
	"errors"
	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"

	"github.com/google/uuid"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) userService {
	return userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if user.Email == "" {
		return nil, errors.New("User email is required")
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetByEmail(ctx, email)
}
