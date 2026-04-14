package main

import (
	"log"
	"net/http"

	"project-tracker/server/internal/config"
	httplayer "project-tracker/server/internal/http"
	"project-tracker/server/internal/repositories"
	"project-tracker/server/internal/services"
)

func main() {
	cfg := config.Load()

	projectRepo := repositories.NewProjectRepositoryStub()
	taskRepo := repositories.NewTaskRepositoryStub()

	projectService := services.NewProjectService(projectRepo, taskRepo)
	taskService := services.NewTaskService(taskRepo)

	router := httplayer.NewRouter(cfg, projectService, taskService)

	log.Printf("server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
