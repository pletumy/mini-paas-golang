package services

import (
	"fmt"
	"mini-paas/backend/internal/repository"
)

type FakeGitService struct{}

func (f *FakeGitService) Clone(repoURL string) (string, error) {
	return "/fake/path/to/repo", nil
}

type FakeBuildService struct{}

func (f *FakeBuildService) Build(path string, appName string) (string, error) {
	return fmt.Sprintf("%s:latest", appName), nil
}

type FakeRegistryService struct{}

func (f *FakeRegistryService) Push(imageName string) (string, error) {
	return fmt.Sprintf("registry.local/%s", imageName), nil
}

type FakeDeployService struct{}

func (f *FakeDeployService) Deploy(imageURL string, appName string) (string, error) {
	return fmt.Sprintf("http://%s.local", appName), nil
}

func NewMockAppService(repo repository.ApplicationRepository) *AppService {
	return NewAppService(
		&FakeGitService{},
		&FakeBuildService{},
		&FakeRegistryService{},
		&FakeDeployService{},
		repo,
	)
}
