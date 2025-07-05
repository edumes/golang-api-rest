package seeds

import (
	"context"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/edumes/golang-api-rest/internal/infrastructure"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeed struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserSeed(db *gorm.DB) *UserSeed {
	return &UserSeed{
		db:     db,
		logger: logrus.New(),
	}
}

func (s *UserSeed) Run(ctx context.Context) error {
	s.logger.Info("Starting user seeds...")

	users := []domain.User{
		{
			ID:           uuid.New(),
			Name:         "Admin User",
			Email:        "admin@example.com",
			PasswordHash: s.hashPassword("admin123"),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Name:         "John Doe",
			Email:        "john.doe@example.com",
			PasswordHash: s.hashPassword("password123"),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Name:         "Jane Smith",
			Email:        "jane.smith@example.com",
			PasswordHash: s.hashPassword("password123"),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Name:         "Bob Johnson",
			Email:        "bob.johnson@example.com",
			PasswordHash: s.hashPassword("password123"),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           uuid.New(),
			Name:         "Alice Brown",
			Email:        "alice.brown@example.com",
			PasswordHash: s.hashPassword("password123"),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	repository := infrastructure.NewPostgresUserRepository(s.db)

	for _, user := range users {
		existingUser, err := repository.GetByID(ctx, user.ID)
		if err == nil && existingUser != nil {
			s.logger.WithFields(logrus.Fields{
				"user_id": user.ID,
				"email":   user.Email,
			}).Info("User already exists, skipping...")
			continue
		}

		err = repository.Create(ctx, &user)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"error":   err.Error(),
				"user_id": user.ID,
				"email":   user.Email,
			}).Error("Failed to create user seed")
			return err
		}

		s.logger.WithFields(logrus.Fields{
			"user_id": user.ID,
			"email":   user.Email,
			"name":    user.Name,
		}).Info("User seed created successfully")
	}

	s.logger.Info("User seeds completed successfully")
	return nil
}

func (s *UserSeed) hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to hash password")
		return ""
	}
	return string(hashedPassword)
}
