package services

import (
	"project-tracker/server/internal/domain"
	"project-tracker/server/internal/repositories"
)

type taskService struct {
	taskRepo repositories.TaskRepository
}

func NewTaskService(taskRepo repositories.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) UpdateTask(id string, t domain.Task) (*domain.Task, error) {
	return s.taskRepo.UpdateTask(id, t)
}

func (s *taskService) DeleteTask(id string) error {
	return s.taskRepo.DeleteTask(id)
}
