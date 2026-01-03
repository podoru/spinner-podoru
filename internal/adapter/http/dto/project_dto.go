package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

// ProjectResponse represents project data in API responses
type ProjectResponse struct {
	ID           uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	TeamID       uuid.UUID `json:"team_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name         string    `json:"name" example:"My Web App"`
	Slug         string    `json:"slug" example:"my-web-app"`
	Description  *string   `json:"description,omitempty" example:"A production web application"`
	GithubRepo   *string   `json:"github_repo,omitempty" example:"https://github.com/user/repo"`
	GithubBranch string    `json:"github_branch" example:"main"`
	AutoDeploy   bool      `json:"auto_deploy" example:"true"`
	CreatedAt    time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// CreateProjectRequest represents the project creation payload
type CreateProjectRequest struct {
	Name         string  `json:"name" validate:"required,min=2,max=100" example:"My Web App"`
	Slug         string  `json:"slug" validate:"required,slug,min=2,max=100" example:"my-web-app"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=500" example:"A production web application"`
	GithubRepo   *string `json:"github_repo,omitempty" validate:"omitempty,url" example:"https://github.com/user/repo"`
	GithubBranch *string `json:"github_branch,omitempty" example:"main"`
	GithubToken  *string `json:"github_token,omitempty" example:"ghp_xxxxxxxxxxxxxxxxxxxx"`
	AutoDeploy   *bool   `json:"auto_deploy,omitempty" example:"false"`
}

// UpdateProjectRequest represents the project update payload
type UpdateProjectRequest struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=2,max=100" example:"Updated Project Name"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=500" example:"Updated description"`
	GithubRepo   *string `json:"github_repo,omitempty" validate:"omitempty,url" example:"https://github.com/user/new-repo"`
	GithubBranch *string `json:"github_branch,omitempty" example:"develop"`
	GithubToken  *string `json:"github_token,omitempty" example:"ghp_xxxxxxxxxxxxxxxxxxxx"`
	AutoDeploy   *bool   `json:"auto_deploy,omitempty" example:"true"`
}

func ToProjectResponse(project *entity.Project) ProjectResponse {
	return ProjectResponse{
		ID:           project.ID,
		TeamID:       project.TeamID,
		Name:         project.Name,
		Slug:         project.Slug,
		Description:  project.Description,
		GithubRepo:   project.GithubRepo,
		GithubBranch: project.GithubBranch,
		AutoDeploy:   project.AutoDeploy,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
	}
}

func ToProjectsResponse(projects []entity.Project) []ProjectResponse {
	responses := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		responses[i] = ToProjectResponse(&p)
	}
	return responses
}
