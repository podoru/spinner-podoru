package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/podoru/podoru/internal/domain/entity"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *entity.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Project, error)
	GetByTeamAndSlug(ctx context.Context, teamID uuid.UUID, slug string) (*entity.Project, error)
	Update(ctx context.Context, project *entity.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByTeamID(ctx context.Context, teamID uuid.UUID) ([]entity.Project, error)
	ExistsByTeamAndSlug(ctx context.Context, teamID uuid.UUID, slug string) (bool, error)
}
