package services

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"

	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type deploymentService struct {
	repo   repository.DeploymentRepository
	client *kubernetes.Clientset
}

func NewDeploymentService(repo repository.DeploymentRepository) DeploymentService {
	client, err := getK8SClient()
	if err != nil {
		return nil
	}
	return &deploymentService{repo: repo, client: client}
}

func (s *deploymentService) CreateDeployment(ctx context.Context, dep *models.Deployment) (*models.Deployment, error) {
	if dep.AppID == uuid.Nil {
		return nil, errors.New("Deployment AppID is required")
	}
	if err := s.repo.Create(ctx, dep); err != nil {
		return nil, err
	}
	return dep, nil
}

func (s *deploymentService) GetDeploymentStatus(ctx context.Context, id uuid.UUID) (string, error) {
	deploy, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	return deploy.Status, nil
}

func (s *deploymentService) GetDeploymentByID(ctx context.Context, id uuid.UUID) (*models.Deployment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *deploymentService) ListAllDeployments(ctx context.Context, f repository.DeploymentFilter, page repository.Page, sort repository.Sort) (repository.ListResult[models.Deployment], error) {
	return s.repo.List(ctx, f, page, sort)
}

func getK8SClient() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}
	return kubernetes.NewForConfig(config)
}

func (s *deploymentService) DeployApp(ctx context.Context, app models.Application) (*models.Deployment, error) {
	// 1. save deployment record
	deploy := &models.Deployment{
		ID:     uuid.New(),
		AppID:  app.ID,
		Status: "PENDING",
	}
	if err := s.repo.Create(ctx, deploy); err != nil {
		return nil, err
	}

	// 2. create resrouce in K8S
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("app-%s", deploy.ID.String()[0:8]),
			Labels: map[string]string{
				"app": app.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": app.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": app.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  app.Name,
							Image: app.ImageURL,
							Ports: []corev1.ContainerPort{{ContainerPort: 8080}},
						},
					},
				},
			},
		},
	}

	_, err := s.client.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		_ = s.repo.UpdateStatus(ctx, deploy.ID, "FAILED")
		return nil, fmt.Errorf("failed to create k8s deployment: %w", err)
	}

	// 3. update status = DEPLOYING
	if err := s.repo.UpdateStatus(ctx, deploy.ID, "DEPLOYING"); err != nil {
		return nil, err
	}

	// 4. Start background goroutine to track status
	go s.trackDeployment(ctx, deploy.ID, app.Name)

	return deploy, nil
}

func (s *deploymentService) trackDeployment(ctx context.Context, deployID uuid.UUID, appName string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// check pods status
		pods, err := s.client.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
			LabelSelector: fmt.Sprintf("app=%s", appName),
		})
		if err != nil {
			_ = s.repo.UpdateStatus(ctx, deployID, "FAILED")
			return
		}

		ready := true
		for _, pod := range pods.Items {
			if pod.Status.Phase != corev1.PodRunning {
				ready = false
				break
			}
		}

		if len(pods.Items) > 0 && ready {
			_ = s.repo.UpdateStatus(ctx, deployID, "RUNNING")
			return
		}
	}
}

// int32Ptr returns a pointer to the given int32 value.
func int32Ptr(i int32) *int32 {
	return &i
}
