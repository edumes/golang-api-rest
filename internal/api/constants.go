package api

// Route constants
const (
	// API Version
	APIVersion = "/v1"

	// Health check endpoints
	HealthLive  = "/health/live"
	HealthReady = "/health/ready"

	// Auth endpoints
	AuthLogin = "/auth/login"

	// User endpoints
	UsersEndpoint = "/users"
	UserByID      = "/users/:id"

	// Product endpoints
	ProductsEndpoint     = "/products"
	ProductByID          = "/products/:id"
	ProductStockEndpoint = "/products/:id/stock"
	ProductBySKUEndpoint = "/products/sku/:sku"

	// Project endpoints
	ProjectsEndpoint = "/projects"
	ProjectByID      = "/projects/:id"

	// Project Item endpoints
	ProjectItemsEndpoint  = "/project-items"
	ProjectItemByID       = "/project-items/:id"
	ProjectItemsByProject = "/project-items/project/:projectId"

	// Swagger documentation
	SwaggerEndpoint = "/swagger/*any"
)

// HTTP Status codes
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNoContent           = 204
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusNotFound            = 404
	StatusInternalServerError = 500
)
