package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/edumes/golang-api-rest/docs"
	"github.com/edumes/golang-api-rest/internal/api"
	"github.com/edumes/golang-api-rest/internal/application"
	"github.com/edumes/golang-api-rest/internal/config"
	"github.com/edumes/golang-api-rest/internal/domain"
	"github.com/edumes/golang-api-rest/internal/infrastructure"
	"github.com/edumes/golang-api-rest/internal/observability"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title Golang API REST
// @version 1.0
// @description API REST in Go with Clean Architecture
// @host localhost:8080
// @BasePath /

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger := infrastructure.GetColoredLogger()
	logger.Info("Starting Golang API REST application")

	logger.Info("Loading configuration")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to load configuration")
	}

	logger.Info("Configuring application logging")
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	switch cfg.Logging.Level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	logger.Info("Gin mode configured")

	logger.Info("Initializing database connection")
	db, err := infrastructure.NewPostgresDBWithConfig(cfg)
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

	logger.Info("Setting up observability")
	observability.SetupMetrics()
	healthHandler := observability.NewHealthHandler(db)
	logger.Info("Observability setup completed")

	logger.Info("Setting up application router")
	router := api.NewRouter()
	router.SetupRoutes(userService, productService, projectService, projectItemService, healthHandler)
	r := router.GetEngine()
	logger.Info("Router setup completed")

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	logger.WithFields(logrus.Fields{
		"port":          cfg.Server.Port,
		"read_timeout":  cfg.Server.ReadTimeout,
		"write_timeout": cfg.Server.WriteTimeout,
		"idle_timeout":  cfg.Server.IdleTimeout,
	}).Info("Starting HTTP server")

	go func() {
		logger.Info("HTTP server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("HTTP server failed to start")
		}
	}()

	logger.Info("HTTP server started successfully")

	<-ctx.Done()

	logger.Info("Shutdown signal received, starting graceful shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Info("Shutting down HTTP server")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Server forced to shutdown")
	}

	logger.Info("Server exited gracefully")
}
