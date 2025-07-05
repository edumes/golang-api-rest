package infrastructure

import (
	"context"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostgresProjectItemRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewPostgresProjectItemRepository(db *gorm.DB) *PostgresProjectItemRepository {
	return &PostgresProjectItemRepository{
		db:     db,
		logger: logrus.New(),
	}
}

func (r *PostgresProjectItemRepository) Create(ctx context.Context, item *domain.ProjectItem) error {
	r.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Debug("Creating project item in database")

	err := r.db.WithContext(ctx).Create(item).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"item_id":    item.ID,
			"name":       item.Name,
			"project_id": item.ProjectID,
		}).Error("Failed to create project item in database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Debug("Project item created successfully in database")

	return nil
}

func (r *PostgresProjectItemRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProjectItem, error) {
	r.logger.WithFields(logrus.Fields{
		"item_id": id,
	}).Debug("Getting project item by ID from database")

	var item domain.ProjectItem
	err := r.db.WithContext(ctx).First(&item, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"item_id": id,
		}).Warn("Project item not found in database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Debug("Project item retrieved successfully from database")

	return &item, nil
}

func (r *PostgresProjectItemRepository) List(ctx context.Context, filter domain.ProjectItemParams, pagination domain.Pagination) ([]domain.ProjectItem, error) {
	r.logger.WithFields(logrus.Fields{
		"filter_name":     filter.Name,
		"filter_status":   filter.Status,
		"filter_priority": filter.Priority,
		"limit":           pagination.Limit,
		"offset":          pagination.Offset,
		"sort":            pagination.Sort,
	}).Debug("Listing project items from database with filters")

	var items []domain.ProjectItem
	db := r.db.WithContext(ctx).Model(&domain.ProjectItem{})

	if filter.ProjectID != nil {
		r.logger.WithFields(logrus.Fields{
			"filter_project_id": filter.ProjectID,
		}).Debug("Applying project_id filter")
		db = db.Where("project_id = ?", filter.ProjectID)
	}

	if filter.Name != "" {
		r.logger.WithFields(logrus.Fields{
			"filter_name": filter.Name,
		}).Debug("Applying name filter")
		db = db.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.Status != "" {
		r.logger.WithFields(logrus.Fields{
			"filter_status": filter.Status,
		}).Debug("Applying status filter")
		db = db.Where("status = ?", filter.Status)
	}

	if filter.Priority != "" {
		r.logger.WithFields(logrus.Fields{
			"filter_priority": filter.Priority,
		}).Debug("Applying priority filter")
		db = db.Where("priority = ?", filter.Priority)
	}

	if filter.AssignedTo != nil {
		r.logger.WithFields(logrus.Fields{
			"filter_assigned_to": filter.AssignedTo,
		}).Debug("Applying assigned_to filter")
		db = db.Where("assigned_to = ?", filter.AssignedTo)
	}

	if filter.DueDateFrom != nil {
		r.logger.WithFields(logrus.Fields{
			"due_date_from": filter.DueDateFrom,
		}).Debug("Applying due_date_from filter")
		db = db.Where("due_date >= ?", *filter.DueDateFrom)
	}

	if filter.DueDateTo != nil {
		r.logger.WithFields(logrus.Fields{
			"due_date_to": filter.DueDateTo,
		}).Debug("Applying due_date_to filter")
		db = db.Where("due_date <= ?", *filter.DueDateTo)
	}

	if filter.EstimatedHoursFrom != nil {
		r.logger.WithFields(logrus.Fields{
			"estimated_hours_from": filter.EstimatedHoursFrom,
		}).Debug("Applying estimated_hours_from filter")
		db = db.Where("estimated_hours >= ?", *filter.EstimatedHoursFrom)
	}

	if filter.EstimatedHoursTo != nil {
		r.logger.WithFields(logrus.Fields{
			"estimated_hours_to": filter.EstimatedHoursTo,
		}).Debug("Applying estimated_hours_to filter")
		db = db.Where("estimated_hours <= ?", *filter.EstimatedHoursTo)
	}

	if filter.ActualHoursFrom != nil {
		r.logger.WithFields(logrus.Fields{
			"actual_hours_from": filter.ActualHoursFrom,
		}).Debug("Applying actual_hours_from filter")
		db = db.Where("actual_hours >= ?", *filter.ActualHoursFrom)
	}

	if filter.ActualHoursTo != nil {
		r.logger.WithFields(logrus.Fields{
			"actual_hours_to": filter.ActualHoursTo,
		}).Debug("Applying actual_hours_to filter")
		db = db.Where("actual_hours <= ?", *filter.ActualHoursTo)
	}

	if filter.CreatedAtFrom != nil {
		r.logger.WithFields(logrus.Fields{
			"created_at_from": filter.CreatedAtFrom,
		}).Debug("Applying created_at_from filter")
		db = db.Where("created_at >= ?", *filter.CreatedAtFrom)
	}

	if filter.CreatedAtTo != nil {
		r.logger.WithFields(logrus.Fields{
			"created_at_to": filter.CreatedAtTo,
		}).Debug("Applying created_at_to filter")
		db = db.Where("created_at <= ?", *filter.CreatedAtTo)
	}

	db = db.Where("deleted_at IS NULL")

	if pagination.Sort != "" {
		r.logger.WithFields(logrus.Fields{
			"sort": pagination.Sort,
		}).Debug("Applying sort")
		db = db.Order(pagination.Sort)
	}

	if pagination.Limit > 0 {
		r.logger.WithFields(logrus.Fields{
			"limit": pagination.Limit,
		}).Debug("Applying limit")
		db = db.Limit(pagination.Limit)
	}

	if pagination.Offset > 0 {
		r.logger.WithFields(logrus.Fields{
			"offset": pagination.Offset,
		}).Debug("Applying offset")
		db = db.Offset(pagination.Offset)
	}

	if err := db.Find(&items).Error; err != nil {
		r.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list project items from database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"count": len(items),
	}).Debug("Project items listed successfully from database")

	return items, nil
}

func (r *PostgresProjectItemRepository) Update(ctx context.Context, item *domain.ProjectItem) error {
	r.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"status":     item.Status,
		"project_id": item.ProjectID,
	}).Debug("Updating project item in database")

	err := r.db.WithContext(ctx).Model(item).Updates(item).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"item_id": item.ID,
		}).Error("Failed to update project item in database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Debug("Project item updated successfully in database")

	return nil
}

func (r *PostgresProjectItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.WithFields(logrus.Fields{
		"item_id": id,
	}).Debug("Soft deleting project item in database")

	err := r.db.WithContext(ctx).Model(&domain.ProjectItem{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"item_id": id,
		}).Error("Failed to delete project item from database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"item_id": id,
	}).Debug("Project item soft deleted successfully in database")

	return nil
}

func (r *PostgresProjectItemRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]domain.ProjectItem, error) {
	r.logger.WithFields(logrus.Fields{
		"project_id": projectID,
	}).Debug("Getting project items by project ID from database")

	var items []domain.ProjectItem
	err := r.db.WithContext(ctx).Where("project_id = ? AND deleted_at IS NULL", projectID).Find(&items).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": projectID,
		}).Error("Failed to get project items by project ID from database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"count":      len(items),
	}).Debug("Project items retrieved successfully by project ID from database")

	return items, nil
}

func (r *PostgresProjectItemRepository) GetByAssignedTo(ctx context.Context, assignedTo uuid.UUID) ([]domain.ProjectItem, error) {
	r.logger.WithFields(logrus.Fields{
		"assigned_to": assignedTo,
	}).Debug("Getting project items by assigned user from database")

	var items []domain.ProjectItem
	err := r.db.WithContext(ctx).Where("assigned_to = ? AND deleted_at IS NULL", assignedTo).Find(&items).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":       err.Error(),
			"assigned_to": assignedTo,
		}).Error("Failed to get project items by assigned user from database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"assigned_to": assignedTo,
		"count":       len(items),
	}).Debug("Project items retrieved successfully by assigned user from database")

	return items, nil
}
