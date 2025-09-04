package models

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DeploymentID uuid.UUID `gorm:"type:uuid;not null"`
	Message      string    `gorm:"not null"`
	Level        string    `gorm:"default:INFO"`
	Timestamp    time.Time
}
