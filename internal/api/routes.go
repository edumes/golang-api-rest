package api

import (
	"github.com/edumes/golang-api-rest/internal/application"
	"github.com/gin-contrib/cors"
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

func (r *Router) SetupRoutes(userService *application.UserService, productService *application.ProductService, projectService *application.ProjectService, projectItemService *application.ProjectItemService) {
	r.logger.Info("Setting up application routes")

	r.engine.Use(gin.Recovery())
	r.engine.Use(cors.Default())
	r.engine.Use(LoggingMiddleware())
	r.engine.Use(ErrorRecoveryMiddleware())

	r.logger.Debug("Middleware configured successfully")

	r.engine.GET(SwaggerEndpoint, ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.logger.Debug("Swagger endpoint configured")

	r.setupHealthRoutes()
	r.logger.Debug("Health routes configured")

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

func (r *Router) setupHealthRoutes() {
	r.logger.Debug("Setting up health check routes")

	health := r.engine.Group("/health")
	{
		// @Summary Health live check
		// @Description Check if the application is alive
		// @Tags health
		// @Produce json
		// @Success 200 "OK"
		// @Router /health/live [get]
		health.GET("/live", func(c *gin.Context) {
			r.logger.Debug("Health live check requested")
			c.Status(StatusOK)
		})

		// @Summary Health ready check
		// @Description Check if the application is ready to serve requests
		// @Tags health
		// @Produce json
		// @Success 200 "OK"
		// @Router /health/ready [get]
		health.GET("/ready", func(c *gin.Context) {
			r.logger.Debug("Health ready check requested")
			c.Status(StatusOK)
		})
	}
}

func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
