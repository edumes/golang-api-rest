package api

import (
	"github.com/edumes/golang-api-rest/internal/application"
	"github.com/edumes/golang-api-rest/internal/observability"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	engine *gin.Engine
	logger *logrus.Logger
}

func NewRouter() *Router {
	return &Router{
		engine: gin.New(),
		logger: logrus.New(),
	}
}

func (r *Router) SetupRoutes(userService *application.UserService, productService *application.ProductService, projectService *application.ProjectService, projectItemService *application.ProjectItemService, healthHandler *observability.HealthHandler) {
	r.logger.Info("Setting up application routes")

	r.engine.Use(RequestIDMiddleware())
	r.engine.Use(LoggerMiddleware())
	r.engine.Use(ErrorHandlerMiddleware())
	r.engine.Use(RecoveryMiddleware())
	r.engine.Use(CORSMiddleware())
	r.engine.Use(observability.MetricsMiddleware())

	r.logger.Debug("Middleware configured successfully")

	r.engine.GET("/metrics", observability.MetricsHandler())
	healthHandler.RegisterHealthRoutes(r.engine.Group(""))

	r.engine.GET(SwaggerEndpoint, ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.logger.Debug("Swagger endpoint configured")

	userHandler := NewUserHandler(userService)
	authHandler := NewAuthHandler(userService)
	productHandler := NewProductHandler(productService)
	projectHandler := NewProjectHandler(projectService)
	projectItemHandler := NewProjectItemHandler(projectItemService)

	r.logger.Debug("Handlers created successfully")

	r.setupV1Routes(userHandler, authHandler, productHandler, projectHandler, projectItemHandler)

	r.logger.Info("All routes configured successfully")
}

func (r *Router) setupV1Routes(userHandler *UserHandler, authHandler *AuthHandler, productHandler *ProductHandler, projectHandler *ProjectHandler, projectItemHandler *ProjectItemHandler) {
	r.logger.Info("Setting up v1 API routes")

	v1 := r.engine.Group(APIVersion)

	r.logger.Info("Registering public routes")
	authHandler.RegisterRoutes(v1)

	r.logger.Info("Registering protected routes")
	protected := v1.Group("")
	protected.Use(AuthMiddleware())
	userHandler.RegisterRoutes(protected)
	productHandler.RegisterRoutes(protected)
	projectHandler.RegisterRoutes(protected)
	projectItemHandler.RegisterRoutes(protected)
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
