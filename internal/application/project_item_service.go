package application

import (
	"context"
	"errors"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProjectItemService struct {
	repo   domain.ProjectItemRepository
	logger *logrus.Logger
}

func NewProjectItemService(repo domain.ProjectItemRepository) *ProjectItemService {
	return &ProjectItemService{
		repo:   repo,
		logger: logrus.New(),
	}
}

func (s *ProjectItemService) CreateProjectItem(ctx context.Context, projectID uuid.UUID, name, description, status, priority string, estimatedHours, actualHours *float64, dueDate *time.Time, assignedTo *uuid.UUID) (*domain.ProjectItem, error) {
	s.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"name":       name,
		"status":     status,
		"priority":   priority,
	}).Info("Creating new project item")

	if name == "" {
		s.logger.Warn("Project item name is required")
		return nil, errors.New("project item name is required")
	}

	if status == "" {
		status = "pending"
	}

	if priority == "" {
		priority = "medium"
	}

	item := &domain.ProjectItem{
		ID:             uuid.New(),
		ProjectID:      projectID,
		Name:           name,
		Description:    description,
		Status:         status,
		Priority:       priority,
		EstimatedHours: estimatedHours,
		ActualHours:    actualHours,
		DueDate:        dueDate,
		AssignedTo:     assignedTo,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	s.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Debug("Saving project item to repository")

	if err := s.repo.Create(ctx, item); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"item_id":    item.ID,
			"name":       item.Name,
			"project_id": item.ProjectID,
		}).Error("Failed to create project item in repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Info("Project item created successfully")

	return item, nil
}

func (s *ProjectItemService) GetProjectItemByID(ctx context.Context, id uuid.UUID) (*domain.ProjectItem, error) {
	s.logger.WithFields(logrus.Fields{
		"item_id": id,
	}).Debug("Getting project item by ID")

	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"item_id": id,
		}).Warn("Project item not found by ID")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Debug("Project item retrieved successfully")

	return item, nil
}

func (s *ProjectItemService) ListProjectItems(ctx context.Context, filter domain.ProjectItemParams, pagination domain.Pagination) ([]domain.ProjectItem, error) {
	s.logger.WithFields(logrus.Fields{
		"filter_name":     filter.Name,
		"filter_status":   filter.Status,
		"filter_priority": filter.Priority,
		"limit":           pagination.Limit,
		"offset":          pagination.Offset,
		"sort":            pagination.Sort,
	}).Debug("Listing project items with filters")

	items, err := s.repo.List(ctx, filter, pagination)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list project items from repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"count": len(items),
	}).Info("Project items listed successfully")

	return items, nil
}

func (s *ProjectItemService) UpdateProjectItem(ctx context.Context, item *domain.ProjectItem) error {
	s.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"status":     item.Status,
		"project_id": item.ProjectID,
	}).Info("Updating project item")

	item.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, item)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"item_id": item.ID,
		}).Error("Failed to update project item in repository")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Info("Project item updated successfully")

	return nil
}

func (s *ProjectItemService) DeleteProjectItem(ctx context.Context, id uuid.UUID) error {
	s.logger.WithFields(logrus.Fields{
		"item_id": id,
	}).Info("Deleting project item")

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"item_id": id,
		}).Error("Failed to delete project item from repository")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"item_id": id,
	}).Info("Project item deleted successfully")

	return nil
}

func (s *ProjectItemService) GetProjectItemsByProjectID(ctx context.Context, projectID uuid.UUID) ([]domain.ProjectItem, error) {
	s.logger.WithFields(logrus.Fields{
		"project_id": projectID,
	}).Debug("Getting project items by project ID")

	items, err := s.repo.GetByProjectID(ctx, projectID)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": projectID,
		}).Error("Failed to get project items by project ID from repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"count":      len(items),
	}).Info("Project items retrieved successfully by project ID")

	return items, nil
}

func (s *ProjectItemService) GetProjectItemsByAssignedTo(ctx context.Context, assignedTo uuid.UUID) ([]domain.ProjectItem, error) {
	s.logger.WithFields(logrus.Fields{
		"assigned_to": assignedTo,
	}).Debug("Getting project items by assigned user")

	items, err := s.repo.GetByAssignedTo(ctx, assignedTo)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":       err.Error(),
			"assigned_to": assignedTo,
		}).Error("Failed to get project items by assigned user from repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"assigned_to": assignedTo,
		"count":       len(items),
	}).Info("Project items retrieved successfully by assigned user")

	return items, nil
}
