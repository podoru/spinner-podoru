package entity

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID                   uuid.UUID `json:"id"`
	TeamID               uuid.UUID `json:"team_id"`
	Name                 string    `json:"name"`
	Slug                 string    `json:"slug"`
	Description          *string   `json:"description,omitempty"`
	GithubRepo           *string   `json:"github_repo,omitempty"`
	GithubBranch         string    `json:"github_branch"`
	GithubTokenEncrypted []byte    `json:"-"`
	AutoDeploy           bool      `json:"auto_deploy"`
	WebhookSecret        *string   `json:"-"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type ProjectCreate struct {
	Name         string  `json:"name" validate:"required,min=2,max=100"`
	Slug         string  `json:"slug" validate:"required,min=2,max=100,slug"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=500"`
	GithubRepo   *string `json:"github_repo,omitempty" validate:"omitempty,max=500"`
	GithubBranch *string `json:"github_branch,omitempty" validate:"omitempty,max=100"`
	GithubToken  *string `json:"github_token,omitempty"`
	AutoDeploy   *bool   `json:"auto_deploy,omitempty"`
}

type ProjectUpdate struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description  *string `json:"description,omitempty" validate:"omitempty,max=500"`
	GithubRepo   *string `json:"github_repo,omitempty" validate:"omitempty,max=500"`
	GithubBranch *string `json:"github_branch,omitempty" validate:"omitempty,max=100"`
	GithubToken  *string `json:"github_token,omitempty"`
	AutoDeploy   *bool   `json:"auto_deploy,omitempty"`
}

type ProjectWithServices struct {
	Project
	Services []Service `json:"services,omitempty"`
}
