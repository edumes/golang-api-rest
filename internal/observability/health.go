package observability

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HealthStatus struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

type HealthChecker interface {
	Check(ctx context.Context) HealthStatus
}

type DatabaseHealthChecker struct {
	db *gorm.DB
}

func NewDatabaseHealthChecker(db *gorm.DB) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{db: db}
}

func (d *DatabaseHealthChecker) Check(ctx context.Context) HealthStatus {
	start := time.Now()

	sqlDB, err := d.db.DB()
	if err != nil {
		return HealthStatus{
			Status:    "unhealthy",
			Timestamp: time.Now(),
			Details: map[string]interface{}{
				"error": err.Error(),
			},
		}
	}

	err = sqlDB.PingContext(ctx)
	duration := time.Since(start)

	status := HealthStatus{
		Timestamp: time.Now(),
		Details: map[string]interface{}{
			"response_time_ms": duration.Milliseconds(),
		},
	}

	if err != nil {
		status.Status = "unhealthy"
		status.Details["error"] = err.Error()
	} else {
		status.Status = "healthy"
		status.Details["max_open_connections"] = sqlDB.Stats().MaxOpenConnections
		status.Details["open_connections"] = sqlDB.Stats().OpenConnections
		status.Details["in_use_connections"] = sqlDB.Stats().InUse
		status.Details["idle_connections"] = sqlDB.Stats().Idle
	}

	return status
}

type HealthHandler struct {
	dbChecker *DatabaseHealthChecker
	logger    *logrus.Logger
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		dbChecker: NewDatabaseHealthChecker(db),
		logger:    logrus.New(),
	}
}

func (h *HealthHandler) LiveCheck(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"endpoint": "/health/live",
		"ip":       c.ClientIP(),
	}).Debug("Health live check requested")

	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now(),
		"service":   "golang-api-rest",
	})
}

func (h *HealthHandler) ReadyCheck(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"endpoint": "/health/ready",
		"ip":       c.ClientIP(),
	}).Debug("Health ready check requested")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	dbStatus := h.dbChecker.Check(ctx)

	overallStatus := "ready"
	httpStatus := http.StatusOK

	if dbStatus.Status != "healthy" {
		overallStatus = "not_ready"
		httpStatus = http.StatusServiceUnavailable
	}

	response := gin.H{
		"status":    overallStatus,
		"timestamp": time.Now(),
		"service":   "golang-api-rest",
		"checks": gin.H{
			"database": dbStatus,
		},
	}

	c.JSON(httpStatus, response)
}

func (h *HealthHandler) DetailedCheck(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"endpoint": "/health/detailed",
		"ip":       c.ClientIP(),
	}).Debug("Health detailed check requested")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	dbStatus := h.dbChecker.Check(ctx)

	systemInfo := gin.H{
		"go_version": "1.24.4",
		"uptime":     time.Since(time.Now()).String(),
		"memory": gin.H{
			"alloc": "N/A",
			"total": "N/A",
		},
	}

	overallStatus := "healthy"
	httpStatus := http.StatusOK

	if dbStatus.Status != "healthy" {
		overallStatus = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	response := gin.H{
		"status":      overallStatus,
		"timestamp":   time.Now(),
		"service":     "golang-api-rest",
		"version":     "1.0.0",
		"system_info": systemInfo,
		"checks": gin.H{
			"database": dbStatus,
		},
	}

	c.JSON(httpStatus, response)
}

func (h *HealthHandler) RegisterHealthRoutes(router *gin.RouterGroup) {
	health := router.Group("/health")
	{
		// @Summary Health live check
		// @Description Check if the application is alive
		// @Tags health
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Router /health/live [get]
		health.GET("/live", h.LiveCheck)

		// @Summary Health ready check
		// @Description Check if the application is ready to serve requests
		// @Tags health
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 503 {object} map[string]interface{}
		// @Router /health/ready [get]
		health.GET("/ready", h.ReadyCheck)

		// @Summary Health detailed check
		// @Description Get detailed health information
		// @Tags health
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 503 {object} map[string]interface{}
		// @Router /health/detailed [get]
		health.GET("/detailed", h.DetailedCheck)
	}
}
