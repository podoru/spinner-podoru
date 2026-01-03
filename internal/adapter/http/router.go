package http

import (
	"github.com/gin-gonic/gin"

	"github.com/podoru/spinner-podoru/internal/adapter/http/handler"
	"github.com/podoru/spinner-podoru/internal/adapter/http/middleware"
	"github.com/podoru/spinner-podoru/internal/infrastructure/logger"
)

type Router struct {
	engine         *gin.Engine
	log            *logger.Logger
	authMiddleware *middleware.AuthMiddleware
	authHandler    *handler.AuthHandler
	userHandler    *handler.UserHandler
	teamHandler    *handler.TeamHandler
	projectHandler *handler.ProjectHandler
	serviceHandler *handler.ServiceHandler
	docsHandler    *handler.DocsHandler
}

type RouterConfig struct {
	Logger         *logger.Logger
	AuthMiddleware *middleware.AuthMiddleware
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	TeamHandler    *handler.TeamHandler
	ProjectHandler *handler.ProjectHandler
	ServiceHandler *handler.ServiceHandler
	DocsHandler    *handler.DocsHandler
}

func NewRouter(cfg *RouterConfig) *Router {
	return &Router{
		log:            cfg.Logger,
		authMiddleware: cfg.AuthMiddleware,
		authHandler:    cfg.AuthHandler,
		userHandler:    cfg.UserHandler,
		teamHandler:    cfg.TeamHandler,
		projectHandler: cfg.ProjectHandler,
		serviceHandler: cfg.ServiceHandler,
		docsHandler:    cfg.DocsHandler,
	}
}

func (r *Router) Setup(mode string) *gin.Engine {
	if mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	r.engine = engine

	engine.Use(middleware.RequestID())
	engine.Use(middleware.CORS(nil))
	engine.Use(middleware.Logger(r.log))
	engine.Use(middleware.Recovery(r.log))

	r.setupRoutes()

	return engine
}

func (r *Router) setupRoutes() {
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	api := r.engine.Group("/api/v1")

	r.setupDocsRoutes(api)
	r.setupAuthRoutes(api)
	r.setupUserRoutes(api)
	r.setupTeamRoutes(api)
	r.setupProjectRoutes(api)
	r.setupServiceRoutes(api)
}

func (r *Router) setupDocsRoutes(api *gin.RouterGroup) {
	if r.docsHandler == nil {
		return
	}

	docs := api.Group("/docs")
	{
		docs.GET("", r.docsHandler.Scalar)
		docs.GET("/openapi.json", r.docsHandler.OpenAPISpec)
	}
}

func (r *Router) setupAuthRoutes(api *gin.RouterGroup) {
	if r.authHandler == nil {
		return
	}

	auth := api.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/refresh", r.authHandler.Refresh)
		auth.POST("/logout", r.authHandler.Logout)
	}
}

func (r *Router) setupUserRoutes(api *gin.RouterGroup) {
	if r.userHandler == nil {
		return
	}

	users := api.Group("/users")
	users.Use(r.authMiddleware.RequireAuth())
	{
		users.GET("/me", r.userHandler.GetMe)
		users.PUT("/me", r.userHandler.UpdateMe)
		users.PUT("/me/password", r.userHandler.UpdatePassword)
	}
}

func (r *Router) setupTeamRoutes(api *gin.RouterGroup) {
	if r.teamHandler == nil {
		return
	}

	teams := api.Group("/teams")
	teams.Use(r.authMiddleware.RequireAuth())
	{
		teams.GET("", r.teamHandler.List)
		teams.POST("", r.teamHandler.Create)
		teams.GET("/:teamId", r.teamHandler.Get)
		teams.PUT("/:teamId", r.teamHandler.Update)
		teams.DELETE("/:teamId", r.teamHandler.Delete)
		teams.GET("/:teamId/members", r.teamHandler.ListMembers)
		teams.POST("/:teamId/members", r.teamHandler.AddMember)
		teams.PUT("/:teamId/members/:userId", r.teamHandler.UpdateMember)
		teams.DELETE("/:teamId/members/:userId", r.teamHandler.RemoveMember)

		teams.GET("/:teamId/projects", r.projectHandler.ListByTeam)
		teams.POST("/:teamId/projects", r.projectHandler.Create)
	}
}

func (r *Router) setupProjectRoutes(api *gin.RouterGroup) {
	if r.projectHandler == nil {
		return
	}

	projects := api.Group("/projects")
	projects.Use(r.authMiddleware.RequireAuth())
	{
		projects.GET("/:projectId", r.projectHandler.Get)
		projects.PUT("/:projectId", r.projectHandler.Update)
		projects.DELETE("/:projectId", r.projectHandler.Delete)
		projects.POST("/:projectId/deploy", r.projectHandler.Deploy)

		projects.GET("/:projectId/services", r.serviceHandler.ListByProject)
		projects.POST("/:projectId/services", r.serviceHandler.Create)
	}
}

func (r *Router) setupServiceRoutes(api *gin.RouterGroup) {
	if r.serviceHandler == nil {
		return
	}

	services := api.Group("/services")
	services.Use(r.authMiddleware.RequireAuth())
	{
		services.GET("/:serviceId", r.serviceHandler.Get)
		services.PUT("/:serviceId", r.serviceHandler.Update)
		services.DELETE("/:serviceId", r.serviceHandler.Delete)
		services.POST("/:serviceId/deploy", r.serviceHandler.Deploy)
		services.POST("/:serviceId/start", r.serviceHandler.Start)
		services.POST("/:serviceId/stop", r.serviceHandler.Stop)
		services.POST("/:serviceId/restart", r.serviceHandler.Restart)
		services.POST("/:serviceId/scale", r.serviceHandler.Scale)
		services.GET("/:serviceId/logs", r.serviceHandler.Logs)

		// Domain routes
		services.GET("/:serviceId/domains", r.serviceHandler.ListDomains)
		services.POST("/:serviceId/domains", r.serviceHandler.AddDomain)
		services.DELETE("/:serviceId/domains/:domainId", r.serviceHandler.DeleteDomain)
	}
}
