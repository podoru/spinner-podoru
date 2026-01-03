package deployment

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"

	domainDocker "github.com/podoru/spinner-podoru/internal/domain/docker"
	"github.com/podoru/spinner-podoru/internal/domain/entity"
	"github.com/podoru/spinner-podoru/internal/domain/repository"
	"github.com/podoru/spinner-podoru/internal/infrastructure/config"
)

var (
	ErrServiceNotFound    = errors.New("service not found")
	ErrProjectNotFound    = errors.New("project not found")
	ErrNotTeamMember      = errors.New("not a team member")
	ErrNoImageSpecified   = errors.New("no image specified for deployment")
	ErrServiceNotDeployed = errors.New("service not deployed yet")
	ErrAlreadyDeploying   = errors.New("deployment already in progress")
)

// UseCase handles deployment operations
type UseCase struct {
	serviceRepo      repository.ServiceRepository
	projectRepo      repository.ProjectRepository
	teamMemberRepo   repository.TeamMemberRepository
	deploymentRepo   repository.DeploymentRepository
	domainRepo       repository.DomainRepository
	containerManager domainDocker.ContainerManager
	traefikConfig    *config.TraefikConfig
}

// NewUseCase creates a new deployment use case
func NewUseCase(
	serviceRepo repository.ServiceRepository,
	projectRepo repository.ProjectRepository,
	teamMemberRepo repository.TeamMemberRepository,
	deploymentRepo repository.DeploymentRepository,
	domainRepo repository.DomainRepository,
	containerManager domainDocker.ContainerManager,
	traefikConfig *config.TraefikConfig,
) *UseCase {
	return &UseCase{
		serviceRepo:      serviceRepo,
		projectRepo:      projectRepo,
		teamMemberRepo:   teamMemberRepo,
		deploymentRepo:   deploymentRepo,
		domainRepo:       domainRepo,
		containerManager: containerManager,
		traefikConfig:    traefikConfig,
	}
}

// Deploy deploys a service
func (uc *UseCase) Deploy(ctx context.Context, userID, serviceID uuid.UUID) (*entity.Deployment, error) {
	// 1. Validate access
	service, err := uc.validateAccess(ctx, userID, serviceID)
	if err != nil {
		return nil, err
	}

	// 2. Check if already deploying
	if service.Status == entity.ServiceStatusDeploying {
		return nil, ErrAlreadyDeploying
	}

	// 3. Validate deploy type - only support image for now
	if service.DeployType != entity.DeployTypeImage {
		return nil, fmt.Errorf("deploy type %s not yet supported", service.DeployType)
	}

	if service.Image == nil || *service.Image == "" {
		return nil, ErrNoImageSpecified
	}

	// 4. Create deployment record
	deployment := &entity.Deployment{
		ID:          uuid.New(),
		ServiceID:   serviceID,
		TriggeredBy: &userID,
		Status:      entity.DeploymentStatusPending,
		StartedAt:   time.Now(),
	}
	if err := uc.deploymentRepo.Create(ctx, deployment); err != nil {
		return nil, err
	}

	// 5. Execute deployment in goroutine
	go uc.executeDeployment(context.Background(), service, deployment)

	return deployment, nil
}

func (uc *UseCase) executeDeployment(ctx context.Context, service *entity.Service, deployment *entity.Deployment) {
	var deployErr error

	defer func() {
		now := time.Now()
		deployment.FinishedAt = &now

		if deployErr != nil {
			deployment.Status = entity.DeploymentStatusFailed
			logs := deployErr.Error()
			deployment.Logs = &logs
			uc.serviceRepo.UpdateStatus(ctx, service.ID, entity.ServiceStatusFailed)
		} else {
			deployment.Status = entity.DeploymentStatusSuccess
			uc.serviceRepo.UpdateStatus(ctx, service.ID, entity.ServiceStatusRunning)
		}

		uc.deploymentRepo.Update(ctx, deployment)
	}()

	// Update status to deploying
	deployment.Status = entity.DeploymentStatusDeploying
	uc.deploymentRepo.Update(ctx, deployment)
	uc.serviceRepo.UpdateStatus(ctx, service.ID, entity.ServiceStatusDeploying)

	// Stop and remove existing container if exists
	if service.ContainerID != nil && *service.ContainerID != "" {
		_ = uc.containerManager.StopContainer(ctx, *service.ContainerID, nil)
		_ = uc.containerManager.RemoveContainer(ctx, *service.ContainerID, true)
	}

	// Pull the image
	if err := uc.containerManager.PullImage(ctx, *service.Image); err != nil {
		deployErr = fmt.Errorf("failed to pull image: %w", err)
		return
	}

	// Fetch domains for the service
	domains, err := uc.domainRepo.ListByServiceID(ctx, service.ID)
	if err != nil {
		deployErr = fmt.Errorf("failed to fetch domains: %w", err)
		return
	}

	// Build container config
	containerName := fmt.Sprintf("podoru-%s", service.Slug)

	var memLimit *int64
	if service.MemoryLimit != nil {
		mem := int64(*service.MemoryLimit) * 1024 * 1024 // Convert MB to bytes
		memLimit = &mem
	}

	// Build labels with Traefik configuration
	labels := uc.buildContainerLabels(service, domains)

	// Determine network
	var networkID string
	if uc.traefikConfig != nil && uc.traefikConfig.Enabled && len(domains) > 0 {
		networkID = uc.traefikConfig.Network

		// Validate network exists before deployment
		if err := uc.containerManager.ValidateNetwork(ctx, networkID); err != nil {
			deployErr = fmt.Errorf("traefik network '%s' not found - ensure Traefik is running: %w", networkID, err)
			return
		}
	}

	config := &domainDocker.ContainerConfig{
		Name:          containerName,
		Image:         *service.Image,
		Env:           nil, // TODO: decrypt env vars
		PortMappings:  nil, // TODO: add port mappings
		Volumes:       nil, // TODO: add volumes
		CPULimit:      service.CPULimit,
		MemoryLimit:   memLimit,
		RestartPolicy: service.RestartPolicy,
		Labels:        labels,
		NetworkID:     networkID,
	}

	// Create container
	containerID, err := uc.containerManager.CreateContainer(ctx, config)
	if err != nil {
		deployErr = fmt.Errorf("failed to create container: %w", err)
		return
	}

	// Update container ID in database
	if err := uc.serviceRepo.UpdateContainerID(ctx, service.ID, &containerID); err != nil {
		deployErr = fmt.Errorf("failed to update container ID: %w", err)
		return
	}

	// Start container
	if err := uc.containerManager.StartContainer(ctx, containerID); err != nil {
		deployErr = fmt.Errorf("failed to start container: %w", err)
		return
	}
}

// Start starts a deployed service
func (uc *UseCase) Start(ctx context.Context, userID, serviceID uuid.UUID) error {
	service, err := uc.validateAccess(ctx, userID, serviceID)
	if err != nil {
		return err
	}

	if service.ContainerID == nil || *service.ContainerID == "" {
		return ErrServiceNotDeployed
	}

	if err := uc.containerManager.StartContainer(ctx, *service.ContainerID); err != nil {
		return err
	}

	return uc.serviceRepo.UpdateStatus(ctx, serviceID, entity.ServiceStatusRunning)
}

// Stop stops a running service
func (uc *UseCase) Stop(ctx context.Context, userID, serviceID uuid.UUID) error {
	service, err := uc.validateAccess(ctx, userID, serviceID)
	if err != nil {
		return err
	}

	if service.ContainerID == nil || *service.ContainerID == "" {
		return ErrServiceNotDeployed
	}

	if err := uc.containerManager.StopContainer(ctx, *service.ContainerID, nil); err != nil {
		return err
	}

	return uc.serviceRepo.UpdateStatus(ctx, serviceID, entity.ServiceStatusStopped)
}

// Restart restarts a running service
func (uc *UseCase) Restart(ctx context.Context, userID, serviceID uuid.UUID) error {
	service, err := uc.validateAccess(ctx, userID, serviceID)
	if err != nil {
		return err
	}

	if service.ContainerID == nil || *service.ContainerID == "" {
		return ErrServiceNotDeployed
	}

	return uc.containerManager.RestartContainer(ctx, *service.ContainerID, nil)
}

// GetLogs retrieves logs from a service
func (uc *UseCase) GetLogs(ctx context.Context, userID, serviceID uuid.UUID, tail, since string) (string, error) {
	service, err := uc.validateAccess(ctx, userID, serviceID)
	if err != nil {
		return "", err
	}

	if service.ContainerID == nil || *service.ContainerID == "" {
		return "", ErrServiceNotDeployed
	}

	if tail == "" {
		tail = "100"
	}

	opts := &domainDocker.LogOptions{
		Tail:  tail,
		Since: since,
	}

	reader, err := uc.containerManager.GetLogs(ctx, *service.ContainerID, opts)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// Read logs (with size limit)
	buf := make([]byte, 1024*1024) // 1MB max
	n, err := io.ReadFull(reader, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return "", err
	}

	return string(buf[:n]), nil
}

// buildContainerLabels builds Docker labels including Traefik configuration
func (uc *UseCase) buildContainerLabels(service *entity.Service, domains []entity.Domain) map[string]string {
	labels := map[string]string{
		"podoru.service.id": service.ID.String(),
		"podoru.project.id": service.ProjectID.String(),
		"podoru.managed":    "true",
	}

	// Add Traefik labels if enabled and domains exist
	if uc.traefikConfig == nil || !uc.traefikConfig.Enabled || len(domains) == 0 {
		return labels
	}

	routerName := fmt.Sprintf("podoru-%s", service.Slug)

	// Enable Traefik for this container
	labels["traefik.enable"] = "true"

	// Build Host rules for all domains
	var hostRules []string
	for _, d := range domains {
		hostRules = append(hostRules, fmt.Sprintf("Host(`%s`)", d.Domain))
	}
	hostRule := strings.Join(hostRules, " || ")

	// HTTP router
	labels[fmt.Sprintf("traefik.http.routers.%s.rule", routerName)] = hostRule
	labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", routerName)] = "web"

	// Check if any domain has SSL enabled
	hasSSL := false
	hasAutoSSL := false
	for _, d := range domains {
		if d.SSLEnabled {
			hasSSL = true
			if d.SSLAuto {
				hasAutoSSL = true
				break
			}
		}
	}

	// HTTPS router with SSL
	if hasSSL {
		httpsRouterName := routerName + "-secure"
		labels[fmt.Sprintf("traefik.http.routers.%s.rule", httpsRouterName)] = hostRule
		labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", httpsRouterName)] = "websecure"
		labels[fmt.Sprintf("traefik.http.routers.%s.tls", httpsRouterName)] = "true"

		if hasAutoSSL {
			labels[fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", httpsRouterName)] = "letsencrypt"
		}

		// HTTP to HTTPS redirect
		labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", routerName)] = routerName + "-redirect"
		labels[fmt.Sprintf("traefik.http.middlewares.%s-redirect.redirectscheme.scheme", routerName)] = "https"
		labels[fmt.Sprintf("traefik.http.middlewares.%s-redirect.redirectscheme.permanent", routerName)] = "true"
	}

	// Service configuration - use first port mapping or default to 80
	serviceName := routerName
	labels[fmt.Sprintf("traefik.http.routers.%s.service", routerName)] = serviceName
	if hasSSL {
		labels[fmt.Sprintf("traefik.http.routers.%s-secure.service", routerName)] = serviceName
	}

	// Default to port 80 if no health check path specified
	port := "80"
	if service.HealthCheckPath != nil && *service.HealthCheckPath != "" {
		// Use port 80 as default, can be enhanced to use port mappings
		port = "80"
	}
	labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", serviceName)] = port

	return labels
}

// Destroy stops and removes the container for a service (used before deletion)
func (uc *UseCase) Destroy(ctx context.Context, userID, serviceID uuid.UUID) error {
	service, err := uc.validateAccess(ctx, userID, serviceID)
	if err != nil {
		return err
	}

	// If no container, nothing to destroy
	if service.ContainerID == nil || *service.ContainerID == "" {
		return nil
	}

	// Stop container (ignore errors - it might already be stopped)
	_ = uc.containerManager.StopContainer(ctx, *service.ContainerID, nil)

	// Remove container
	if err := uc.containerManager.RemoveContainer(ctx, *service.ContainerID, true); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	return nil
}

// Helper methods

func (uc *UseCase) validateAccess(ctx context.Context, userID, serviceID uuid.UUID) (*entity.Service, error) {
	service, err := uc.serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrServiceNotFound
	}

	project, err := uc.projectRepo.GetByID(ctx, service.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, project.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	return service, nil
}
