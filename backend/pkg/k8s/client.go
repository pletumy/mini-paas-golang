package k8s

import (
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewClientFromKubeConfig() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("build kubeconfig: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("new clientset: %w", err)
	}
	return clientset, nil
}
