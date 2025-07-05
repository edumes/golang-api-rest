package api

import (
	"strconv"

	"github.com/edumes/golang-api-rest/internal/application"
	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/edumes/golang-api-rest/internal/infrastructure"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	service *application.ProductService
	logger  *logrus.Logger
}

func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  infrastructure.GetColoredLogger(),
	}
}

func (h *ProductHandler) RegisterRoutes(r *gin.RouterGroup) {
	h.logger.Info("Registering product routes")
	r.POST(ProductsEndpoint, h.CreateProduct)
	r.GET(ProductsEndpoint, h.ListProducts)
	r.GET(ProductByID, h.GetProduct)
	r.PUT(ProductByID, h.UpdateProduct)
	r.DELETE(ProductByID, h.DeleteProduct)
	r.PATCH(ProductStockEndpoint, h.UpdateProductStock)
	r.GET(ProductBySKUEndpoint, h.GetProductBySKU)
}

type createProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	Category    string  `json:"category"`
	SKU         string  `json:"sku" binding:"required"`
}

type updateProductStockRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createProductRequest true "Product data"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /v1/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Creating new product")

	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		}).Warn("Invalid request body for product creation")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"name":     req.Name,
		"sku":      req.SKU,
		"price":    req.Price,
		"stock":    req.Stock,
		"category": req.Category,
	}).Debug("Processing product creation request")

	product, err := h.service.CreateProduct(c.Request.Context(), req.Name, req.Description, req.Category, req.SKU, req.Price, req.Stock)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"sku":   req.SKU,
		}).Error("Failed to create product")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Info("Product created successfully")

	c.JSON(StatusCreated, product)
}

// @Summary List products
// @Description Get a list of products with optional filtering and pagination
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Filter by name"
// @Param category query string false "Filter by category"
// @Param sku query string false "Filter by SKU"
// @Param price_from query number false "Minimum price filter"
// @Param price_to query number false "Maximum price filter"
// @Param stock_from query integer false "Minimum stock filter"
// @Param stock_to query integer false "Maximum stock filter"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Param sort query string false "Sort order (default: created_at desc)"
// @Success 200 {array} domain.Product
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /v1/products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Listing products")

	var priceFrom, priceTo *float64
	if priceFromStr := c.Query("price_from"); priceFromStr != "" {
		if val, err := strconv.ParseFloat(priceFromStr, 64); err == nil {
			priceFrom = &val
		}
	}
	if priceToStr := c.Query("price_to"); priceToStr != "" {
		if val, err := strconv.ParseFloat(priceToStr, 64); err == nil {
			priceTo = &val
		}
	}

	var stockFrom, stockTo *int
	if stockFromStr := c.Query("stock_from"); stockFromStr != "" {
		if val, err := strconv.Atoi(stockFromStr); err == nil {
			stockFrom = &val
		}
	}
	if stockToStr := c.Query("stock_to"); stockToStr != "" {
		if val, err := strconv.Atoi(stockToStr); err == nil {
			stockTo = &val
		}
	}

	filter := domain.ProductParams{
		Name:      c.Query("name"),
		Category:  c.Query("category"),
		SKU:       c.Query("sku"),
		PriceFrom: priceFrom,
		PriceTo:   priceTo,
		StockFrom: stockFrom,
		StockTo:   stockTo,
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	pagination := domain.Pagination{
		Limit:  limit,
		Offset: offset,
		Sort:   c.DefaultQuery("sort", "created_at desc"),
	}

	h.logger.WithFields(logrus.Fields{
		"filter_name":     filter.Name,
		"filter_category": filter.Category,
		"filter_sku":      filter.SKU,
		"limit":           limit,
		"offset":          offset,
		"sort":            pagination.Sort,
	}).Debug("🔍 List products with filters and pagination")

	products, err := h.service.ListProducts(c.Request.Context(), filter, pagination)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list products")
		c.JSON(StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"count": len(products),
	}).Info("Products listed successfully")

	c.JSON(StatusOK, products)
}

// @Summary Get product by ID
// @Description Get a specific product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /v1/products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid product ID format")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"product_id": id,
		"ip":         c.ClientIP(),
	}).Info("Getting product by ID")

	product, err := h.service.GetProductByID(c.Request.Context(), id)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
			"client_ip":  c.ClientIP(),
		}).Warn("Product not found")
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Info("Product retrieved successfully")

	c.JSON(StatusOK, product)
}

// @Summary Get product by SKU
// @Description Get a specific product by its SKU
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param sku path string true "Product SKU"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /v1/products/sku/{sku} [get]
func (h *ProductHandler) GetProductBySKU(c *gin.Context) {
	sku := c.Param("sku")
	if sku == "" {
		h.logger.WithFields(logrus.Fields{
			"client_ip": c.ClientIP(),
		}).Warn("Empty SKU parameter")
		c.JSON(StatusBadRequest, gin.H{"error": "sku parameter is required"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"sku":    sku,
		"ip":     c.ClientIP(),
	}).Info("Getting product by SKU")

	product, err := h.service.GetProductBySKU(c.Request.Context(), sku)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"sku":       sku,
			"client_ip": c.ClientIP(),
		}).Warn("Product not found by SKU")
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Info("Product retrieved successfully by SKU")

	c.JSON(StatusOK, product)
}

// @Summary Update product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param product body domain.Product true "Product data"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid product ID format for update")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"product_id": id,
		"ip":         c.ClientIP(),
	}).Info("Updating product")

	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
			"client_ip":  c.ClientIP(),
		}).Warn("Invalid request body for product update")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.ID = id
	if err := h.service.UpdateProduct(c.Request.Context(), &product); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
			"client_ip":  c.ClientIP(),
		}).Error("Failed to update product")
		c.JSON(StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"product_id": product.ID,
		"sku":        product.SKU,
	}).Info("Product updated successfully")

	c.JSON(StatusOK, product)
}

// @Summary Delete product
// @Description Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid product ID format for deletion")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"product_id": id,
		"ip":         c.ClientIP(),
	}).Info("Deleting product")

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
			"client_ip":  c.ClientIP(),
		}).Error("Failed to delete product")
		c.JSON(StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"product_id": id,
	}).Info("Product deleted successfully")

	c.JSON(StatusNoContent, nil)
}

// @Summary Update product stock
// @Description Update the stock quantity of a product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body updateProductStockRequest true "Stock update data"
// @Success 200 "OK"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /v1/products/{id}/stock [patch]
func (h *ProductHandler) UpdateProductStock(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid product ID format for stock update")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"product_id": id,
		"ip":         c.ClientIP(),
	}).Info("Updating product stock")

	var req updateProductStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
			"client_ip":  c.ClientIP(),
		}).Warn("Invalid request body for stock update")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProductStock(c.Request.Context(), id, req.Quantity); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"product_id": id,
			"quantity":   req.Quantity,
			"client_ip":  c.ClientIP(),
		}).Error("Failed to update product stock")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"product_id": id,
		"quantity":   req.Quantity,
	}).Info("Product stock updated successfully")

	c.JSON(StatusOK, gin.H{"message": "Product stock updated successfully"})
}
