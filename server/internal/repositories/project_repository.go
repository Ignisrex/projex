package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"project-tracker/server/internal/dbgen"
	"project-tracker/server/internal/domain"
)

type projectRepository struct {
	q *dbgen.Queries
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{q: dbgen.New(db)}
}

func (r *projectRepository) ListProjects() ([]domain.Project, error) {
	rows, err := r.q.ListProjects(context.TODO())
	if err != nil {
		return nil, err
	}
	projects := make([]domain.Project, 0, len(rows))
	for _, row := range rows {
		projects = append(projects, projectFromDB(row))
	}
	return projects, nil
}

func (r *projectRepository) GetProject(id string) (*domain.Project, error) {
	dbID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrNotFound
	}
	row, err := r.q.GetProject(context.TODO(), dbID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	p := projectFromDB(row)
	return &p, nil
}

func (r *projectRepository) CreateProject(p domain.Project) (*domain.Project, error) {
	row, err := r.q.CreateProject(context.TODO(), dbgen.CreateProjectParams{
		Name:        p.Name,
		Description: p.Description,
	})
	if err != nil {
		return nil, err
	}
	created := projectFromDB(row)
	return &created, nil
}

func (r *projectRepository) UpdateProject(id string, p domain.Project) (*domain.Project, error) {
	dbID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrNotFound
	}
	row, err := r.q.UpdateProject(context.TODO(), dbgen.UpdateProjectParams{
		ID:          dbID,
		Name:        p.Name,
		Description: p.Description,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	updated := projectFromDB(row)
	return &updated, nil
}

func (r *projectRepository) DeleteProject(id string) error {
	dbID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return ErrNotFound
	}
	affected, err := r.q.DeleteProject(context.TODO(), dbID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func projectFromDB(row dbgen.Project) domain.Project {
	return domain.Project{
		ID:          strconv.FormatInt(row.ID, 10),
		Name:        row.Name,
		Description: row.Description,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
