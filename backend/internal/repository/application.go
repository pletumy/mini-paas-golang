package repository

import "time"

// entity object
type Application struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GitURL      string    `json:"git_url"`
	ImageURL    string    `json:"image_url"`
	DeployURL   string    `json:"deploy_url"`
	Status      string    `json:"status"`
}

// work with DB
type ApplicationRepository interface {
	Create(app Application) error
	GetByID(id string) (*Application, error)
	Update(app Application) error
	UpdateStatus(id string, status string) error
	Delete(id string) error
	ListAll() ([]Application, error)
}
