package models

import "time"

type Application struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	NAME        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	GitURL      string    `gorm:"type:varchar(255)" json:"git_url"`
	ImageURL    string    `gorm:"type:varchar(255)" json:"image_url"`
	DeployURL   string    `gorm:"type:varchar(255)" json:"deploy_url"`
	Status      string    `gorm:"type:varchar(50);default:'pending'" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
