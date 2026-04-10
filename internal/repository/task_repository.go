package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Bharat1Rajput/taskflow-backend/internal/model"

	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, t model.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (model.Task, error)
	ListByProject(ctx context.Context, projectID uuid.UUID, status *string, assignee *uuid.UUID) ([]model.Task, error)
	Update(ctx context.Context, t model.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, t model.Task) error {
	query := `
		INSERT INTO tasks 
		(id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW())
	`

	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.Title, t.Description, t.Status, t.Priority,
		t.ProjectID, t.AssigneeID, t.DueDate,
	)

	return err
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Task, error) {
	query := `
		SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
		FROM tasks WHERE id=$1
	`

	var t model.Task

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
		&t.ProjectID, &t.AssigneeID, &t.DueDate,
		&t.CreatedAt, &t.UpdatedAt,
	)

	return t, err
}

func (r *taskRepository) ListByProject(ctx context.Context, projectID uuid.UUID, status *string, assignee *uuid.UUID) ([]model.Task, error) {
	query := `
		SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at
		FROM tasks
		WHERE project_id = $1
	`

	args := []interface{}{projectID}
	argIdx := 2

	if status != nil {
		query += " AND status = $" + fmt.Sprint(argIdx)
		args = append(args, *status)
		argIdx++
	}

	if assignee != nil {
		query += " AND assignee_id = $" + fmt.Sprint(argIdx)
		args = append(args, *assignee)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var t model.Task
		err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority,
			&t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, t model.Task) error {
	query := `
		UPDATE tasks
		SET title = $1,
		    description = $2,
		    status = $3,
		    priority = $4,
		    assignee_id = $5,
		    due_date = $6,
		    updated_at = NOW()
		WHERE id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		t.Title,
		t.Description,
		t.Status,
		t.Priority,
		t.AssigneeID,
		t.DueDate,
		t.ID,
	)

	return err
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
