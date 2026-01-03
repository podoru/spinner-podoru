package entity

import (
	"time"

	"github.com/google/uuid"
)

type NetworkDriver string

const (
	NetworkDriverBridge  NetworkDriver = "bridge"
	NetworkDriverOverlay NetworkDriver = "overlay"
)

type Network struct {
	ID              uuid.UUID     `json:"id"`
	ProjectID       *uuid.UUID    `json:"project_id,omitempty"`
	Name            string        `json:"name"`
	DockerNetworkID *string       `json:"docker_network_id,omitempty"`
	Driver          NetworkDriver `json:"driver"`
	Subnet          *string       `json:"subnet,omitempty"`
	Gateway         *string       `json:"gateway,omitempty"`
	IsDefault       bool          `json:"is_default"`
	CreatedAt       time.Time     `json:"created_at"`
}

type NetworkCreate struct {
	Name    string         `json:"name" validate:"required,min=1,max=255"`
	Driver  *NetworkDriver `json:"driver,omitempty" validate:"omitempty,oneof=bridge overlay"`
	Subnet  *string        `json:"subnet,omitempty" validate:"omitempty,cidr"`
	Gateway *string        `json:"gateway,omitempty" validate:"omitempty,ip"`
}
