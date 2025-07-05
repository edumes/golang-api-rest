package infrastructure

import (
	"context"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostgresProductRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewPostgresProductRepository(db *gorm.DB) *PostgresProductRepository {
	return &PostgresProductRepository{
		db:     db,
		logger: logrus.New(),
	}
}

func (r *PostgresProductRepository) Create(ctx context.Context, product *domain.Product) error {
	r.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
		"name":       product.Name,
		"price":      product.Price,
		"stock":      product.Stock,
	}).Debug("Creating product in database")

	err := r.db.WithContext(ctx).Create(product).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": product.ID,
			"sku":        product.SKU,
		}).Error("Failed to create product in database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Debug("Product created successfully in database")

	return nil
}

func (r *PostgresProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	r.logger.WithFields(logrus.Fields{
		"product_id": id,
	}).Debug("Getting product by ID from database")

	var product domain.Product
	err := r.db.WithContext(ctx).First(&product, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
		}).Warn("Product not found in database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Debug("Product retrieved successfully from database")

	return &product, nil
}

func (r *PostgresProductRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	r.logger.WithFields(logrus.Fields{
		"sku": sku,
	}).Debug("Getting product by SKU from database")

	var product domain.Product
	err := r.db.WithContext(ctx).First(&product, "sku = ? AND deleted_at IS NULL", sku).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"sku":   sku,
		}).Warn("Product not found by SKU in database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Debug("Product retrieved successfully by SKU from database")

	return &product, nil
}

func (r *PostgresProductRepository) List(ctx context.Context, filter domain.ProductParams, pagination domain.Pagination) ([]domain.Product, error) {
	r.logger.WithFields(logrus.Fields{
		"filter_name":     filter.Name,
		"filter_category": filter.Category,
		"filter_sku":      filter.SKU,
		"limit":           pagination.Limit,
		"offset":          pagination.Offset,
		"sort":            pagination.Sort,
	}).Debug("Listing products from database with filters")

	var products []domain.Product
	db := r.db.WithContext(ctx).Model(&domain.Product{})

	if filter.Name != "" {
		r.logger.WithFields(logrus.Fields{
			"filter_name": filter.Name,
		}).Debug("Applying name filter")
		db = db.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.Category != "" {
		r.logger.WithFields(logrus.Fields{
			"filter_category": filter.Category,
		}).Debug("Applying category filter")
		db = db.Where("category ILIKE ?", "%"+filter.Category+"%")
	}

	if filter.SKU != "" {
		r.logger.WithFields(logrus.Fields{
			"filter_sku": filter.SKU,
		}).Debug("Applying SKU filter")
		db = db.Where("sku ILIKE ?", "%"+filter.SKU+"%")
	}

	if filter.PriceFrom != nil {
		r.logger.WithFields(logrus.Fields{
			"price_from": *filter.PriceFrom,
		}).Debug("Applying price_from filter")
		db = db.Where("price >= ?", *filter.PriceFrom)
	}

	if filter.PriceTo != nil {
		r.logger.WithFields(logrus.Fields{
			"price_to": *filter.PriceTo,
		}).Debug("Applying price_to filter")
		db = db.Where("price <= ?", *filter.PriceTo)
	}

	if filter.StockFrom != nil {
		r.logger.WithFields(logrus.Fields{
			"stock_from": *filter.StockFrom,
		}).Debug("Applying stock_from filter")
		db = db.Where("stock >= ?", *filter.StockFrom)
	}

	if filter.StockTo != nil {
		r.logger.WithFields(logrus.Fields{
			"stock_to": *filter.StockTo,
		}).Debug("Applying stock_to filter")
		db = db.Where("stock <= ?", *filter.StockTo)
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

	if err := db.Find(&products).Error; err != nil {
		r.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list products from database")
		return nil, err
	}

	r.logger.WithFields(logrus.Fields{
		"count": len(products),
	}).Debug("Products listed successfully from database")

	return products, nil
}

func (r *PostgresProductRepository) Update(ctx context.Context, product *domain.Product) error {
	r.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
		"name":       product.Name,
		"price":      product.Price,
		"stock":      product.Stock,
	}).Debug("Updating product in database")

	err := r.db.WithContext(ctx).Model(product).Updates(product).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": product.ID,
		}).Error("Failed to update product in database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Debug("Product updated successfully in database")

	return nil
}

func (r *PostgresProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.WithFields(logrus.Fields{
		"product_id": id,
	}).Debug("Soft deleting product in database")

	err := r.db.WithContext(ctx).Model(&domain.Product{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
		}).Error("Failed to delete product from database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"product_id": id,
	}).Debug("Product soft deleted successfully in database")

	return nil
}

func (r *PostgresProductRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	r.logger.WithFields(logrus.Fields{
		"product_id": id,
		"quantity":   quantity,
	}).Debug("Updating product stock in database")

	err := r.db.WithContext(ctx).Model(&domain.Product{}).Where("id = ?", id).Update("stock", quantity).Error
	if err != nil {
		r.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
		}).Error("Failed to update product stock in database")
		return err
	}

	r.logger.WithFields(logrus.Fields{
		"product_id": id,
		"new_stock":  quantity,
	}).Debug("Product stock updated successfully in database")

	return nil
}
