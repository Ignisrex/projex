package services

import "project-tracker/server/internal/domain"

type ProjectService interface {
	ListProjects() ([]domain.Project, error)
	GetProject(id string) (*domain.Project, error)
	CreateProject(p domain.Project) (*domain.Project, error)
	UpdateProject(id string, p domain.Project) (*domain.Project, error)
	DeleteProject(id string) error
	ListTasksByProject(projectID string) ([]domain.Task, error)
	CreateTask(projectID string, t domain.Task) (*domain.Task, error)
}

type TaskService interface {
	UpdateTask(id string, t domain.Task) (*domain.Task, error)
	DeleteTask(id string) error
}
