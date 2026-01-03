package entity

import (
	"time"

	"github.com/google/uuid"
)

type PortMapping struct {
	ID            uuid.UUID `json:"id"`
	ServiceID     uuid.UUID `json:"service_id"`
	ContainerPort int       `json:"container_port"`
	HostPort      *int      `json:"host_port,omitempty"`
	Protocol      string    `json:"protocol"`
	CreatedAt     time.Time `json:"created_at"`
}

type PortMappingCreate struct {
	ContainerPort int    `json:"container_port" validate:"required,min=1,max=65535"`
	HostPort      *int   `json:"host_port,omitempty" validate:"omitempty,min=1,max=65535"`
	Protocol      string `json:"protocol,omitempty" validate:"omitempty,oneof=tcp udp"`
}
