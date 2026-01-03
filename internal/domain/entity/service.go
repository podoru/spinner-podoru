package entity

import (
	"time"

	"github.com/google/uuid"
)

type DeployType string

const (
	DeployTypeImage      DeployType = "image"
	DeployTypeDockerfile DeployType = "dockerfile"
	DeployTypeCompose    DeployType = "compose"
)

type ServiceStatus string

const (
	ServiceStatusStopped   ServiceStatus = "stopped"
	ServiceStatusRunning   ServiceStatus = "running"
	ServiceStatusDeploying ServiceStatus = "deploying"
	ServiceStatusFailed    ServiceStatus = "failed"
)

type RestartPolicy string

const (
	RestartPolicyNo            RestartPolicy = "no"
	RestartPolicyAlways        RestartPolicy = "always"
	RestartPolicyOnFailure     RestartPolicy = "on-failure"
	RestartPolicyUnlessStopped RestartPolicy = "unless-stopped"
)

type Service struct {
	ID                  uuid.UUID      `json:"id"`
	ProjectID           uuid.UUID      `json:"project_id"`
	Name                string         `json:"name"`
	Slug                string         `json:"slug"`
	DeployType          DeployType     `json:"deploy_type"`
	Image               *string        `json:"image,omitempty"`
	DockerfilePath      string         `json:"dockerfile_path"`
	BuildContext        string         `json:"build_context"`
	ComposeFile         *string        `json:"compose_file,omitempty"`
	EnvVarsEncrypted    []byte         `json:"-"`
	Replicas            int            `json:"replicas"`
	CPULimit            *float64       `json:"cpu_limit,omitempty"`
	MemoryLimit         *int           `json:"memory_limit,omitempty"`
	HealthCheckPath     *string        `json:"health_check_path,omitempty"`
	HealthCheckInterval int            `json:"health_check_interval"`
	RestartPolicy       RestartPolicy  `json:"restart_policy"`
	Status              ServiceStatus  `json:"status"`
	ContainerID         *string        `json:"container_id,omitempty"`
	SwarmServiceID      *string        `json:"swarm_service_id,omitempty"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

type ServiceCreate struct {
	Name                string         `json:"name" validate:"required,min=2,max=100"`
	Slug                string         `json:"slug" validate:"required,min=2,max=100,slug"`
	DeployType          DeployType     `json:"deploy_type" validate:"required,oneof=image dockerfile compose"`
	Image               *string        `json:"image,omitempty" validate:"omitempty,max=500"`
	DockerfilePath      *string        `json:"dockerfile_path,omitempty" validate:"omitempty,max=500"`
	BuildContext        *string        `json:"build_context,omitempty" validate:"omitempty,max=500"`
	ComposeFile         *string        `json:"compose_file,omitempty" validate:"omitempty,max=500"`
	EnvVars             map[string]string `json:"env_vars,omitempty"`
	Replicas            *int           `json:"replicas,omitempty" validate:"omitempty,min=1,max=100"`
	CPULimit            *float64       `json:"cpu_limit,omitempty" validate:"omitempty,min=0.1,max=128"`
	MemoryLimit         *int           `json:"memory_limit,omitempty" validate:"omitempty,min=32,max=524288"`
	HealthCheckPath     *string        `json:"health_check_path,omitempty" validate:"omitempty,max=255"`
	HealthCheckInterval *int           `json:"health_check_interval,omitempty" validate:"omitempty,min=5,max=300"`
	RestartPolicy       *RestartPolicy `json:"restart_policy,omitempty" validate:"omitempty,oneof=no always on-failure unless-stopped"`
}

type ServiceUpdate struct {
	Name                *string        `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Image               *string        `json:"image,omitempty" validate:"omitempty,max=500"`
	DockerfilePath      *string        `json:"dockerfile_path,omitempty" validate:"omitempty,max=500"`
	BuildContext        *string        `json:"build_context,omitempty" validate:"omitempty,max=500"`
	ComposeFile         *string        `json:"compose_file,omitempty" validate:"omitempty,max=500"`
	EnvVars             map[string]string `json:"env_vars,omitempty"`
	Replicas            *int           `json:"replicas,omitempty" validate:"omitempty,min=1,max=100"`
	CPULimit            *float64       `json:"cpu_limit,omitempty" validate:"omitempty,min=0.1,max=128"`
	MemoryLimit         *int           `json:"memory_limit,omitempty" validate:"omitempty,min=32,max=524288"`
	HealthCheckPath     *string        `json:"health_check_path,omitempty" validate:"omitempty,max=255"`
	HealthCheckInterval *int           `json:"health_check_interval,omitempty" validate:"omitempty,min=5,max=300"`
	RestartPolicy       *RestartPolicy `json:"restart_policy,omitempty" validate:"omitempty,oneof=no always on-failure unless-stopped"`
}

type ServiceScale struct {
	Replicas int `json:"replicas" validate:"required,min=0,max=100"`
}

type ServiceWithDetails struct {
	Service
	EnvVars      map[string]string `json:"env_vars,omitempty"`
	Domains      []Domain          `json:"domains,omitempty"`
	PortMappings []PortMapping     `json:"port_mappings,omitempty"`
	Volumes      []Volume          `json:"volumes,omitempty"`
}
