package infrastructure

import (
	"context"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db:     db,
		logger: logrus.New(),
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
	}).Debug("Creating user in database")

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
			"email":   user.Email,
		}).Error("Failed to create user in database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Debug("User created successfully in database")

	return nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	r.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Debug("Getting user by ID from database")

	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": id,
		}).Warn("User not found in database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Debug("User retrieved successfully from database")

	return &user, nil
}

func (r *PostgresUserRepository) List(ctx context.Context, filter domain.Params, pagination domain.Pagination) ([]domain.User, error) {
	r.logger.WithFields(logrus.Fields{
		"filter_name":  filter.Name,
		"filter_email": filter.Email,
		"limit":        pagination.Limit,
		"offset":       pagination.Offset,
		"sort":         pagination.Sort,
	}).Debug("Listing users from database with filters")

	var users []domain.User
	db := r.db.WithContext(ctx).Model(&domain.User{})

	if filter.Name != "" {
		r.logger.WithFields(logrus.Fields{
			"filter_name": filter.Name,
		}).Debug("Applying name filter")
		db = db.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.Email != "" {
		r.logger.WithFields(logrus.Fields{
			"filter_email": filter.Email,
		}).Debug("Applying email filter")
		db = db.Where("email = ?", filter.Email)
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

	if err := db.Find(&users).Error; err != nil {
		r.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list users from database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"count": len(users),
	}).Debug("Users listed successfully from database")

	return users, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.Name,
	}).Debug("Updating user in database")

	err := r.db.WithContext(ctx).Model(user).Updates(user).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		}).Error("Failed to update user in database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Debug("User updated successfully in database")

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Debug("Soft deleting user in database")

	err := r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Failed to delete user from database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Debug("User soft deleted successfully in database")

	return nil
}
