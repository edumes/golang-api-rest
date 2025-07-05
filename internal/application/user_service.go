package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   domain.UserRepository
	logger *logrus.Logger
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{
		repo:   repo,
		logger: logrus.New(),
	}
}

func (s *UserService) CreateUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	s.logger.WithFields(logrus.Fields{
		"email": email,
		"name":  name,
	}).Info("Creating new user")

	if !strings.Contains(email, "@") {
		s.logger.WithFields(logrus.Fields{
			"email": email,
		}).Warn("Invalid email format")
		return nil, errors.New("invalid email")
	}

	if len(password) < 6 {
		s.logger.WithFields(logrus.Fields{
			"password_length": len(password),
		}).Warn("Password too short")
		return nil, errors.New("password too short")
	}

	s.logger.Debug("Generating password hash")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to hash password")
		return nil, err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Debug("Saving user to repository")

	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
			"email":   user.Email,
		}).Error("Failed to create user in repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User created successfully")

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	s.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Debug("Getting user by ID")

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": id,
		}).Warn("User not found by ID")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Debug("User retrieved successfully")

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, filter domain.Params, pagination domain.Pagination) ([]domain.User, error) {
	s.logger.WithFields(logrus.Fields{
		"filter_name":  filter.Name,
		"filter_email": filter.Email,
		"limit":        pagination.Limit,
		"offset":       pagination.Offset,
		"sort":         pagination.Sort,
	}).Debug("Listing users with filters")

	users, err := s.repo.List(ctx, filter, pagination)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list users from repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"count": len(users),
	}).Info("Users listed successfully")

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	s.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("Updating user")

	user.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, user)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": user.ID,
		}).Error("Failed to update user in repository")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User updated successfully")

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	s.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("Deleting user")

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"user_id": id,
		}).Error("Failed to delete user from repository")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User deleted successfully")

	return nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	s.logger.WithFields(logrus.Fields{
		"email": email,
	}).Debug("Getting user by email")

	users, err := s.repo.List(ctx, domain.Params{Email: email}, domain.Pagination{Limit: 1})
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": email,
		}).Error("Failed to get user by email from repository")
		return nil, errors.New("user not found")
	}

	if len(users) == 0 {
		s.logger.WithFields(logrus.Fields{
			"email": email,
		}).Warn("User not found by email")
		return nil, errors.New("user not found")
	}

	user := &users[0]
	s.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Debug("User found by email")

	return user, nil
}

func (s *UserService) CheckPassword(user *domain.User, password string) bool {
	s.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Debug("Checking password")

	isValid := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil

	if isValid {
		s.logger.WithFields(logrus.Fields{
			"user_id": user.ID,
			"email":   user.Email,
		}).Debug("Password check successful")
	} else {
		s.logger.WithFields(logrus.Fields{
			"user_id": user.ID,
			"email":   user.Email,
		}).Warn("Password check failed")
	}

	return isValid
}
