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

type UserHandler struct {
	service *application.UserService
	logger  *logrus.Logger
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  infrastructure.GetColoredLogger(),
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	h.logger.Info("Registering user routes")
	r.POST(UsersEndpoint, h.CreateUser)
	r.GET(UsersEndpoint, h.ListUsers)
	r.GET(UserByID, h.GetUser)
	r.PUT(UserByID, h.UpdateUser)
	r.DELETE(UserByID, h.DeleteUser)
}

type createUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body createUserRequest true "User data"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Creating new user")

	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		}).Warn("Invalid request body for user creation")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"email": req.Email,
		"name":  req.Name,
	}).Debug("Processing user creation request")

	user, err := h.service.CreateUser(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.Email,
		}).Error("Failed to create user")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User created successfully")

	c.JSON(StatusCreated, user)
}

// @Summary List users
// @Description Get a list of users with optional filtering and pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Filter by name"
// @Param email query string false "Filter by email"
// @Param limit query int false "Number of items per page (default: 20)"
// @Param offset query int false "Number of items to skip (default: 0)"
// @Param sort query string false "Sort order (default: created_at desc)"
// @Success 200 {array} domain.User
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Listing users")

	filter := domain.Params{
		Name:  c.Query("name"),
		Email: c.Query("email"),
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	pagination := domain.Pagination{
		Limit:  limit,
		Offset: offset,
		Sort:   c.DefaultQuery("sort", "created_at desc"),
	}

	h.logger.WithFields(logrus.Fields{
		"filter_name":  filter.Name,
		"filter_email": filter.Email,
		"limit":        limit,
		"offset":       offset,
		"sort":         pagination.Sort,
	}).Debug("List users with filters and pagination")

	users, err := h.service.ListUsers(c.Request.Context(), filter, pagination)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to list users")
		c.JSON(StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"count": len(users),
	}).Info("Users listed successfully")

	c.JSON(StatusOK, users)
}

// @Summary Get user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid user ID format")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":  c.Request.Method,
		"path":    c.Request.URL.Path,
		"user_id": id,
		"ip":      c.ClientIP(),
	}).Info("Getting user by ID")

	user, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"user_id":   id,
			"client_ip": c.ClientIP(),
		}).Warn("User not found")
		c.JSON(StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User retrieved successfully")

	c.JSON(StatusOK, user)
}

// @Summary Update user
// @Description Update an existing user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body domain.User true "User data"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid user ID format for update")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":  c.Request.Method,
		"path":    c.Request.URL.Path,
		"user_id": id,
		"ip":      c.ClientIP(),
	}).Info("Updating user")

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"user_id":   id,
			"client_ip": c.ClientIP(),
		}).Warn("Invalid request body for user update")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id
	if err := h.service.UpdateUser(c.Request.Context(), &user); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"user_id":   id,
			"client_ip": c.ClientIP(),
		}).Error("Failed to update user")
		c.JSON(StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
	}).Info("User updated successfully")

	c.JSON(StatusOK, user)
}

// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"param_id":  c.Param("id"),
			"client_ip": c.ClientIP(),
		}).Warn("Invalid user ID format for deletion")
		c.JSON(StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"method":  c.Request.Method,
		"path":    c.Request.URL.Path,
		"user_id": id,
		"ip":      c.ClientIP(),
	}).Info("Deleting user")

	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"user_id":   id,
			"client_ip": c.ClientIP(),
		}).Error("Failed to delete user")
		c.JSON(StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("User deleted successfully")

	c.Status(StatusNoContent)
}
