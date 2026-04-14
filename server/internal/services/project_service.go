package services

import (
	"project-tracker/server/internal/domain"
	"project-tracker/server/internal/repositories"
)

type projectService struct {
	projectRepo repositories.ProjectRepository
	taskRepo    repositories.TaskRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository, taskRepo repositories.TaskRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
	}
}

func (s *projectService) ListProjects() ([]domain.Project, error) {
	return s.projectRepo.ListProjects()
}

func (s *projectService) GetProject(id string) (*domain.Project, error) {
	return s.projectRepo.GetProject(id)
}

func (s *projectService) CreateProject(p domain.Project) (*domain.Project, error) {
	return s.projectRepo.CreateProject(p)
}

func (s *projectService) UpdateProject(id string, p domain.Project) (*domain.Project, error) {
	return s.projectRepo.UpdateProject(id, p)
}

func (s *projectService) DeleteProject(id string) error {
	return s.projectRepo.DeleteProject(id)
}

func (s *projectService) ListTasksByProject(projectID string) ([]domain.Task, error) {
	return s.taskRepo.ListTasksByProject(projectID)
}

func (s *projectService) CreateTask(projectID string, t domain.Task) (*domain.Task, error) {
	return s.taskRepo.CreateTask(projectID, t)
}
