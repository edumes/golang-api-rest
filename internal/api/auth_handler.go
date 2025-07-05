package api

import (
	"time"

	"github.com/edumes/golang-api-rest/internal/application"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type AuthHandler struct {
	service *application.UserService
	logger  *logrus.Logger
}

func NewAuthHandler(service *application.UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logrus.New(),
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	h.logger.Info("Registering auth routes")
	r.POST(AuthLogin, h.Login)
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login credentials"
// @Success 200 {object} loginResponse
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	h.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	}).Info("Login attempt")

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"ip":    c.ClientIP(),
		}).Warn("Invalid login request body")
		c.JSON(StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"email": req.Email,
		"ip":    c.ClientIP(),
	}).Debug("Processing login request")

	user, err := h.service.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"email": req.Email,
			"ip":    c.ClientIP(),
		}).Warn("Login failed - user not found")
		c.JSON(StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if !h.service.CheckPassword(user, req.Password) {
		h.logger.WithFields(logrus.Fields{
			"user_id": user.ID,
			"email":   req.Email,
			"ip":      c.ClientIP(),
		}).Warn("Login failed - invalid password")
		c.JSON(StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
		"ip":      c.ClientIP(),
	}).Info("User authenticated successfully")

	secret := viper.GetString("APP_JWT_SECRET")
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"error":     err.Error(),
			"user_id":   user.ID,
			"client_ip": c.ClientIP(),
		}).Error("Failed to generate JWT token")
		c.JSON(StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id": user.ID,
		"email":   user.Email,
		"ip":      c.ClientIP(),
	}).Info("JWT token generated successfully")

	c.JSON(StatusOK, loginResponse{Token: tokenStr})
}
