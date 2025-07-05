package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Name         string     `json:"name"`
	Email        string     `json:"email" gorm:"uniqueIndex"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`
}

type Params struct {
	Name          string
	Email         string
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
}

type Pagination struct {
	Limit  int
	Offset int
	Sort   string
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	List(ctx context.Context, filter Params, pagination Pagination) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
