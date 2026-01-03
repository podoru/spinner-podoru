package team

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
	"github.com/podoru/spinner-podoru/internal/domain/repository"
)

var (
	ErrTeamNotFound      = errors.New("team not found")
	ErrSlugAlreadyExists = errors.New("slug already exists")
	ErrNotTeamMember     = errors.New("not a team member")
	ErrNotTeamOwner      = errors.New("not team owner")
	ErrNotTeamAdmin      = errors.New("requires admin or owner role")
	ErrCannotRemoveOwner = errors.New("cannot remove team owner")
	ErrUserNotFound      = errors.New("user not found")
	ErrAlreadyMember     = errors.New("user is already a team member")
)

type UseCase struct {
	teamRepo       repository.TeamRepository
	teamMemberRepo repository.TeamMemberRepository
	userRepo       repository.UserRepository
}

func NewUseCase(
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	userRepo repository.UserRepository,
) *UseCase {
	return &UseCase{
		teamRepo:       teamRepo,
		teamMemberRepo: teamMemberRepo,
		userRepo:       userRepo,
	}
}

func (uc *UseCase) Create(ctx context.Context, userID uuid.UUID, input *entity.TeamCreate) (*entity.Team, error) {
	exists, err := uc.teamRepo.ExistsBySlug(ctx, input.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugAlreadyExists
	}

	now := time.Now()
	team := &entity.Team{
		ID:          uuid.New(),
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		OwnerID:     userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.teamRepo.Create(ctx, team); err != nil {
		return nil, err
	}

	member := &entity.TeamMember{
		ID:        uuid.New(),
		TeamID:    team.ID,
		UserID:    userID,
		Role:      entity.TeamRoleOwner,
		CreatedAt: now,
	}

	if err := uc.teamMemberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	return team, nil
}

func (uc *UseCase) GetByID(ctx context.Context, userID, teamID uuid.UUID) (*entity.Team, error) {
	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	team, err := uc.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrTeamNotFound
	}

	return team, nil
}

func (uc *UseCase) List(ctx context.Context, userID uuid.UUID) ([]entity.TeamWithRole, error) {
	return uc.teamRepo.ListByUserID(ctx, userID)
}

func (uc *UseCase) Update(ctx context.Context, userID, teamID uuid.UUID, input *entity.TeamUpdate) (*entity.Team, error) {
	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}
	if member.Role != entity.TeamRoleOwner && member.Role != entity.TeamRoleAdmin {
		return nil, ErrNotTeamAdmin
	}

	team, err := uc.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrTeamNotFound
	}

	if input.Name != nil {
		team.Name = *input.Name
	}
	if input.Description != nil {
		team.Description = input.Description
	}
	team.UpdatedAt = time.Now()

	if err := uc.teamRepo.Update(ctx, team); err != nil {
		return nil, err
	}

	return team, nil
}

func (uc *UseCase) Delete(ctx context.Context, userID, teamID uuid.UUID) error {
	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return ErrNotTeamMember
	}
	if member.Role != entity.TeamRoleOwner {
		return ErrNotTeamOwner
	}

	return uc.teamRepo.Delete(ctx, teamID)
}

func (uc *UseCase) ListMembers(ctx context.Context, userID, teamID uuid.UUID) ([]entity.TeamMember, error) {
	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	return uc.teamMemberRepo.ListByTeamID(ctx, teamID)
}

func (uc *UseCase) AddMember(ctx context.Context, userID, teamID uuid.UUID, input *entity.TeamMemberCreate) (*entity.TeamMember, error) {
	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}
	if member.Role != entity.TeamRoleOwner && member.Role != entity.TeamRoleAdmin {
		return nil, ErrNotTeamAdmin
	}

	newUser, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if newUser == nil {
		return nil, ErrUserNotFound
	}

	existing, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, newUser.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrAlreadyMember
	}

	newMember := &entity.TeamMember{
		ID:        uuid.New(),
		TeamID:    teamID,
		UserID:    newUser.ID,
		Role:      input.Role,
		CreatedAt: time.Now(),
		User:      newUser,
	}

	if err := uc.teamMemberRepo.Create(ctx, newMember); err != nil {
		return nil, err
	}

	return newMember, nil
}

func (uc *UseCase) UpdateMember(ctx context.Context, userID, teamID, targetUserID uuid.UUID, input *entity.TeamMemberUpdate) (*entity.TeamMember, error) {
	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}
	if member.Role != entity.TeamRoleOwner && member.Role != entity.TeamRoleAdmin {
		return nil, ErrNotTeamAdmin
	}

	targetMember, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, targetUserID)
	if err != nil {
		return nil, err
	}
	if targetMember == nil {
		return nil, ErrNotTeamMember
	}

	if targetMember.Role == entity.TeamRoleOwner {
		return nil, ErrCannotRemoveOwner
	}

	targetMember.Role = input.Role

	if err := uc.teamMemberRepo.Update(ctx, targetMember); err != nil {
		return nil, err
	}

	return targetMember, nil
}

func (uc *UseCase) RemoveMember(ctx context.Context, userID, teamID, targetUserID uuid.UUID) error {
	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return ErrNotTeamMember
	}
	if member.Role != entity.TeamRoleOwner && member.Role != entity.TeamRoleAdmin {
		return ErrNotTeamAdmin
	}

	targetMember, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, targetUserID)
	if err != nil {
		return err
	}
	if targetMember == nil {
		return ErrNotTeamMember
	}

	if targetMember.Role == entity.TeamRoleOwner {
		return ErrCannotRemoveOwner
	}

	return uc.teamMemberRepo.Delete(ctx, teamID, targetUserID)
}

func (uc *UseCase) CheckMembership(ctx context.Context, userID, teamID uuid.UUID) (*entity.TeamMember, error) {
	return uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, userID)
}
