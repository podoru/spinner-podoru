package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

// TeamResponse represents team data in API responses
type TeamResponse struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"My Awesome Team"`
	Slug        string    `json:"slug" example:"my-awesome-team"`
	Description *string   `json:"description,omitempty" example:"A team for awesome projects"`
	OwnerID     uuid.UUID `json:"owner_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// TeamWithRoleResponse represents team with user's role
type TeamWithRoleResponse struct {
	ID          uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" example:"My Awesome Team"`
	Slug        string    `json:"slug" example:"my-awesome-team"`
	Description *string   `json:"description,omitempty" example:"A team for awesome projects"`
	OwnerID     uuid.UUID `json:"owner_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Role        string    `json:"role" example:"owner"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// CreateTeamRequest represents the team creation payload
type CreateTeamRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100" example:"My Awesome Team"`
	Slug        string  `json:"slug" validate:"required,slug,min=2,max=100" example:"my-awesome-team"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500" example:"A team for awesome projects"`
}

// UpdateTeamRequest represents the team update payload
type UpdateTeamRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100" example:"Updated Team Name"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500" example:"Updated description"`
}

// TeamMemberResponse represents team member data in API responses
type TeamMemberResponse struct {
	ID        uuid.UUID     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	TeamID    uuid.UUID     `json:"team_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	UserID    uuid.UUID     `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Role      string        `json:"role" example:"member"`
	User      *UserResponse `json:"user,omitempty"`
	CreatedAt time.Time     `json:"created_at" example:"2024-01-15T10:30:00Z"`
}

// AddTeamMemberRequest represents the add member payload
type AddTeamMemberRequest struct {
	Email string `json:"email" validate:"required,email" example:"newmember@example.com"`
	Role  string `json:"role" validate:"required,oneof=admin member" example:"member"`
}

// UpdateTeamMemberRequest represents the update member role payload
type UpdateTeamMemberRequest struct {
	Role string `json:"role" validate:"required,oneof=admin member" example:"admin"`
}

func ToTeamResponse(team *entity.Team) TeamResponse {
	return TeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		Slug:        team.Slug,
		Description: team.Description,
		OwnerID:     team.OwnerID,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}
}

func ToTeamWithRoleResponse(twr *entity.TeamWithRole) TeamWithRoleResponse {
	return TeamWithRoleResponse{
		ID:          twr.ID,
		Name:        twr.Name,
		Slug:        twr.Slug,
		Description: twr.Description,
		OwnerID:     twr.OwnerID,
		Role:        string(twr.Role),
		CreatedAt:   twr.CreatedAt,
		UpdatedAt:   twr.UpdatedAt,
	}
}

func ToTeamMemberResponse(member *entity.TeamMember) TeamMemberResponse {
	resp := TeamMemberResponse{
		ID:        member.ID,
		TeamID:    member.TeamID,
		UserID:    member.UserID,
		Role:      string(member.Role),
		CreatedAt: member.CreatedAt,
	}
	if member.User != nil {
		userResp := ToUserResponse(member.User)
		resp.User = &userResp
	}
	return resp
}
