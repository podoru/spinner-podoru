// Podoru API
//
//	@title						Podoru API
//	@version					1.0.0
//	@description				Podoru is a Docker container management platform similar to EasyPanel.
//	@description				It provides APIs for managing containers, services, networks, and Docker Swarm clusters.
//
//	@contact.name				Podoru Team
//	@contact.url				https://github.com/podoru/spinner-podoru
//
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//
//	@host						localhost:8080
//	@BasePath					/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpAdapter "github.com/podoru/spinner-podoru/internal/adapter/http"
	"github.com/podoru/spinner-podoru/internal/adapter/http/handler"
	"github.com/podoru/spinner-podoru/internal/adapter/http/middleware"
	"github.com/podoru/spinner-podoru/internal/adapter/repository/postgres"
	"github.com/podoru/spinner-podoru/internal/infrastructure/config"
	"github.com/podoru/spinner-podoru/internal/infrastructure/database"
	"github.com/podoru/spinner-podoru/internal/infrastructure/docker"
	"github.com/podoru/spinner-podoru/internal/infrastructure/logger"
	"github.com/podoru/spinner-podoru/internal/usecase/auth"
	"github.com/podoru/spinner-podoru/internal/usecase/deployment"
	"github.com/podoru/spinner-podoru/internal/usecase/project"
	"github.com/podoru/spinner-podoru/internal/usecase/service"
	"github.com/podoru/spinner-podoru/internal/usecase/team"
	"github.com/podoru/spinner-podoru/internal/usecase/user"
	"github.com/podoru/spinner-podoru/pkg/crypto"
	"github.com/podoru/spinner-podoru/pkg/validator"
)

const (
	version = "0.1.0"
)

func main() {
	ctx := context.Background()

	configPath := getEnvOrDefault("CONFIG_PATH", "configs/config.yaml")
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(&cfg.Logger)
	if err != nil {
		fmt.Printf("Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting Podoru",
		"version", version,
		"env", cfg.App.Env,
	)

	db, err := database.New(ctx, &cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Info("Connected to database")

	migrationsPath := getEnvOrDefault("MIGRATIONS_PATH", "migrations")
	migrator, err := database.NewMigrator(&cfg.Database, migrationsPath)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Info("Database migrations completed")

	encryptor, err := crypto.NewEncryptor(cfg.Encryption.Key)
	if err != nil {
		log.Fatalf("Failed to create encryptor: %v", err)
	}

	v, err := validator.New()
	if err != nil {
		log.Fatalf("Failed to create validator: %v", err)
	}

	// Initialize Docker client
	dockerClient, err := docker.NewClient(&cfg.Docker)
	if err != nil {
		log.Warnf("Failed to create Docker client: %v", err)
	} else {
		if err := dockerClient.Ping(ctx); err != nil {
			log.Warnf("Docker daemon not available: %v", err)
		} else {
			log.Info("Connected to Docker daemon")
		}
		defer dockerClient.Close()
	}

	containerManager := docker.NewContainerManager(dockerClient)

	userRepo := postgres.NewUserRepository(db.Pool)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(db.Pool)
	teamRepo := postgres.NewTeamRepository(db.Pool)
	teamMemberRepo := postgres.NewTeamMemberRepository(db.Pool)
	projectRepo := postgres.NewProjectRepository(db.Pool)
	serviceRepo := postgres.NewServiceRepository(db.Pool)
	deploymentRepo := postgres.NewDeploymentRepository(db.Pool)
	domainRepo := postgres.NewDomainRepository(db.Pool)

	authUseCase := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, &cfg.JWT, &cfg.App)
	userUseCase := user.NewUseCase(userRepo)
	teamUseCase := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)
	projectUseCase := project.NewUseCase(projectRepo, teamMemberRepo, encryptor)
	serviceUseCase := service.NewUseCase(serviceRepo, projectRepo, teamMemberRepo, domainRepo, encryptor)
	deploymentUseCase := deployment.NewUseCase(serviceRepo, projectRepo, teamMemberRepo, deploymentRepo, domainRepo, containerManager, &cfg.Traefik)

	authMiddleware := middleware.NewAuthMiddleware(authUseCase)
	authHandler := handler.NewAuthHandler(authUseCase, v)
	userHandler := handler.NewUserHandler(userUseCase, v)
	teamHandler := handler.NewTeamHandler(teamUseCase, v)
	projectHandler := handler.NewProjectHandler(projectUseCase, v)
	serviceHandler := handler.NewServiceHandler(serviceUseCase, deploymentUseCase, v)
	docsHandler := handler.NewDocsHandler()

	router := httpAdapter.NewRouter(&httpAdapter.RouterConfig{
		Logger:         log,
		AuthMiddleware: authMiddleware,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		TeamHandler:    teamHandler,
		ProjectHandler: projectHandler,
		ServiceHandler: serviceHandler,
		DocsHandler:    docsHandler,
	})

	engine := router.Setup(cfg.App.Env)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Infof("Server listening on port %d", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited properly")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
