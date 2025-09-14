package services

import (
	"context"
	"io"

	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
)

type AppService interface {
	CreateApp(ctx context.Context, app *models.Application) (*models.Application, error)
	GetAppByID(ctx context.Context, id uuid.UUID) (*models.Application, error)
	ListApps(ctx context.Context, f repository.AppFilter, page repository.Page, sort repository.Sort) (repository.ListResult[models.Application], error)
	DeleteApp(ctx context.Context, id uuid.UUID) error
}

type DeploymentService interface {
	CreateDeployment(ctx context.Context, dep *models.Deployment) (*models.Deployment, error)
	GetDeploymentByID(ctx context.Context, id uuid.UUID) (*models.Deployment, error)
	ListAllDeployments(ctx context.Context, f repository.DeploymentFilter, page repository.Page, sort repository.Sort) (repository.ListResult[models.Deployment], error)
	// ListDeploymentsByApp(ctx context.Context, appID uuid.UUID, page repository.Page, sort repository.Sort) (repository.ListResult[models.Deployment], error)
	// k8
	DeployApp(ctx context.Context, app models.Application) (*models.Deployment, error)
	GetDeploymentStatus(ctx context.Context, id uuid.UUID) (string, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type LogService interface {
	CreateLog(ctx context.Context, log *models.Log) (*models.Log, error)
	ListAllLogs(ctx context.Context, f repository.LogFilter, limit int) ([]models.Log, error)
	// ListLogsByDeployment(ctx context.Context, depID uuid.UUID) (repository.ListResult[models.Log], error)
	// k8
	StreamPodLogs(ctx context.Context, namespace, podName string, follow bool, tailLines *int64) (<-chan string, error)
}

type K8sLogService interface {
	FindPodsForDeployment(ctx context.Context, deploymentName, namespace string) ([]corev1.Pod, error)
	StreamPodLogs(ctx context.Context, namespace, podName string, follow bool, tailLine *int64) (io.ReadCloser, error)
}
