package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

// RegisterRequest represents the registration payload
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"securepassword123"`
	Name     string `json:"name" validate:"required,min=2,max=100" example:"John Doe"`
}

// LoginRequest represents the login payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"securepassword123"`
}

// RefreshRequest represents the token refresh payload
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// TokenResponse represents JWT token pair response
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.Rq8IjqbeD8FxP0eKPlM1rRq8IjqbeD8FxP0eKPlM1rQ"`
	ExpiresIn    int64  `json:"expires_in" example:"900"`
	TokenType    string `json:"token_type" example:"Bearer"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string    `json:"email" example:"user@example.com"`
	Name      string    `json:"name" example:"John Doe"`
	Role      string    `json:"role" example:"user"`
	AvatarURL *string   `json:"avatar_url,omitempty" example:"https://example.com/avatar.png"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// AuthResponse represents authentication response with user and tokens
type AuthResponse struct {
	User   UserResponse  `json:"user"`
	Tokens TokenResponse `json:"tokens"`
}

// UpdateUserRequest represents the user update payload
type UpdateUserRequest struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,min=2,max=100" example:"John Doe Updated"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url" example:"https://example.com/new-avatar.png"`
}

// UpdatePasswordRequest represents the password change payload
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"oldpassword123"`
	NewPassword     string `json:"new_password" validate:"required,min=8" example:"newpassword456"`
}

func ToUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      string(user.Role),
		AvatarURL: user.AvatarURL,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToTokenResponse(tokens *entity.TokenPair) TokenResponse {
	return TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		TokenType:    tokens.TokenType,
	}
}

func ToAuthResponse(user *entity.User, tokens *entity.TokenPair) AuthResponse {
	return AuthResponse{
		User:   ToUserResponse(user),
		Tokens: ToTokenResponse(tokens),
	}
}

func (r *RegisterRequest) ToEntity() *entity.UserCreate {
	return &entity.UserCreate{
		Email:    r.Email,
		Password: r.Password,
		Name:     r.Name,
	}
}

func (r *LoginRequest) ToEntity() *entity.LoginRequest {
	return &entity.LoginRequest{
		Email:    r.Email,
		Password: r.Password,
	}
}
