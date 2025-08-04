package services

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"mini-paas-golang/backend/internal/models"
	"mini-paas-golang/backend/internal/repository"
	"mini-paas-golang/backend/pkg/docker"
	"mini-paas-golang/backend/pkg/kubernetes"
)

// ApplicationService handles business logic for applications
type ApplicationService struct {
	appRepo     *repository.ApplicationRepository
	dockerSvc   *docker.Service
	k8sSvc      *kubernetes.Service
}

// NewApplicationService creates a new application service
func NewApplicationService(
	appRepo *repository.ApplicationRepository,
	dockerSvc *docker.Service,
	k8sSvc *kubernetes.Service,
) *ApplicationService {
	return &ApplicationService{
		appRepo:   appRepo,
		dockerSvc: dockerSvc,
		k8sSvc:    k8sSvc,
	}
}

// CreateApplication creates a new application and starts deployment process
func (s *ApplicationService) CreateApplication(ctx context.Context, createReq *models.ApplicationCreate, userID uuid.UUID) (*models.Application, error) {
	// Create application in database
	app := &models.Application{
		Name:        createReq.Name,
		Description: createReq.Description,
		GitURL:      createReq.GitURL,
		Branch:      createReq.Branch,
		Port:        createReq.Port,
		Environment: createReq.Environment,
		Replicas:    createReq.Replicas,
		UserID:      userID,
	}

	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	// Start async deployment process
	go s.deployApplication(context.Background(), app)

	return app, nil
}

// GetApplication retrieves an application by ID
func (s *ApplicationService) GetApplication(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	return s.appRepo.GetByID(ctx, id)
}

// GetUserApplications retrieves all applications for a user
func (s *ApplicationService) GetUserApplications(ctx context.Context, userID uuid.UUID, page, limit int) (*models.ApplicationList, error) {
	return s.appRepo.GetByUserID(ctx, userID, page, limit)
}

// UpdateApplication updates an application
func (s *ApplicationService) UpdateApplication(ctx context.Context, id uuid.UUID, update *models.ApplicationUpdate) error {
	// Check if application exists
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Update in database
	if err := s.appRepo.Update(ctx, id, update); err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	// If significant changes were made, redeploy
	if shouldRedeploy(update) {
		go s.redeployApplication(context.Background(), app)
	}

	return nil
}

// DeleteApplication deletes an application and cleans up resources
func (s *ApplicationService) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	// Get application details for cleanup
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Delete from Kubernetes
	if err := s.k8sSvc.DeleteDeployment(ctx, app.Name); err != nil {
		log.Printf("Failed to delete Kubernetes deployment: %v", err)
	}

	// Delete from database
	if err := s.appRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	return nil
}

// DeployApplication manually triggers deployment
func (s *ApplicationService) DeployApplication(ctx context.Context, id uuid.UUID) error {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	go s.deployApplication(context.Background(), app)
	return nil
}

// GetApplicationLogs retrieves application logs
func (s *ApplicationService) GetApplicationLogs(ctx context.Context, id uuid.UUID, lines int) ([]string, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	return s.k8sSvc.GetPodLogs(ctx, app.Name, lines)
}

// GetApplicationMetrics retrieves application metrics
func (s *ApplicationService) GetApplicationMetrics(ctx context.Context, id uuid.UUID) (*models.ApplicationMetrics, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("application not found: %w", err)
	}

	// Get metrics from Kubernetes
	cpuUsage, memoryUsage, err := s.k8sSvc.GetPodMetrics(ctx, app.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics: %w", err)
	}

	return &models.ApplicationMetrics{
		ID:            uuid.New(),
		ApplicationID: app.ID,
		CPUUsage:      cpuUsage,
		MemoryUsage:   memoryUsage,
		Timestamp:     app.UpdatedAt,
	}, nil
}

// deployApplication handles the deployment process
func (s *ApplicationService) deployApplication(ctx context.Context, app *models.Application) {
	// Update status to building
	if err := s.appRepo.UpdateStatus(ctx, app.ID, "building"); err != nil {
		log.Printf("Failed to update status to building: %v", err)
		return
	}

	// Build Docker image
	imageURL, err := s.dockerSvc.BuildImage(ctx, app)
	if err != nil {
		log.Printf("Failed to build Docker image: %v", err)
		s.appRepo.UpdateStatus(ctx, app.ID, "failed")
		return
	}

	// Update image URL
	app.ImageURL = imageURL
	if err := s.appRepo.Update(ctx, app.ID, &models.ApplicationUpdate{
		ImageURL: &imageURL,
	}); err != nil {
		log.Printf("Failed to update image URL: %v", err)
	}

	// Update status to deploying
	if err := s.appRepo.UpdateStatus(ctx, app.ID, "deploying"); err != nil {
		log.Printf("Failed to update status to deploying: %v", err)
		return
	}

	// Deploy to Kubernetes
	deploymentURL, err := s.k8sSvc.DeployApplication(ctx, app)
	if err != nil {
		log.Printf("Failed to deploy to Kubernetes: %v", err)
		s.appRepo.UpdateStatus(ctx, app.ID, "failed")
		return
	}

	// Update deployment URL
	app.DeploymentURL = deploymentURL
	if err := s.appRepo.Update(ctx, app.ID, &models.ApplicationUpdate{
		DeploymentURL: &deploymentURL,
	}); err != nil {
		log.Printf("Failed to update deployment URL: %v", err)
	}

	// Update status to running
	if err := s.appRepo.UpdateStatus(ctx, app.ID, "running"); err != nil {
		log.Printf("Failed to update status to running: %v", err)
	}
}

// redeployApplication handles the redeployment process
func (s *ApplicationService) redeployApplication(ctx context.Context, app *models.Application) {
	// Update status to deploying
	if err := s.appRepo.UpdateStatus(ctx, app.ID, "deploying"); err != nil {
		log.Printf("Failed to update status to deploying: %v", err)
		return
	}

	// Redeploy to Kubernetes
	deploymentURL, err := s.k8sSvc.DeployApplication(ctx, app)
	if err != nil {
		log.Printf("Failed to redeploy to Kubernetes: %v", err)
		s.appRepo.UpdateStatus(ctx, app.ID, "failed")
		return
	}

	// Update deployment URL if changed
	if deploymentURL != app.DeploymentURL {
		if err := s.appRepo.Update(ctx, app.ID, &models.ApplicationUpdate{
			DeploymentURL: &deploymentURL,
		}); err != nil {
			log.Printf("Failed to update deployment URL: %v", err)
		}
	}

	// Update status to running
	if err := s.appRepo.UpdateStatus(ctx, app.ID, "running"); err != nil {
		log.Printf("Failed to update status to running: %v", err)
	}
}

// shouldRedeploy determines if changes require redeployment
func shouldRedeploy(update *models.ApplicationUpdate) bool {
	return update.GitURL != nil || update.Branch != nil || update.Port != nil || 
		   update.Environment != nil || update.Replicas != nil
} 