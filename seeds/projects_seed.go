package seeds

import (
	"context"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/google/uuid"
)

func SeedProjects(repo domain.ProjectRepository) error {
	ctx := context.Background()

	projects := []domain.Project{
		{
			ID:          uuid.New(),
			Name:        "E-commerce Platform",
			Description: "A modern e-commerce platform with payment integration",
			Status:      "active",
			StartDate:   &[]time.Time{time.Now().AddDate(0, -2, 0)}[0],
			EndDate:     &[]time.Time{time.Now().AddDate(0, 4, 0)}[0],
			Budget:      &[]float64{50000.0}[0],
			OwnerID:     uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "Mobile App Development",
			Description: "Cross-platform mobile application for iOS and Android",
			Status:      "active",
			StartDate:   &[]time.Time{time.Now().AddDate(0, -1, 0)}[0],
			EndDate:     &[]time.Time{time.Now().AddDate(0, 5, 0)}[0],
			Budget:      &[]float64{75000.0}[0],
			OwnerID:     uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Name:        "API Documentation",
			Description: "Comprehensive API documentation and testing suite",
			Status:      "completed",
			StartDate:   &[]time.Time{time.Now().AddDate(0, -3, 0)}[0],
			EndDate:     &[]time.Time{time.Now().AddDate(0, -1, 0)}[0],
			Budget:      &[]float64{15000.0}[0],
			OwnerID:     uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, project := range projects {
		if err := repo.Create(ctx, &project); err != nil {
			return err
		}
	}

	return nil
}
