package docker

import (
	"context"

	"github.com/docker/docker/client"

	"github.com/podoru/spinner-podoru/internal/infrastructure/config"
)

// Client wraps the Docker client
type Client struct {
	cli *client.Client
}

// NewClient creates a new Docker client
func NewClient(cfg *config.DockerConfig) (*Client, error) {
	opts := []client.Opt{
		client.WithAPIVersionNegotiation(),
	}

	if cfg.Host != "" {
		opts = append(opts, client.WithHost(cfg.Host))
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	return &Client{cli: cli}, nil
}

// Close closes the Docker client
func (c *Client) Close() error {
	return c.cli.Close()
}

// Ping checks if the Docker daemon is accessible
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx)
	return err
}

// Client returns the underlying Docker client
func (c *Client) Client() *client.Client {
	return c.cli
}
