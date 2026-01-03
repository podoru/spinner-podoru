package entity

import (
	"time"

	"github.com/google/uuid"
)

type TeamRole string

const (
	TeamRoleOwner  TeamRole = "owner"
	TeamRoleAdmin  TeamRole = "admin"
	TeamRoleMember TeamRole = "member"
)

type Team struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TeamMember struct {
	ID        uuid.UUID `json:"id"`
	TeamID    uuid.UUID `json:"team_id"`
	UserID    uuid.UUID `json:"user_id"`
	Role      TeamRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

type TeamCreate struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Slug        string  `json:"slug" validate:"required,min=2,max=100,slug"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

type TeamUpdate struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
}

type TeamMemberCreate struct {
	Email string   `json:"email" validate:"required,email"`
	Role  TeamRole `json:"role" validate:"required,oneof=admin member"`
}

type TeamMemberUpdate struct {
	Role TeamRole `json:"role" validate:"required,oneof=admin member"`
}

type TeamWithRole struct {
	Team
	Role TeamRole `json:"role"`
}
