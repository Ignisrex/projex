package repositories

import (
	"errors"

	"project-tracker/server/internal/domain"
)

var ErrNotImplemented = errors.New("repository method not implemented")

type projectRepositoryStub struct{}

func NewProjectRepositoryStub() ProjectRepository {
	return &projectRepositoryStub{}
}

func (r *projectRepositoryStub) ListProjects() ([]domain.Project, error) {
	return nil, ErrNotImplemented
}

func (r *projectRepositoryStub) GetProject(id string) (*domain.Project, error) {
	return nil, ErrNotImplemented
}

func (r *projectRepositoryStub) CreateProject(p domain.Project) (*domain.Project, error) {
	return nil, ErrNotImplemented
}

func (r *projectRepositoryStub) UpdateProject(id string, p domain.Project) (*domain.Project, error) {
	return nil, ErrNotImplemented
}

func (r *projectRepositoryStub) DeleteProject(id string) error {
	return ErrNotImplemented
}

type taskRepositoryStub struct{}

func NewTaskRepositoryStub() TaskRepository {
	return &taskRepositoryStub{}
}

func (r *taskRepositoryStub) ListTasksByProject(projectID string) ([]domain.Task, error) {
	return nil, ErrNotImplemented
}

func (r *taskRepositoryStub) CreateTask(projectID string, t domain.Task) (*domain.Task, error) {
	return nil, ErrNotImplemented
}

func (r *taskRepositoryStub) UpdateTask(id string, t domain.Task) (*domain.Task, error) {
	return nil, ErrNotImplemented
}

func (r *taskRepositoryStub) DeleteTask(id string) error {
	return ErrNotImplemented
}
