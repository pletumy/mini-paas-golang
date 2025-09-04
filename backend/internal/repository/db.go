package repository

import (
	"mini-paas/backend/internal/models"

	"gorm.io/gorm"
)

type DBRepository struct {
	db *gorm.DB
}

func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}

func (r *DBRepository) Create(app models.Application) error {
	return r.db.Create(&app).Error
}

func (r *DBRepository) Update(app models.Application) error {
	return r.db.Save(&app).Error
}

func (r *DBRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Application{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *DBRepository) GetByID(id string) (*models.Application, error) {
	var app models.Application
	err := r.db.First(&app, "id = ?", id).Error
	return &app, err
}

func (r *DBRepository) ListAll() ([]models.Application, error) {
	var apps []models.Application
	err := r.db.Find(&apps).Error
	return apps, err
}

func (r *DBRepository) Delete(id string) error {
	return r.db.Delete(&models.Application{}, "id = ?", id).Error
}

// sql thuan
// func (r *DBRepository) Create(app Application) error {
// 	query := `INSERT INTO applications (id, name, description, git_url, image_url, deploy_url, status, created_at, updated_at)
//               VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
// 	_, err := r.db.ExecContext(context.Background(),
// 		query, app.ID, app.Name, app.Description, app.GitURL,
// 		app.ImageURL, app.DeployURL, app.Status,
// 		app.CreatedAt, app.UpdatedAt)
// 	return err
// }

// func (r *DBRepository) Update(app Application) error {
// 	query := `UPDATE applications
// 				SET name=$2, description=$3, git_url=$4, image_url=$5, deploy_url=$6, status=$7, updated_at=$8
// 				where id=$1`
// 	_, err := r.db.ExecContext(context.Background(),
// 		query,
// 		app.ID, app.Name, app.Description, app.GitURL,
// 		app.ImageURL, app.DeployURL, app.Status,
// 		app.UpdatedAt,
// 	)
// 	return err
// }

// func (r *DBRepository) UpdateStatus(id string, status string) error {
// 	query := `UPDATE applications SET status=$2, updated_at=NOW() WHERE id=$1`
// 	_, err := r.db.ExecContext(context.Background(), query, id, status)
// 	return err
// }

// func (r *DBRepository) GetByID(id string) (Application, error) {
// 	query := `SELECT id, name, description, git_url, image_url, deploy_url, status, created_at, updated_at
// 				FROM applications WHERE id=$1`
// 	row := r.db.QueryRowContext(context.Background(), query, id)
// 	var app Application
// 	err := row.Scan(
// 		&app.ID, &app.Name, &app.Description, &app.GitURL, &app.ImageURL, &app.DeployURL, &app.Status, &app.CreatedAt, &app.UpdatedAt
// 	)
// 	return app, err
// }

// func (r *DBRepository) ListAll() ([]Application, error) {
// 	query := `SELECT id, name, description, git_url, image_url, deploy_url, status, created_at, updated_at FROM applications`
// 	rows, err := r.db.QueryContext(context.Background(), query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var apps []Application
// 	for rows.Next() {
// 		var app Application
// 		if err := rows.Scan(
// 			&app.ID, &app.Name, &app.Description, &app.GitURL, &app.ImageURL, &app.DeployURL, &app.Status, &app.CreatedAt, &app.UpdatedAt,
// 		); err != nil {
// 			return nil, err
// 		}
// 		apps = append(apps, app)
// 	}
// 	return apps, nil
// }

// func (r *DBRepository) Delete(id string) error {
// 	query := `DELETE FROM applications where id = $1`
// 	_, err := r.db.ExecContext(context.Background(), query, id)
// 	return err
// }
