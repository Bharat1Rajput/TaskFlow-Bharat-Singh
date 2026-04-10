package service

import (
	"context"
	"errors"

	"github.com/Bharat1Rajput/taskflow-backend/internal/model"
	"github.com/Bharat1Rajput/taskflow-backend/internal/repository"

	"github.com/google/uuid"
)

var ErrForbidden = errors.New("forbidden")

type ProjectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(r repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: r}
}

func (s *ProjectService) Create(ctx context.Context, name string, desc *string, userID uuid.UUID) error {
	project := model.Project{
		ID:          uuid.New(),
		Name:        name,
		Description: desc,
		OwnerID:     userID,
	}

	return s.repo.Create(ctx, project)
}

func (s *ProjectService) List(ctx context.Context, userID uuid.UUID) ([]model.Project, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, name string, desc *string, userID uuid.UUID) error {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return ErrForbidden
	}

	project.Name = name
	project.Description = desc

	return s.repo.Update(ctx, project)
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}
