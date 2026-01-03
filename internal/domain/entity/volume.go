package entity

import (
	"time"

	"github.com/google/uuid"
)

type Volume struct {
	ID        uuid.UUID `json:"id"`
	ServiceID uuid.UUID `json:"service_id"`
	Name      string    `json:"name"`
	MountPath string    `json:"mount_path"`
	HostPath  *string   `json:"host_path,omitempty"`
	Driver    string    `json:"driver"`
	CreatedAt time.Time `json:"created_at"`
}

type VolumeCreate struct {
	Name      string  `json:"name" validate:"required,min=1,max=255"`
	MountPath string  `json:"mount_path" validate:"required,min=1,max=500"`
	HostPath  *string `json:"host_path,omitempty" validate:"omitempty,max=500"`
	Driver    *string `json:"driver,omitempty" validate:"omitempty,oneof=local"`
}
