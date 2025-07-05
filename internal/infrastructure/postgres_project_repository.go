package infrastructure

import (
	"context"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostgresProjectRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewPostgresProjectRepository(db *gorm.DB) *PostgresProjectRepository {
	return &PostgresProjectRepository{
		db:     db,
		logger: logrus.New(),
	}
}

func (r *PostgresProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	r.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
		"owner_id":   project.OwnerID,
	}).Debug("Creating project in database")

	err := r.db.WithContext(ctx).Create(project).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": project.ID,
			"name":       project.Name,
		}).Error("Failed to create project in database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
	}).Debug("Project created successfully in database")

	return nil
}

func (r *PostgresProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	r.logger.WithFields(logrus.Fields{
		"project_id": id,
	}).Debug("Getting project by ID from database")

	var project domain.Project
	err := r.db.WithContext(ctx).First(&project, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": id,
		}).Warn("Project not found in database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
	}).Debug("Project retrieved successfully from database")

	return &project, nil
}

func (r *PostgresProjectRepository) List(ctx context.Context, filter domain.ProjectParams, pagination domain.Pagination) ([]domain.Project, error) {
	r.logger.WithFields(logrus.Fields{
		"filter_name":   filter.Name,
		"filter_status": filter.Status,
		"limit":         pagination.Limit,
		"offset":        pagination.Offset,
		"sort":          pagination.Sort,
	}).Debug("Listing projects from database with filters")

	var projects []domain.Project
	db := r.db.WithContext(ctx).Model(&domain.Project{})

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

	if filter.OwnerID != nil {
		r.logger.WithFields(logrus.Fields{
			"filter_owner_id": filter.OwnerID,
		}).Debug("Applying owner_id filter")
		db = db.Where("owner_id = ?", filter.OwnerID)
	}

	if filter.StartDateFrom != nil {
		r.logger.WithFields(logrus.Fields{
			"start_date_from": filter.StartDateFrom,
		}).Debug("Applying start_date_from filter")
		db = db.Where("start_date >= ?", *filter.StartDateFrom)
	}

	if filter.StartDateTo != nil {
		r.logger.WithFields(logrus.Fields{
			"start_date_to": filter.StartDateTo,
		}).Debug("Applying start_date_to filter")
		db = db.Where("start_date <= ?", *filter.StartDateTo)
	}

	if filter.EndDateFrom != nil {
		r.logger.WithFields(logrus.Fields{
			"end_date_from": filter.EndDateFrom,
		}).Debug("Applying end_date_from filter")
		db = db.Where("end_date >= ?", *filter.EndDateFrom)
	}

	if filter.EndDateTo != nil {
		r.logger.WithFields(logrus.Fields{
			"end_date_to": filter.EndDateTo,
		}).Debug("Applying end_date_to filter")
		db = db.Where("end_date <= ?", *filter.EndDateTo)
	}

	if filter.BudgetFrom != nil {
		r.logger.WithFields(logrus.Fields{
			"budget_from": filter.BudgetFrom,
		}).Debug("Applying budget_from filter")
		db = db.Where("budget >= ?", *filter.BudgetFrom)
	}

	if filter.BudgetTo != nil {
		r.logger.WithFields(logrus.Fields{
			"budget_to": filter.BudgetTo,
		}).Debug("Applying budget_to filter")
		db = db.Where("budget <= ?", *filter.BudgetTo)
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

	if err := db.Find(&projects).Error; err != nil {
		r.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list projects from database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"count": len(projects),
	}).Debug("Projects listed successfully from database")

	return projects, nil
}

func (r *PostgresProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	r.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
		"status":     project.Status,
	}).Debug("Updating project in database")

	err := r.db.WithContext(ctx).Model(project).Updates(project).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": project.ID,
		}).Error("Failed to update project in database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
	}).Debug("Project updated successfully in database")

	return nil
}

func (r *PostgresProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.WithFields(logrus.Fields{
		"project_id": id,
	}).Debug("Soft deleting project in database")

	err := r.db.WithContext(ctx).Model(&domain.Project{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": id,
		}).Error("Failed to delete project from database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"project_id": id,
	}).Debug("Project soft deleted successfully in database")

	return nil
}

func (r *PostgresProjectRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]domain.Project, error) {
	r.logger.WithFields(logrus.Fields{
		"owner_id": ownerID,
	}).Debug("Getting projects by owner ID from database")

	var projects []domain.Project
	err := r.db.WithContext(ctx).Where("owner_id = ? AND deleted_at IS NULL", ownerID).Find(&projects).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":    err.Error(),
			"owner_id": ownerID,
		}).Error("Failed to get projects by owner ID from database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"owner_id": ownerID,
		"count":    len(projects),
	}).Debug("Projects retrieved successfully by owner ID from database")

	return projects, nil
}
