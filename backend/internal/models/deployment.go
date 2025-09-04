package models

import (
	"time"

	"github.com/google/uuid"
)

type Deployment struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	AppID      uuid.UUID `gorm:"type:uuid;not null"`
	Version    string    `gorm:"not null"`
	ImageURL   string
	Status     string `gorm:"default:pending"`
	DeployedAt time.Time
	CreatedAt  time.Time
}
