package main

import (
	"fmt"
	"mini-paas/backend/internal/repository"
	"mini-paas/backend/internal/services"
)

func main() {
	git := &services.FakeGitService{}
	build := &services.FakeBuildService{}
	registry := &services.FakeRegistryService{}
	deploy := &services.FakeDeployService{}
	repo := repository.NewMemoryRepo()

	appService := services.NewAppService(git, build, registry, deploy, repo)

	id, err := appService.Deploy(services.DeployRequest{
		Name:   "testapp",
		GitURL: "https://github.com/example/test.git",
	})

	if err != nil {
		panic(err)
	}

	app, _ := repo.GetByID(id)
	fmt.Printf("App deployed: %+v\n", app)
}
