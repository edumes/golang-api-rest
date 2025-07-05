package seeds

import (
	"context"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/google/uuid"
)

func SeedProjectItems(repo domain.ProjectItemRepository, projectRepo domain.ProjectRepository) error {
	ctx := context.Background()

	projects, err := projectRepo.List(ctx, domain.ProjectParams{}, domain.Pagination{Limit: 10})
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		return nil
	}

	projectID := projects[0].ID

	items := []domain.ProjectItem{
		{
			ID:             uuid.New(),
			ProjectID:      projectID,
			Name:           "Database Design",
			Description:    "Design and implement the database schema",
			Status:         "completed",
			Priority:       "high",
			EstimatedHours: &[]float64{16.0}[0],
			ActualHours:    &[]float64{18.0}[0],
			DueDate:        &[]time.Time{time.Now().AddDate(0, -1, 0)}[0],
			AssignedTo:     &[]uuid.UUID{uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")}[0],
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			ProjectID:      projectID,
			Name:           "User Authentication",
			Description:    "Implement user registration and login system",
			Status:         "in_progress",
			Priority:       "high",
			EstimatedHours: &[]float64{24.0}[0],
			ActualHours:    &[]float64{12.0}[0],
			DueDate:        &[]time.Time{time.Now().AddDate(0, 1, 0)}[0],
			AssignedTo:     &[]uuid.UUID{uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")}[0],
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			ProjectID:      projectID,
			Name:           "Payment Integration",
			Description:    "Integrate payment gateway (Stripe/PayPal)",
			Status:         "pending",
			Priority:       "medium",
			EstimatedHours: &[]float64{32.0}[0],
			ActualHours:    nil,
			DueDate:        &[]time.Time{time.Now().AddDate(0, 2, 0)}[0],
			AssignedTo:     &[]uuid.UUID{uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")}[0],
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			ProjectID:      projectID,
			Name:           "Frontend Development",
			Description:    "Build responsive user interface",
			Status:         "pending",
			Priority:       "medium",
			EstimatedHours: &[]float64{40.0}[0],
			ActualHours:    nil,
			DueDate:        &[]time.Time{time.Now().AddDate(0, 3, 0)}[0],
			AssignedTo:     &[]uuid.UUID{uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")}[0],
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		{
			ID:             uuid.New(),
			ProjectID:      projectID,
			Name:           "Testing & QA",
			Description:    "Comprehensive testing and quality assurance",
			Status:         "pending",
			Priority:       "low",
			EstimatedHours: &[]float64{20.0}[0],
			ActualHours:    nil,
			DueDate:        &[]time.Time{time.Now().AddDate(0, 4, 0)}[0],
			AssignedTo:     &[]uuid.UUID{uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")}[0],
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	for _, item := range items {
		if err := repo.Create(ctx, &item); err != nil {
			return err
		}
	}

	return nil
}
