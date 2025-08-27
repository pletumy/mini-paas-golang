package services

import (
	"mini-paas/backend/internal/repository"
	"time"

	"github.com/google/uuid"
)

type AppService struct {
	gitService      GitService
	buildService    BuildService
	registryService RegistryService
	deployService   DeployService
	repo            repository.ApplicationRepository
}

func NewAppService(
	g GitService,
	b BuildService,
	r RegistryService,
	d DeployService,
	repo repository.ApplicationRepository,
) *AppService {
	return &AppService{
		gitService:      g,
		buildService:    b,
		registryService: r,
		deployService:   d,
		repo:            repo,
	}
}

type DeployRequest struct {
	Name   string
	GitURL string
}

// func (s *AppService) Repo() repository.ApplicationRepository {
// 	return s.repo
// }

func (s *AppService) Deploy(req DeployRequest) (string, error) {
	id := uuid.New().String()
	app := repository.Application{
		ID:          id,
		Name:        req.Name,
		Description: "Deployed via mini-paas",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		GitURL:      req.GitURL,
		Status:      "PENDING",
	}

	if err := s.repo.Create(app); err != nil {
		return "", err
	}

	// clone repo
	path, err := s.gitService.Clone(req.GitURL)
	if err != nil {
		s.repo.UpdateStatus(id, "CLONE_FAILED")
		return "", err
	}

	// build img
	imageName, err := s.buildService.Build(path, req.Name)
	if err != nil {
		s.repo.UpdateStatus(id, "BUILD_FAILED")
		return "", err
	}

	// push img
	imageURL, err := s.registryService.Push(imageName)
	if err != nil {
		s.repo.UpdateStatus(id, "PUSH_FAILED")
		return "", err
	}

	// deploy to k8s
	deployURL, err := s.deployService.Deploy(imageURL, req.Name)
	if err != nil {
		s.repo.UpdateStatus(id, "DEPLOY_FAILED")
		return "", err
	}

	// update to DB
	app.ImageURL = imageURL
	app.DeployURL = deployURL
	app.Status = "DEPLOYED"
	if err := s.repo.Update(app); err != nil {
		return "", err
	}

	return id, nil
}

func (s *AppService) GetAllApps() ([]repository.Application, error) {
	return s.repo.ListAll()
}

func (s *AppService) GetAppByID(id string) (*repository.Application, error) {
	return s.repo.GetByID(id)
}

func (s *AppService) DeleteApp(id string) error {
	return s.repo.Delete(id)
}
