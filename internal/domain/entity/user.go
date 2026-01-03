package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleSuperAdmin UserRole = "superadmin"
	UserRoleUser       UserRole = "user"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	Role         UserRole  `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) IsSuperAdmin() bool {
	return u.Role == UserRoleSuperAdmin
}

type UserCreate struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
}

type UserUpdate struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}

type UserPasswordUpdate struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}
