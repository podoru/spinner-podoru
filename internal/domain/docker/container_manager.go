package docker

import (
	"context"
	"io"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

// ContainerConfig holds configuration for creating a container
type ContainerConfig struct {
	Name          string
	Image         string
	Env           []string
	PortMappings  []entity.PortMapping
	Volumes       []entity.Volume
	CPULimit      *float64
	MemoryLimit   *int64 // in bytes
	RestartPolicy entity.RestartPolicy
	Labels        map[string]string
	NetworkID     string
}

// ContainerInfo holds information about a container
type ContainerInfo struct {
	ID     string
	Status string
	State  string
}

// LogOptions for retrieving container logs
type LogOptions struct {
	Tail   string
	Since  string
	Follow bool
}

// ContainerManager interface for container operations
type ContainerManager interface {
	// Image operations
	PullImage(ctx context.Context, imageName string) error

	// Container operations
	CreateContainer(ctx context.Context, config *ContainerConfig) (string, error)
	StartContainer(ctx context.Context, containerID string) error
	StopContainer(ctx context.Context, containerID string, timeout *int) error
	RestartContainer(ctx context.Context, containerID string, timeout *int) error
	RemoveContainer(ctx context.Context, containerID string, force bool) error
	InspectContainer(ctx context.Context, containerID string) (*ContainerInfo, error)

	// Logs
	GetLogs(ctx context.Context, containerID string, opts *LogOptions) (io.ReadCloser, error)
}

// SwarmServiceConfig holds configuration for creating a swarm service
type SwarmServiceConfig struct {
	Name          string
	Image         string
	Replicas      uint64
	Env           []string
	PortMappings  []entity.PortMapping
	Volumes       []entity.Volume
	CPULimit      *int64 // in nanocores
	MemoryLimit   *int64 // in bytes
	RestartPolicy entity.RestartPolicy
	Labels        map[string]string
	Networks      []string
}

// SwarmManager interface for swarm operations
type SwarmManager interface {
	// Swarm mode check
	IsSwarmMode(ctx context.Context) (bool, error)

	// Service operations
	CreateService(ctx context.Context, config *SwarmServiceConfig) (string, error)
	UpdateService(ctx context.Context, serviceID string, config *SwarmServiceConfig) error
	RemoveService(ctx context.Context, serviceID string) error
	ScaleService(ctx context.Context, serviceID string, replicas uint64) error

	// Logs
	GetServiceLogs(ctx context.Context, serviceID string, opts *LogOptions) (io.ReadCloser, error)
}
