package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"project-tracker/server/internal/dbgen"
	"project-tracker/server/internal/domain"
)

type taskRepository struct {
	q *dbgen.Queries
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{q: dbgen.New(db)}
}

func (r *taskRepository) ListTasksByProject(projectID string) ([]domain.Task, error) {
	dbID, err := strconv.ParseInt(projectID, 10, 64)
	if err != nil {
		return nil, ErrNotFound
	}
	rows, err := r.q.ListTasksByProject(context.TODO(), dbID)
	if err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, 0, len(rows))
	for _, row := range rows {
		tasks = append(tasks, taskFromDB(row))
	}
	return tasks, nil
}

func (r *taskRepository) CreateTask(projectID string, t domain.Task) (*domain.Task, error) {
	dbProjectID, err := strconv.ParseInt(projectID, 10, 64)
	if err != nil {
		return nil, ErrNotFound
	}
	row, err := r.q.CreateTask(context.TODO(), dbgen.CreateTaskParams{
		ProjectID:   dbProjectID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Priority:    int32(t.Priority),
	})
	if err != nil {
		return nil, err
	}
	created := taskFromDB(row)
	return &created, nil
}

func (r *taskRepository) UpdateTask(id string, t domain.Task) (*domain.Task, error) {
	dbID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrNotFound
	}
	row, err := r.q.UpdateTask(context.TODO(), dbgen.UpdateTaskParams{
		ID:          dbID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Priority:    int32(t.Priority),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	updated := taskFromDB(row)
	return &updated, nil
}

func (r *taskRepository) DeleteTask(id string) error {
	dbID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return ErrNotFound
	}
	affected, err := r.q.DeleteTask(context.TODO(), dbID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func taskFromDB(row dbgen.Task) domain.Task {
	return domain.Task{
		ID:          strconv.FormatInt(row.ID, 10),
		ProjectID:   strconv.FormatInt(row.ProjectID, 10),
		Title:       row.Title,
		Description: row.Description,
		Status:      row.Status,
		Priority:    int(row.Priority),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
