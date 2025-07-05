package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProjectItem struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	ProjectID      uuid.UUID  `json:"project_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	EstimatedHours *float64   `json:"estimated_hours"`
	ActualHours    *float64   `json:"actual_hours"`
	DueDate        *time.Time `json:"due_date"`
	AssignedTo     *uuid.UUID `json:"assigned_to"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" gorm:"index"`
}

type ProjectItemParams struct {
	ProjectID          *uuid.UUID
	Name               string
	Status             string
	Priority           string
	AssignedTo         *uuid.UUID
	DueDateFrom        *time.Time
	DueDateTo          *time.Time
	EstimatedHoursFrom *float64
	EstimatedHoursTo   *float64
	ActualHoursFrom    *float64
	ActualHoursTo      *float64
	CreatedAtFrom      *time.Time
	CreatedAtTo        *time.Time
}

type ProjectItemRepository interface {
	Create(ctx context.Context, item *ProjectItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*ProjectItem, error)
	List(ctx context.Context, filter ProjectItemParams, pagination Pagination) ([]ProjectItem, error)
	Update(ctx context.Context, item *ProjectItem) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]ProjectItem, error)
	GetByAssignedTo(ctx context.Context, assignedTo uuid.UUID) ([]ProjectItem, error)
}
