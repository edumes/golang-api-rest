package seeds

import (
	"context"
	"fmt"

	"github.com/edumes/golang-api-rest/internal/infrastructure"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Seeder struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		db:     db,
		logger: logrus.New(),
	}
}

func (s *Seeder) RunAll(ctx context.Context) error {
	s.logger.Info("Starting all seeds...")

	userSeed := NewUserSeed(s.db)
	if err := userSeed.Run(ctx); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to run user seeds")
		return fmt.Errorf("failed to run user seeds: %w", err)
	}

	projectRepo := infrastructure.NewPostgresProjectRepository(s.db)
	if err := SeedProjects(projectRepo); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to run project seeds")
		return fmt.Errorf("failed to run project seeds: %w", err)
	}

	projectItemRepo := infrastructure.NewPostgresProjectItemRepository(s.db)
	if err := SeedProjectItems(projectItemRepo, projectRepo); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to run project item seeds")
		return fmt.Errorf("failed to run project item seeds: %w", err)
	}

	s.logger.Info("All seeds completed successfully")
	return nil
}

func (s *Seeder) RunUsers(ctx context.Context) error {
	s.logger.Info("Starting user seeds...")

	userSeed := NewUserSeed(s.db)
	if err := userSeed.Run(ctx); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to run user seeds")
		return fmt.Errorf("failed to run user seeds: %w", err)
	}

	s.logger.Info("User seeds completed successfully")
	return nil
}

func (s *Seeder) RunProjects(ctx context.Context) error {
	s.logger.Info("Starting project seeds...")

	projectRepo := infrastructure.NewPostgresProjectRepository(s.db)
	if err := SeedProjects(projectRepo); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to run project seeds")
		return fmt.Errorf("failed to run project seeds: %w", err)
	}

	s.logger.Info("Project seeds completed successfully")
	return nil
}

func (s *Seeder) RunProjectItems(ctx context.Context) error {
	s.logger.Info("Starting project item seeds...")

	projectRepo := infrastructure.NewPostgresProjectRepository(s.db)
	projectItemRepo := infrastructure.NewPostgresProjectItemRepository(s.db)
	if err := SeedProjectItems(projectItemRepo, projectRepo); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to run project item seeds")
		return fmt.Errorf("failed to run project item seeds: %w", err)
	}

	s.logger.Info("Project item seeds completed successfully")
	return nil
}
