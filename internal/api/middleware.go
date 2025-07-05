package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func AuthMiddleware() gin.HandlerFunc {
	logger := logrus.New()

	return func(c *gin.Context) {
		logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
		}).Debug("Processing authentication middleware")

		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			logger.WithFields(logrus.Fields{
				"ip":     c.ClientIP(),
				"path":   c.Request.URL.Path,
				"header": header,
			}).Warn("Missing or invalid Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		secret := viper.GetString("APP_JWT_SECRET")

		logger.WithFields(logrus.Fields{
			"ip":   c.ClientIP(),
			"path": c.Request.URL.Path,
		}).Debug("Parsing JWT token")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"ip":    c.ClientIP(),
				"path":  c.Request.URL.Path,
			}).Warn("Invalid JWT token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := claims["sub"]
			userEmail := claims["email"]

			logger.WithFields(logrus.Fields{
				"user_id":    userID,
				"user_email": userEmail,
				"ip":         c.ClientIP(),
				"path":       c.Request.URL.Path,
			}).Info("User authenticated successfully")

			c.Set("user_id", userID)
			c.Set("user_email", userEmail)
		}

		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	logger := logrus.New()

	return func(c *gin.Context) {
		start := time.Now()

		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("Incoming request")

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		traceID := c.GetHeader("X-Trace-Id")

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
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     status,
			"latency":    latency,
			"trace_id":   traceID,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		if userID, exists := c.Get("user_id"); exists {
			fields["user_id"] = userID
		}
		if userEmail, exists := c.Get("user_email"); exists {
			fields["user_email"] = userEmail
		}

		logger.WithFields(fields).Log(logLevel, "Request completed")
	}
}

func ErrorRecoveryMiddleware() gin.HandlerFunc {
	logger := logrus.New()

	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.WithFields(logrus.Fields{
				"error":      err,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"ip":         c.ClientIP(),
				"user_agent": c.Request.UserAgent(),
			}).Error("Panic recovered")
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
	})
}
