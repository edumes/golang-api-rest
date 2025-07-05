package api

import (
	"strconv"
	"time"

	"github.com/edumes/golang-api-rest/internal/application"
	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/edumes/golang-api-rest/internal/infrastructure"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ProjectItemHandler struct {
	service *application.ProjectItemService
	logger  *logrus.Logger
}

func NewProjectItemHandler(service *application.ProjectItemService) *ProjectItemHandler {
	return &ProjectItemHandler{
		service: service,
		logger:  infrastructure.GetColoredLogger(),
	}
}

func (h *ProjectItemHandler) RegisterRoutes(r *gin.RouterGroup) {
	h.logger.Info("Registering project item routes")
	r.POST(ProjectItemsEndpoint, h.CreateProjectItem)
	r.GET(ProjectItemsEndpoint, h.ListProjectItems)
	r.GET(ProjectItemByID, h.GetProjectItem)
	r.PUT(ProjectItemByID, h.UpdateProjectItem)
	r.DELETE(ProjectItemByID, h.DeleteProjectItem)
	r.GET(ProjectItemsByProject, h.GetProjectItemsByProject)
}

type createProjectItemRequest struct {
	ProjectID      uuid.UUID  `json:"project_id" binding:"required"`
	Name           string     `json:"name" binding:"required"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	Priority       string     `json:"priority"`
	EstimatedHours *float64   `json:"estimated_hours"`
	ActualHours    *float64   `json:"actual_hours"`
	DueDate        *time.Time `json:"due_date"`
	AssignedTo     *uuid.UUID `json:"assigned_to"`
}

// @Summary Create project item
// @Description Create a new project item
// @Tags project-items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createProjectItemRequest true "Project item data"
// @Success 201 {object} domain.ProjectItem
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /v1/project-items [post]
func (h *ProjectItemHandler) CreateProjectItem(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Creating new project item")

	var req createProjectItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		}).Warn("Invalid request body for project item creation")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"name":       req.Name,
		"status":     req.Status,
		"priority":   req.Priority,
		"project_id": req.ProjectID,
	}).Debug("Processing project item creation request")

	item, err := h.service.CreateProjectItem(c.Request.Context(), req.ProjectID, req.Name, req.Description, req.Status, req.Priority, req.EstimatedHours, req.ActualHours, req.DueDate, req.AssignedTo)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"name":  req.Name,
		}).Error("Failed to create project item")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Info("Project item created successfully")

	c.JSON(StatusCreated, item)
}

// @Summary List project items
// @Description Get a list of project items with optional filtering and pagination
// @Tags project-items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id query string false "Filter by project ID"
// @Param name query string false "Filter by name"
// @Param status query string false "Filter by status"
// @Param priority query string false "Filter by priority"
// @Param assigned_to query string false "Filter by assigned user ID"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Param sort query string false "Sort order (default: created_at desc)"
// @Success 200 {array} domain.ProjectItem
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /v1/project-items [get]
func (h *ProjectItemHandler) ListProjectItems(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Listing project items")

	filter := domain.ProjectItemParams{
		Name:     c.Query("name"),
		Status:   c.Query("status"),
		Priority: c.Query("priority"),
	}

	if projectIDStr := c.Query("project_id"); projectIDStr != "" {
		if projectID, err := uuid.Parse(projectIDStr); err == nil {
			filter.ProjectID = &projectID
		}
	}

	if assignedToStr := c.Query("assigned_to"); assignedToStr != "" {
		if assignedTo, err := uuid.Parse(assignedToStr); err == nil {
			filter.AssignedTo = &assignedTo
		}
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
		"filter_status":   filter.Status,
		"filter_priority": filter.Priority,
		"limit":           limit,
		"offset":          offset,
		"sort":            pagination.Sort,
	}).Debug("List project items with filters and pagination")

	items, err := h.service.ListProjectItems(c.Request.Context(), filter, pagination)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list project items")
		c.JSON(StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"count": len(items),
	}).Info("Project items listed successfully")

	c.JSON(StatusOK, items)
}

// @Summary Get project item by ID
// @Description Get a specific project item by its ID
// @Tags project-items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project item ID"
// @Success 200 {object} domain.ProjectItem
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /v1/project-items/{id} [get]
func (h *ProjectItemHandler) GetProjectItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid project item ID format")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":  c.Request.Method,
		"path":    c.Request.URL.Path,
		"item_id": id,
		"ip":      c.ClientIP(),
	}).Info("Getting project item by ID")

	item, err := h.service.GetProjectItemByID(c.Request.Context(), id)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"item_id":   id,
			"client_ip": c.ClientIP(),
		}).Warn("Project item not found")
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Info("Project item retrieved successfully")

	c.JSON(StatusOK, item)
}

// @Summary Update project item
// @Description Update an existing project item
// @Tags project-items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project item ID"
// @Param request body domain.ProjectItem true "Project item data"
// @Success 200 {object} domain.ProjectItem
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /v1/project-items/{id} [put]
func (h *ProjectItemHandler) UpdateProjectItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid project item ID format")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":  c.Request.Method,
		"path":    c.Request.URL.Path,
		"item_id": id,
		"ip":      c.ClientIP(),
	}).Info("Updating project item")

	var item domain.ProjectItem
	if err := c.ShouldBindJSON(&item); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		}).Warn("Invalid request body for project item update")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.ID = id

	err = h.service.UpdateProjectItem(c.Request.Context(), &item)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"item_id": id,
		}).Error("Failed to update project item")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"item_id":    item.ID,
		"name":       item.Name,
		"project_id": item.ProjectID,
	}).Info("Project item updated successfully")

	c.JSON(StatusOK, item)
}

// @Summary Delete project item
// @Description Delete a project item (soft delete)
// @Tags project-items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project item ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /v1/project-items/{id} [delete]
func (h *ProjectItemHandler) DeleteProjectItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid project item ID format")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":  c.Request.Method,
		"path":    c.Request.URL.Path,
		"item_id": id,
		"ip":      c.ClientIP(),
	}).Info("Deleting project item")

	err = h.service.DeleteProjectItem(c.Request.Context(), id)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"item_id": id,
		}).Error("Failed to delete project item")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"item_id": id,
	}).Info("Project item deleted successfully")

	c.JSON(StatusNoContent, nil)
}

// @Summary Get project items by project ID
// @Description Get all project items for a specific project
// @Tags project-items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "Project ID"
// @Success 200 {array} domain.ProjectItem
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /v1/project-items/project/{projectId} [get]
func (h *ProjectItemHandler) GetProjectItemsByProject(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":            err.Error(),
			"param_project_id": c.Param("projectId"),
			"client_ip":        c.ClientIP(),
		}).Warn("Invalid project ID format")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"project_id": projectID,
		"ip":         c.ClientIP(),
	}).Info("Getting project items by project ID")

	items, err := h.service.GetProjectItemsByProjectID(c.Request.Context(), projectID)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": projectID,
		}).Error("Failed to get project items by project ID")
		c.JSON(StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"project_id": projectID,
		"count":      len(items),
	}).Info("Project items retrieved successfully by project ID")

	c.JSON(StatusOK, items)
}
