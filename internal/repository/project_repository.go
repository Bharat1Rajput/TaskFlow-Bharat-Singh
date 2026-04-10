package repository

import (
	"context"
	"database/sql"

	"github.com/Bharat1Rajput/taskflow-backend/internal/model"

	"github.com/google/uuid"
)

type ProjectRepository interface {
	Create(ctx context.Context, p model.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (model.Project, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Project, error)
	Update(ctx context.Context, p model.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, p model.Project) error {
	query := `
		INSERT INTO projects (id, name, description, owner_id, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.db.ExecContext(ctx, query, p.ID, p.Name, p.Description, p.OwnerID)
	return err
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Project, error) {
	query := `
		SELECT id, name, description, owner_id, created_at
		FROM projects
		WHERE id = $1
	`

	var p model.Project

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.OwnerID,
		&p.CreatedAt,
	)

	return p, err
}

func (r *projectRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Project, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.owner_id, p.created_at
		FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		WHERE p.owner_id = $1 OR t.assignee_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project

	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (r *projectRepository) Update(ctx context.Context, p model.Project) error {
	query := `
		UPDATE projects
		SET name = $1, description = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, p.Name, p.Description, p.ID)
	return err
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
