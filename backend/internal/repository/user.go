package repository

import (
	"context"

	"mini-paas/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, u *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, u *models.User) error
}

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *models.User) error {
	return getDB(ctx, r.db).Create(u).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var u models.User
	if err := getDB(ctx, r.db).First(&u, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	if err := getDB(ctx, r.db).First(&u, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, u *models.User) error {
	return getDB(ctx, r.db).Model(&models.User{}).
		Where("id = ?", u.ID).
		Updates(map[string]any{
			"name":          u.Name,
			"password_hash": u.PasswordHash,
			"role":          u.Name,
		}).Error
}
