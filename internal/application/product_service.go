package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProductService struct {
	repo   domain.ProductRepository
	logger *logrus.Logger
}

func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{
		repo:   repo,
		logger: logrus.New(),
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, name, description, category, sku string, price float64, stock int) (*domain.Product, error) {
	s.logger.WithFields(logrus.Fields{
		"name":     name,
		"category": category,
		"sku":      sku,
		"price":    price,
		"stock":    stock,
	}).Info("Creating new product")

	if strings.TrimSpace(name) == "" {
		s.logger.WithFields(logrus.Fields{
			"name": name,
		}).Warn("Product name is empty")
		return nil, errors.New("product name is required")
	}

	if strings.TrimSpace(sku) == "" {
		s.logger.WithFields(logrus.Fields{
			"sku": sku,
		}).Warn("Product SKU is empty")
		return nil, errors.New("product SKU is required")
	}

	if price <= 0 {
		s.logger.WithFields(logrus.Fields{
			"price": price,
		}).Warn("Invalid product price")
		return nil, errors.New("product price must be greater than zero")
	}

	if stock < 0 {
		s.logger.WithFields(logrus.Fields{
			"stock": stock,
		}).Warn("Invalid product stock")
		return nil, errors.New("product stock cannot be negative")
	}

	existingProduct, err := s.repo.GetBySKU(ctx, sku)
	if err == nil && existingProduct != nil {
		s.logger.WithFields(logrus.Fields{
			"sku": sku,
		}).Warn("Product SKU already exists")
		return nil, errors.New("product SKU already exists")
	}

	product := &domain.Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		Category:    category,
		SKU:         sku,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Debug("Saving product to repository")

	if err := s.repo.Create(ctx, product); err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": product.ID,
			"sku":        product.SKU,
		}).Error("Failed to create product in repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Info("Product created successfully")

	return product, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	s.logger.WithFields(logrus.Fields{
		"product_id": id,
	}).Debug("Getting product by ID")

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
		}).Warn("Product not found by ID")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Debug("Product retrieved successfully")

	return product, nil
}

func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	s.logger.WithFields(logrus.Fields{
		"sku": sku,
	}).Debug("Getting product by SKU")

	product, err := s.repo.GetBySKU(ctx, sku)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"sku":   sku,
		}).Warn("Product not found by SKU")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Debug("Product retrieved successfully by SKU")

	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context, filter domain.ProductParams, pagination domain.Pagination) ([]domain.Product, error) {
	s.logger.WithFields(logrus.Fields{
		"filter_name":     filter.Name,
		"filter_category": filter.Category,
		"filter_sku":      filter.SKU,
		"limit":           pagination.Limit,
		"offset":          pagination.Offset,
		"sort":            pagination.Sort,
	}).Debug("Listing products with filters")

	products, err := s.repo.List(ctx, filter, pagination)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list products from repository")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"count": len(products),
	}).Info("Products listed successfully")

	return products, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	s.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Info("Updating product")

	if strings.TrimSpace(product.Name) == "" {
		s.logger.WithFields(logrus.Fields{
			"product_id": product.ID,
		}).Warn("Product name is empty")
		return errors.New("product name is required")
	}

	if product.Price <= 0 {
		s.logger.WithFields(logrus.Fields{
			"product_id": product.ID,
			"price":      product.Price,
		}).Warn("Invalid product price")
		return errors.New("product price must be greater than zero")
	}

	if product.Stock < 0 {
		s.logger.WithFields(logrus.Fields{
			"product_id": product.ID,
			"stock":      product.Stock,
		}).Warn("Invalid product stock")
		return errors.New("product stock cannot be negative")
	}

	product.UpdatedAt = time.Now()

	err := s.repo.Update(ctx, product)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": product.ID,
		}).Error("Failed to update product in repository")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Info("Product updated successfully")

	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	s.logger.WithFields(logrus.Fields{
		"product_id": id,
	}).Info("Deleting product")

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
		}).Error("Failed to delete product from repository")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"product_id": id,
	}).Info("Product deleted successfully")

	return nil
}

func (s *ProductService) UpdateProductStock(ctx context.Context, id uuid.UUID, quantity int) error {
	s.logger.WithFields(logrus.Fields{
		"product_id": id,
		"quantity":   quantity,
	}).Info("Updating product stock")

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
		}).Warn("Product not found for stock update")
		return err
	}

	newStock := product.Stock + quantity
	if newStock < 0 {
		s.logger.WithFields(logrus.Fields{
			"product_id":    id,
			"current_stock": product.Stock,
			"quantity":      quantity,
			"new_stock":     newStock,
		}).Warn("Insufficient stock for update")
		return errors.New("insufficient stock")
	}

	err = s.repo.UpdateStock(ctx, id, newStock)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
		}).Error("Failed to update product stock in repository")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"product_id": id,
		"old_stock":  product.Stock,
		"new_stock":  newStock,
	}).Info("Product stock updated successfully")

	return nil
}
