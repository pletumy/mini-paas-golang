package services

type GitService interface {
	Clone(repoURL string) (string, error)
}

type BuildService interface {
	Build(path string, appName string) (string, error)
}

type RegistryService interface {
	Push(imageName string) (string, error)
}

type DeployService interface {
	Deploy(imageURL string, appName string) (string, error)
}
