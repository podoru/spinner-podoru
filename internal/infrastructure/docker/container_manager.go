package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"

	domainDocker "github.com/podoru/spinner-podoru/internal/domain/docker"
	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

// ContainerManagerImpl implements the ContainerManager interface
type ContainerManagerImpl struct {
	client *Client
}

// NewContainerManager creates a new ContainerManager
func NewContainerManager(client *Client) *ContainerManagerImpl {
	return &ContainerManagerImpl{client: client}
}

// PullImage pulls a Docker image
func (m *ContainerManagerImpl) PullImage(ctx context.Context, imageName string) error {
	reader, err := m.client.cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", imageName, err)
	}
	defer reader.Close()

	// Consume the reader to complete the pull
	_, err = io.Copy(io.Discard, reader)
	return err
}

// CreateContainer creates a new container
func (m *ContainerManagerImpl) CreateContainer(ctx context.Context, cfg *domainDocker.ContainerConfig) (string, error) {
	// Build port bindings
	exposedPorts, portBindings := buildPortBindings(cfg.PortMappings)

	// Build mounts
	mounts := buildMounts(cfg.Volumes)

	// Build resource limits
	resources := buildResources(cfg.CPULimit, cfg.MemoryLimit)

	// Build restart policy
	restartPolicy := buildRestartPolicy(cfg.RestartPolicy)

	containerConfig := &container.Config{
		Image:        cfg.Image,
		Env:          cfg.Env,
		ExposedPorts: exposedPorts,
		Labels:       cfg.Labels,
	}

	hostConfig := &container.HostConfig{
		PortBindings:  portBindings,
		Mounts:        mounts,
		Resources:     resources,
		RestartPolicy: restartPolicy,
	}

	networkConfig := &network.NetworkingConfig{}
	if cfg.NetworkID != "" {
		networkConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			cfg.NetworkID: {},
		}
	}

	resp, err := m.client.cli.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, cfg.Name)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	return resp.ID, nil
}

// StartContainer starts a container
func (m *ContainerManagerImpl) StartContainer(ctx context.Context, containerID string) error {
	return m.client.cli.ContainerStart(ctx, containerID, container.StartOptions{})
}

// StopContainer stops a container
func (m *ContainerManagerImpl) StopContainer(ctx context.Context, containerID string, timeout *int) error {
	return m.client.cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: timeout})
}

// RestartContainer restarts a container
func (m *ContainerManagerImpl) RestartContainer(ctx context.Context, containerID string, timeout *int) error {
	return m.client.cli.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: timeout})
}

// RemoveContainer removes a container
func (m *ContainerManagerImpl) RemoveContainer(ctx context.Context, containerID string, force bool) error {
	return m.client.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: force})
}

// InspectContainer inspects a container
func (m *ContainerManagerImpl) InspectContainer(ctx context.Context, containerID string) (*domainDocker.ContainerInfo, error) {
	info, err := m.client.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	return &domainDocker.ContainerInfo{
		ID:     info.ID,
		Status: info.State.Status,
		State:  info.State.Status,
	}, nil
}

// ValidateNetwork checks if a network exists
func (m *ContainerManagerImpl) ValidateNetwork(ctx context.Context, networkName string) error {
	_, err := m.client.cli.NetworkInspect(ctx, networkName, network.InspectOptions{})
	return err
}

// GetLogs retrieves container logs
func (m *ContainerManagerImpl) GetLogs(ctx context.Context, containerID string, opts *domainDocker.LogOptions) (io.ReadCloser, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
	}

	if opts != nil {
		if opts.Tail != "" {
			options.Tail = opts.Tail
		}
		if opts.Since != "" {
			options.Since = opts.Since
		}
		options.Follow = opts.Follow
	}

	return m.client.cli.ContainerLogs(ctx, containerID, options)
}

// Helper functions

func buildPortBindings(portMappings []entity.PortMapping) (nat.PortSet, nat.PortMap) {
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	for _, pm := range portMappings {
		port := nat.Port(fmt.Sprintf("%d/%s", pm.ContainerPort, pm.Protocol))
		exposedPorts[port] = struct{}{}

		if pm.HostPort != nil {
			portBindings[port] = []nat.PortBinding{
				{HostPort: fmt.Sprintf("%d", *pm.HostPort)},
			}
		}
	}

	return exposedPorts, portBindings
}

func buildMounts(volumes []entity.Volume) []mount.Mount {
	var mounts []mount.Mount
	for _, v := range volumes {
		m := mount.Mount{
			Type:   mount.TypeVolume,
			Target: v.MountPath,
		}
		if v.HostPath != nil {
			m.Type = mount.TypeBind
			m.Source = *v.HostPath
		} else {
			m.Source = v.Name
		}
		mounts = append(mounts, m)
	}
	return mounts
}

func buildResources(cpuLimit *float64, memoryLimit *int64) container.Resources {
	resources := container.Resources{}
	if cpuLimit != nil {
		resources.NanoCPUs = int64(*cpuLimit * 1e9)
	}
	if memoryLimit != nil {
		resources.Memory = *memoryLimit
	}
	return resources
}

func buildRestartPolicy(policy entity.RestartPolicy) container.RestartPolicy {
	switch policy {
	case entity.RestartPolicyAlways:
		return container.RestartPolicy{Name: container.RestartPolicyAlways}
	case entity.RestartPolicyOnFailure:
		return container.RestartPolicy{Name: container.RestartPolicyOnFailure}
	case entity.RestartPolicyUnlessStopped:
		return container.RestartPolicy{Name: container.RestartPolicyUnlessStopped}
	default:
		return container.RestartPolicy{Name: container.RestartPolicyDisabled}
	}
}
