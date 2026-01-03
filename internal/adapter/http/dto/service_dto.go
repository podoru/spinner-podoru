package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

// ServiceResponse represents service data in API responses
type ServiceResponse struct {
	ID                 uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProjectID          uuid.UUID `json:"project_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name               string    `json:"name" example:"api-server"`
	Slug               string    `json:"slug" example:"api-server"`
	DeployType         string    `json:"deploy_type" example:"dockerfile"`
	Image              *string   `json:"image,omitempty" example:"nginx:latest"`
	DockerfilePath     string    `json:"dockerfile_path" example:"Dockerfile"`
	BuildContext       string    `json:"build_context" example:"."`
	Replicas           int       `json:"replicas" example:"2"`
	CPULimit           *float64  `json:"cpu_limit,omitempty" example:"0.5"`
	MemoryLimit        *int      `json:"memory_limit,omitempty" example:"512"`
	HealthCheckPath    *string   `json:"health_check_path,omitempty" example:"/health"`
	HealthCheckInterval int      `json:"health_check_interval" example:"30"`
	RestartPolicy      string    `json:"restart_policy" example:"unless-stopped"`
	Status             string    `json:"status" example:"running"`
	ContainerID        *string   `json:"container_id,omitempty" example:"abc123def456"`
	SwarmServiceID     *string   `json:"swarm_service_id,omitempty" example:"svc_abc123"`
	CreatedAt          time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt          time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// CreateServiceRequest represents the service creation payload
type CreateServiceRequest struct {
	Name               string   `json:"name" validate:"required,min=2,max=100" example:"api-server"`
	Slug               string   `json:"slug" validate:"required,slug,min=2,max=100" example:"api-server"`
	DeployType         string   `json:"deploy_type" validate:"required,oneof=image dockerfile compose" example:"dockerfile"`
	Image              *string  `json:"image,omitempty" example:"nginx:latest"`
	DockerfilePath     *string  `json:"dockerfile_path,omitempty" example:"Dockerfile"`
	BuildContext       *string  `json:"build_context,omitempty" example:"."`
	ComposeFile        *string  `json:"compose_file,omitempty" example:"docker-compose.yml"`
	EnvVars            []EnvVar `json:"env_vars,omitempty"`
	Replicas           *int     `json:"replicas,omitempty" example:"1"`
	CPULimit           *float64 `json:"cpu_limit,omitempty" example:"0.5"`
	MemoryLimit        *int     `json:"memory_limit,omitempty" example:"512"`
	HealthCheckPath    *string  `json:"health_check_path,omitempty" example:"/health"`
	HealthCheckInterval *int    `json:"health_check_interval,omitempty" example:"30"`
	RestartPolicy      *string  `json:"restart_policy,omitempty" example:"unless-stopped"`
}

// UpdateServiceRequest represents the service update payload
type UpdateServiceRequest struct {
	Name               *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100" example:"updated-api-server"`
	Image              *string  `json:"image,omitempty" example:"nginx:alpine"`
	DockerfilePath     *string  `json:"dockerfile_path,omitempty" example:"Dockerfile.prod"`
	BuildContext       *string  `json:"build_context,omitempty" example:"./api"`
	EnvVars            []EnvVar `json:"env_vars,omitempty"`
	Replicas           *int     `json:"replicas,omitempty" example:"3"`
	CPULimit           *float64 `json:"cpu_limit,omitempty" example:"1.0"`
	MemoryLimit        *int     `json:"memory_limit,omitempty" example:"1024"`
	HealthCheckPath    *string  `json:"health_check_path,omitempty" example:"/api/health"`
	HealthCheckInterval *int    `json:"health_check_interval,omitempty" example:"60"`
	RestartPolicy      *string  `json:"restart_policy,omitempty" example:"always"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `json:"key" validate:"required" example:"DATABASE_URL"`
	Value string `json:"value" validate:"required" example:"postgresql://localhost:5432/mydb"`
}

// ScaleServiceRequest represents the service scaling payload
type ScaleServiceRequest struct {
	Replicas int `json:"replicas" validate:"required,min=0,max=100" example:"5"`
}

// ServiceLogsResponse represents service logs
type ServiceLogsResponse struct {
	ServiceID uuid.UUID `json:"service_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Logs      string    `json:"logs" example:"[2024-01-15 10:30:00] Server started on port 8080\n[2024-01-15 10:30:01] Connected to database"`
	Timestamp time.Time `json:"timestamp" example:"2024-01-15T10:30:00Z"`
}

// DeploymentResponse represents deployment information
type DeploymentResponse struct {
	ID            uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceID     uuid.UUID  `json:"service_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	TriggeredBy   *uuid.UUID `json:"triggered_by,omitempty" example:"550e8400-e29b-41d4-a716-446655440002"`
	CommitSHA     *string    `json:"commit_sha,omitempty" example:"abc123def456"`
	CommitMessage *string    `json:"commit_message,omitempty" example:"Fix bug in API endpoint"`
	Status        string     `json:"status" example:"success"`
	Logs          *string    `json:"logs,omitempty"`
	StartedAt     time.Time  `json:"started_at" example:"2024-01-15T10:30:00Z"`
	FinishedAt    *time.Time `json:"finished_at,omitempty" example:"2024-01-15T10:32:00Z"`
}

func ToServiceResponse(service *entity.Service) ServiceResponse {
	return ServiceResponse{
		ID:                  service.ID,
		ProjectID:           service.ProjectID,
		Name:                service.Name,
		Slug:                service.Slug,
		DeployType:          string(service.DeployType),
		Image:               service.Image,
		DockerfilePath:      service.DockerfilePath,
		BuildContext:        service.BuildContext,
		Replicas:            service.Replicas,
		CPULimit:            service.CPULimit,
		MemoryLimit:         service.MemoryLimit,
		HealthCheckPath:     service.HealthCheckPath,
		HealthCheckInterval: service.HealthCheckInterval,
		RestartPolicy:       string(service.RestartPolicy),
		Status:              string(service.Status),
		ContainerID:         service.ContainerID,
		SwarmServiceID:      service.SwarmServiceID,
		CreatedAt:           service.CreatedAt,
		UpdatedAt:           service.UpdatedAt,
	}
}

func ToServicesResponse(services []entity.Service) []ServiceResponse {
	responses := make([]ServiceResponse, len(services))
	for i, s := range services {
		responses[i] = ToServiceResponse(&s)
	}
	return responses
}
