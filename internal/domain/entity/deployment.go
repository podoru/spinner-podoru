package entity

import (
	"time"

	"github.com/google/uuid"
)

type DeploymentStatus string

const (
	DeploymentStatusPending   DeploymentStatus = "pending"
	DeploymentStatusBuilding  DeploymentStatus = "building"
	DeploymentStatusDeploying DeploymentStatus = "deploying"
	DeploymentStatusSuccess   DeploymentStatus = "success"
	DeploymentStatusFailed    DeploymentStatus = "failed"
)

type Deployment struct {
	ID            uuid.UUID        `json:"id"`
	ServiceID     uuid.UUID        `json:"service_id"`
	TriggeredBy   *uuid.UUID       `json:"triggered_by,omitempty"`
	CommitSHA     *string          `json:"commit_sha,omitempty"`
	CommitMessage *string          `json:"commit_message,omitempty"`
	Status        DeploymentStatus `json:"status"`
	Logs          *string          `json:"logs,omitempty"`
	StartedAt     time.Time        `json:"started_at"`
	FinishedAt    *time.Time       `json:"finished_at,omitempty"`
}

type DeploymentCreate struct {
	CommitSHA     *string `json:"commit_sha,omitempty"`
	CommitMessage *string `json:"commit_message,omitempty"`
}
