package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	Category    string     `json:"category"`
	SKU         string     `json:"sku" gorm:"uniqueIndex"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

type ProductParams struct {
	Name          string
	Category      string
	SKU           string
	PriceFrom     *float64
	PriceTo       *float64
	StockFrom     *int
	StockTo       *int
	CreatedAtFrom *time.Time
	CreatedAtTo   *time.Time
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetBySKU(ctx context.Context, sku string) (*Product, error)
	List(ctx context.Context, filter ProductParams, pagination Pagination) ([]Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
}
