package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

type TeamRepository interface {
	Create(ctx context.Context, team *entity.Team) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Team, error)
	Update(ctx context.Context, team *entity.Team) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.TeamWithRole, error)
	ExistsBySlug(ctx context.Context, slug string) (bool, error)
}

type TeamMemberRepository interface {
	Create(ctx context.Context, member *entity.TeamMember) error
	GetByTeamAndUser(ctx context.Context, teamID, userID uuid.UUID) (*entity.TeamMember, error)
	Update(ctx context.Context, member *entity.TeamMember) error
	Delete(ctx context.Context, teamID, userID uuid.UUID) error
	ListByTeamID(ctx context.Context, teamID uuid.UUID) ([]entity.TeamMember, error)
	CountByTeamID(ctx context.Context, teamID uuid.UUID) (int, error)
}
