package services

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"mini-paas/backend/internal/repository"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var ErrPodNotFound = errors.New("pod not found")

type k8sLogService struct {
	client     *kubernetes.Clientset
	deployRepo repository.DeploymentRepository
}

func NewK8sLogService(client *kubernetes.Clientset, deployRepo repository.DeploymentRepository) K8sLogService {
	return &k8sLogService{client: client, deployRepo: deployRepo}
}

func (s *k8sLogService) FindPodsForDeployment(ctx context.Context, deploymentName, namespace string) ([]corev1.Pod, error) {
	selector := metav1.LabelSelector{MatchLabels: map[string]string{"app": deploymentName}}
	lo, _ := metav1.LabelSelectorAsSelector(&selector)
	podList, err := s.client.CoreV1().
		Pods(namespace).
		List(
			ctx, metav1.ListOptions{LabelSelector: lo.String()},
		)
	if err != nil {
		return nil, fmt.Errorf("list pods: %w", err)
	}
	if len(podList.Items) == 0 {
		podList, err := s.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{FieldSelector: ""})
		if err != nil {
			return nil, err
		}
		var found []corev1.Pod
		for _, p := range podList.Items {
			if len(p.Name) >= len(deploymentName) && p.Name[:len(deploymentName)] == deploymentName {
				found = append(found, p)
			}
		}

		if len(found) == 0 {
			return nil, ErrPodNotFound
		}
		return found, nil
	}
	return podList.Items, err
}

func (s *k8sLogService) StreamPodLogs(
	ctx context.Context,
	namespace,
	podName string,
	follow bool,
	tailLine *int64,
) (io.ReadCloser, error) {
	opts := &corev1.PodLogOptions{
		Follow:     follow,
		Timestamps: false,
	}
	if tailLine != nil {
		opts.TailLines = tailLine
	}
	req := s.client.CoreV1().Pods(namespace).GetLogs(podName, opts)
	stream, err := req.Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("get logs stream: %w", err)
	}
	return stream, nil
}

func StreamToLines(ctx context.Context, r io.ReadCloser, lineCh chan<- string) {
	defer r.Close()
	defer close(lineCh)
	scanner := bufio.NewScanner(r)

	const maxCapicity = 1024 * 1024
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, maxCapicity)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		case lineCh <- scanner.Text():
		}
	}
}
