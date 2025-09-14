package services

import (
	"bufio"
	"context"

	"mini-paas/backend/internal/models"
	"mini-paas/backend/internal/repository"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type logService struct {
	repo      repository.LogRepository
	clientset *kubernetes.Clientset
}

func NewLogService(repo repository.LogRepository) LogService {
	return &logService{repo: repo}
}

func (s *logService) CreateLog(ctx context.Context, log *models.Log) (*models.Log, error) {
	if err := s.repo.Append(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

func (s *logService) ListAllLogs(ctx context.Context, f repository.LogFilter, limit int) ([]models.Log, error) {
	return s.repo.List(ctx, f, limit)
}

func (s *logService) StreamPodLogs(ctx context.Context, namespace, podName string, follow bool, tailLines *int64) (<-chan string, error) {
	opts := &corev1.PodLogOptions{
		Follow:    follow,
		TailLines: tailLines,
	}

	req := s.clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
	stream, err := req.Stream(ctx)
	if err != nil {
		return nil, err
	}

	logCh := make(chan string)
	go func() {
		defer stream.Close()
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			logCh <- scanner.Text()
		}
		close(logCh)
	}()
	return logCh, nil
}
