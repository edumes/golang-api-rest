package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/edumes/golang-api-rest/docs"
	"github.com/edumes/golang-api-rest/internal/api"
	"github.com/edumes/golang-api-rest/internal/application"
	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/edumes/golang-api-rest/internal/infrastructure"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Golang API REST
// @version 1.0
// @description API REST in Go with Clean Architecture
// @host localhost:8080
// @BasePath /

func main() {
	logger := infrastructure.GetColoredLogger()

	logger.Info("Starting Golang API REST application")

	logger.Info("Loading configuration")
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Failed to read .env file, using environment variables")
	}
	viper.AutomaticEnv()

	logger.Info("Configuring application logging")
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(logrus.DebugLevel)

	gin.SetMode(gin.ReleaseMode)
	logger.Info("Gin mode set to release")

	logger.Info("Initializing database connection")
	db, err := infrastructure.NewPostgresDB()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to connect to database")
	}

	logger.Info("Running database migrations")
	if err := db.AutoMigrate(&domain.User{}, &domain.Product{}, &domain.Project{}, &domain.ProjectItem{}); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to run database migrations")
	}
	logger.Info("Database migrations completed successfully")

	logger.Info("Initializing repositories and services")
	userRepo := infrastructure.NewPostgresUserRepository(db)
	userService := application.NewUserService(userRepo)

	productRepo := infrastructure.NewPostgresProductRepository(db)
	productService := application.NewProductService(productRepo)

	projectRepo := infrastructure.NewPostgresProjectRepository(db)
	projectService := application.NewProjectService(projectRepo)

	projectItemRepo := infrastructure.NewPostgresProjectItemRepository(db)
	projectItemService := application.NewProjectItemService(projectItemRepo)
	logger.Info("Repositories and services initialized successfully")

	logger.Info("Setting up application router")
	router := api.NewRouter()
	router.SetupRoutes(userService, productService, projectService, projectItemService)
	r := router.GetEngine()
	logger.Info("Router setup completed")

	port := viper.GetString("APP_PORT")
	if port == "" {
		port = "8080"
		logger.Warn("APP_PORT not set, using default port 8080")
	}

	logger.WithFields(logrus.Fields{
		"port": port,
	}).Info("Starting HTTP server")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		logger.Info("HTTP server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("HTTP server failed to start")
		}
	}()

	logger.Info("HTTP server started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutdown signal received, starting shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("Shutting down HTTP server")
	if err := srv.Shutdown(ctx); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Server forced to shutdown")
	}

	logger.Info("Server exited")
}
