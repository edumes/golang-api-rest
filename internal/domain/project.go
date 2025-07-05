package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Budget      *float64   `json:"budget"`
	OwnerID     uuid.UUID  `json:"owner_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

type ProjectParams struct {
	Name          string
	Status        string
	OwnerID       *uuid.UUID
	StartDateFrom *time.Time
	StartDateTo   *time.Time
	EndDateFrom   *time.Time
	EndDateTo     *time.Time
	BudgetFrom    *float64
	BudgetTo      *float64
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
}

type ProjectRepository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	List(ctx context.Context, filter ProjectParams, pagination Pagination) ([]Project, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]Project, error)
}
