package repositories

import "project-tracker/server/internal/domain"

type ProjectRepository interface {
	ListProjects() ([]domain.Project, error)
	GetProject(id string) (*domain.Project, error)
	CreateProject(p domain.Project) (*domain.Project, error)
	UpdateProject(id string, p domain.Project) (*domain.Project, error)
	DeleteProject(id string) error
}

type TaskRepository interface {
	ListTasksByProject(projectID string) ([]domain.Task, error)
	CreateTask(projectID string, t domain.Task) (*domain.Task, error)
	UpdateTask(id string, t domain.Task) (*domain.Task, error)
	DeleteTask(id string) error
}
