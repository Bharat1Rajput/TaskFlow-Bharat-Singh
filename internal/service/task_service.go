package service

import (
	"context"

	"github.com/Bharat1Rajput/taskflow-backend/internal/model"
	"github.com/Bharat1Rajput/taskflow-backend/internal/repository"
	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo    repository.TaskRepository
	projectRepo repository.ProjectRepository
}

func NewTaskService(
	taskRepo repository.TaskRepository,
	projectRepo repository.ProjectRepository,
) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (s *TaskService) Create(ctx context.Context, t model.Task) error {
	// Ensure project exists
	_, err := s.projectRepo.GetByID(ctx, t.ProjectID)
	if err != nil {
		return err
	}

	t.ID = uuid.New()

	// Defaults
	if t.Status == "" {
		t.Status = "todo"
	}
	if t.Priority == "" {
		t.Priority = "medium"
	}

	return s.taskRepo.Create(ctx, t)
}

func (s *TaskService) List(
	ctx context.Context,
	projectID uuid.UUID,
	status *string,
	assignee *uuid.UUID,
) ([]model.Task, error) {

	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return s.taskRepo.ListByProject(ctx, projectID, status, assignee)
}

func (s *TaskService) Update(
	ctx context.Context,
	t model.Task,
	userID uuid.UUID,
) error {

	existing, err := s.taskRepo.GetByID(ctx, t.ID)
	if err != nil {
		return err
	}

	project, err := s.projectRepo.GetByID(ctx, existing.ProjectID)
	if err != nil {
		return err
	}

	// 🔐 Only project owner allowed
	if project.OwnerID != userID {
		return ErrForbidden
	}

	// Preserve immutable field
	t.ProjectID = existing.ProjectID

	return s.taskRepo.Update(ctx, t)
}

func (s *TaskService) Delete(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) error {

	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	project, err := s.projectRepo.GetByID(ctx, task.ProjectID)
	if err != nil {
		return err
	}

	// 🔐 Only project owner allowed
	if project.OwnerID != userID {
		return ErrForbidden
	}

	return s.taskRepo.Delete(ctx, id)
}
