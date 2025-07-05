package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	RequestIDKey = "request_id"
	UserIDKey    = "user_id"
	UserEmailKey = "user_email"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	logger := logrus.New()

	return func(c *gin.Context) {
		start := time.Now()
		requestID := c.GetString(RequestIDKey)

		logger.WithFields(logrus.Fields{
			"request_id":   requestID,
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"ip":           c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"content_type": c.GetHeader("Content-Type"),
		}).Info("Incoming request")

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		bodySize := c.Writer.Size()

		var logLevel logrus.Level
		switch {
		case status >= 500:
			logLevel = logrus.ErrorLevel
		case status >= 400:
			logLevel = logrus.WarnLevel
		default:
			logLevel = logrus.InfoLevel
		}

		fields := logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     status,
			"latency":    latency,
			"body_size":  bodySize,
			"ip":         c.ClientIP(),
		}

		if userID, exists := c.Get(UserIDKey); exists {
			fields["user_id"] = userID
		}
		if userEmail, exists := c.Get(UserEmailKey); exists {
			fields["user_email"] = userEmail
		}

		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		logger.WithFields(fields).Log(logLevel, "Request completed")
	}
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := err.(*domain.AppError); ok {
				c.JSON(appErr.Code, gin.H{
					"error":   appErr.Message,
					"code":    appErr.Code,
					"details": appErr.Details,
				})
				return
			}

			if validationErrors, ok := err.(domain.ValidationErrors); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Validation failed",
					"code":    http.StatusBadRequest,
					"details": validationErrors,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"code":  http.StatusInternalServerError,
			})
		}
	}
}

func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func ValidateRequest(validate *validator.Validate, obj interface{}) error {
	if err := validate.Struct(obj); err != nil {
		var validationErrors domain.ValidationErrors

		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, domain.ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   err.Param(),
				Message: getValidationMessage(err.Tag(), err.Field()),
			})
		}

		return validationErrors
	}

	return nil
}

func getValidationMessage(tag, field string) string {
	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + tag + " characters"
	case "max":
		return field + " must be at most " + tag + " characters"
	case "url":
		return field + " must be a valid URL"
	case "numeric":
		return field + " must be numeric"
	case "alpha":
		return field + " must contain only letters"
	case "alphanum":
		return field + " must contain only letters and numbers"
	default:
		return field + " is invalid"
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := viper.GetStringSlice("CORS_ALLOWED_ORIGINS")
		if len(allowedOrigins) == 0 {
			allowedOrigins = []string{"*"}
		}

		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	logger := logrus.New()

	return func(c *gin.Context) {
		requestID := c.GetString(RequestIDKey)

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
		}).Debug("Processing authentication middleware")

		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"ip":         c.ClientIP(),
				"path":       c.Request.URL.Path,
				"header":     header,
			}).Warn("Missing or invalid Authorization header")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid token",
				"code":  http.StatusUnauthorized,
			})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		secret := viper.GetString("APP_JWT_SECRET")

		logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"ip":         c.ClientIP(),
			"path":       c.Request.URL.Path,
		}).Debug("Parsing JWT token")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			logger.WithFields(logrus.Fields{
				"error":      err.Error(),
				"request_id": requestID,
				"ip":         c.ClientIP(),
				"path":       c.Request.URL.Path,
			}).Warn("Invalid JWT token")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
				"code":  http.StatusUnauthorized,
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := claims["sub"]
			userEmail := claims["email"]

			logger.WithFields(logrus.Fields{
				"user_id":    userID,
				"user_email": userEmail,
				"request_id": requestID,
				"ip":         c.ClientIP(),
				"path":       c.Request.URL.Path,
			}).Info("User authenticated successfully")

			c.Set(UserIDKey, userID)
			c.Set(UserEmailKey, userEmail)
		}

		c.Next()
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	logger := logrus.New()

	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := c.GetString(RequestIDKey)

		logger.WithFields(logrus.Fields{
			"error":      recovered,
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Error("Panic recovered")

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  http.StatusInternalServerError,
		})
	})
}
