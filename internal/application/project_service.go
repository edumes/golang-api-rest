package application

import (
	"context"
	"errors"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProjectService struct {
	repo   domain.ProjectRepository
	logger *logrus.Logger
}

func NewProjectService(repo domain.ProjectRepository) *ProjectService {
	return &ProjectService{
		repo:   repo,
		logger: logrus.New(),
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, name, description, status string, startDate, endDate *time.Time, budget *float64, ownerID uuid.UUID) (*domain.Project, error) {
	s.logger.WithFields(logrus.Fields{
		"name":     name,
		"status":   status,
		"owner_id": ownerID,
	}).Info("Creating new project")

	if name == "" {
		s.logger.Warn("Project name is required")
		return nil, errors.New("project name is required")
	}

	if status == "" {
		status = "active"
	}

	project := &domain.Project{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Status:      status,
		StartDate:   startDate,
		EndDate:     endDate,
		Budget:      budget,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
		"owner_id":   project.OwnerID,
	}).Debug("Saving project to repository")

	if err := s.repo.Create(ctx, project); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": project.ID,
			"name":       project.Name,
		}).Error("Failed to create project in repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
		"owner_id":   project.OwnerID,
	}).Info("Project created successfully")

	return project, nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	s.logger.WithFields(logrus.Fields{
		"project_id": id,
	}).Debug("Getting project by ID")

	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": id,
		}).Warn("Project not found by ID")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
		"owner_id":   project.OwnerID,
	}).Debug("Project retrieved successfully")

	return project, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, filter domain.ProjectParams, pagination domain.Pagination) ([]domain.Project, error) {
	s.logger.WithFields(logrus.Fields{
		"filter_name":   filter.Name,
		"filter_status": filter.Status,
		"limit":         pagination.Limit,
		"offset":        pagination.Offset,
		"sort":          pagination.Sort,
	}).Debug("Listing projects with filters")

	projects, err := s.repo.List(ctx, filter, pagination)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list projects from repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"count": len(projects),
	}).Info("Projects listed successfully")

	return projects, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, project *domain.Project) error {
	s.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
		"status":     project.Status,
	}).Info("Updating project")

	project.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, project)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": project.ID,
		}).Error("Failed to update project in repository")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
	}).Info("Project updated successfully")

	return nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	s.logger.WithFields(logrus.Fields{
		"project_id": id,
	}).Info("Deleting project")

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": id,
		}).Error("Failed to delete project from repository")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"project_id": id,
	}).Info("Project deleted successfully")

	return nil
}

func (s *ProjectService) GetProjectsByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]domain.Project, error) {
	s.logger.WithFields(logrus.Fields{
		"owner_id": ownerID,
	}).Debug("Getting projects by owner ID")

	projects, err := s.repo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":    err.Error(),
			"owner_id": ownerID,
		}).Error("Failed to get projects by owner ID from repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"owner_id": ownerID,
		"count":    len(projects),
	}).Info("Projects retrieved successfully by owner ID")

	return projects, nil
}
