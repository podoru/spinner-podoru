package entity

import (
	"time"

	"github.com/google/uuid"
)

type Domain struct {
	ID         uuid.UUID `json:"id"`
	ServiceID  uuid.UUID `json:"service_id"`
	Domain     string    `json:"domain"`
	SSLEnabled bool      `json:"ssl_enabled"`
	SSLAuto    bool      `json:"ssl_auto"`
	CreatedAt  time.Time `json:"created_at"`
}

type DomainCreate struct {
	Domain     string `json:"domain" validate:"required,fqdn,max=255"`
	SSLEnabled *bool  `json:"ssl_enabled,omitempty"`
	SSLAuto    *bool  `json:"ssl_auto,omitempty"`
}
