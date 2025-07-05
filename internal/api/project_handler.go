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

type ProjectHandler struct {
	service *application.ProjectService
	logger  *logrus.Logger
}

func NewProjectHandler(service *application.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		service: service,
		logger:  infrastructure.GetColoredLogger(),
	}
}

func (h *ProjectHandler) RegisterRoutes(r *gin.RouterGroup) {
	h.logger.Info("Registering project routes")
	r.POST(ProjectsEndpoint, h.CreateProject)
	r.GET(ProjectsEndpoint, h.ListProjects)
	r.GET(ProjectByID, h.GetProject)
	r.PUT(ProjectByID, h.UpdateProject)
	r.DELETE(ProjectByID, h.DeleteProject)
}

type createProjectRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Budget      *float64   `json:"budget"`
	OwnerID     uuid.UUID  `json:"owner_id" binding:"required"`
}

// @Summary Create project
// @Description Create a new project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createProjectRequest true "Project data"
// @Success 201 {object} domain.Project
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /v1/projects [post]
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Creating new project")

	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		}).Warn("Invalid request body for project creation")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"name":     req.Name,
		"status":   req.Status,
		"owner_id": req.OwnerID,
	}).Debug("Processing project creation request")

	project, err := h.service.CreateProject(c.Request.Context(), req.Name, req.Description, req.Status, req.StartDate, req.EndDate, req.Budget, req.OwnerID)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"name":  req.Name,
		}).Error("Failed to create project")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
		"owner_id":   project.OwnerID,
	}).Info("Project created successfully")

	c.JSON(StatusCreated, project)
}

// @Summary List projects
// @Description Get a list of projects with optional filtering and pagination
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Filter by name"
// @Param status query string false "Filter by status"
// @Param owner_id query string false "Filter by owner ID"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Param sort query string false "Sort order (default: created_at desc)"
// @Success 200 {array} domain.Project
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /v1/projects [get]
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Listing projects")

	filter := domain.ProjectParams{
		Name:   c.Query("name"),
		Status: c.Query("status"),
	}

	if ownerIDStr := c.Query("owner_id"); ownerIDStr != "" {
		if ownerID, err := uuid.Parse(ownerIDStr); err == nil {
			filter.OwnerID = &ownerID
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
		"filter_name":   filter.Name,
		"filter_status": filter.Status,
		"limit":         limit,
		"offset":        offset,
		"sort":          pagination.Sort,
	}).Debug("List projects with filters and pagination")

	projects, err := h.service.ListProjects(c.Request.Context(), filter, pagination)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list projects")
		c.JSON(StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"count": len(projects),
	}).Info("Projects listed successfully")

	c.JSON(StatusOK, projects)
}

// @Summary Get project by ID
// @Description Get a specific project by its ID
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} domain.Project
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /v1/projects/{id} [get]
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid project ID format")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"project_id": id,
		"ip":         c.ClientIP(),
	}).Info("Getting project by ID")

	project, err := h.service.GetProjectByID(c.Request.Context(), id)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": id,
			"client_ip":  c.ClientIP(),
		}).Warn("Project not found")
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
		"owner_id":   project.OwnerID,
	}).Info("Project retrieved successfully")

	c.JSON(StatusOK, project)
}

// @Summary Update project
// @Description Update an existing project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body domain.Project true "Project data"
// @Success 200 {object} domain.Project
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /v1/projects/{id} [put]
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid project ID format")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"project_id": id,
		"ip":         c.ClientIP(),
	}).Info("Updating project")

	var project domain.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		}).Warn("Invalid request body for project update")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ID = id

	err = h.service.UpdateProject(c.Request.Context(), &project)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": id,
		}).Error("Failed to update project")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"project_id": project.ID,
		"name":       project.Name,
	}).Info("Project updated successfully")

	c.JSON(StatusOK, project)
}

// @Summary Delete project
// @Description Delete a project (soft delete)
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /v1/projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid project ID format")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"project_id": id,
		"ip":         c.ClientIP(),
	}).Info("Deleting project")

	err = h.service.DeleteProject(c.Request.Context(), id)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":      err.Error(),
			"project_id": id,
		}).Error("Failed to delete project")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"project_id": id,
	}).Info("Project deleted successfully")

	c.JSON(StatusNoContent, nil)
}
