package project

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/podoru/podoru/internal/domain/entity"
	"github.com/podoru/podoru/internal/domain/repository"
	"github.com/podoru/podoru/pkg/crypto"
)

var (
	ErrProjectNotFound   = errors.New("project not found")
	ErrSlugAlreadyExists = errors.New("slug already exists")
	ErrNotTeamMember     = errors.New("not a team member")
	ErrNotTeamAdmin      = errors.New("requires admin or owner role")
)

type UseCase struct {
	projectRepo    repository.ProjectRepository
	teamMemberRepo repository.TeamMemberRepository
	encryptor      *crypto.Encryptor
}

func NewUseCase(
	projectRepo repository.ProjectRepository,
	teamMemberRepo repository.TeamMemberRepository,
	encryptor *crypto.Encryptor,
) *UseCase {
	return &UseCase{
		projectRepo:    projectRepo,
		teamMemberRepo: teamMemberRepo,
		encryptor:      encryptor,
	}
}

func (uc *UseCase) Create(ctx context.Context, userID, teamID uuid.UUID, input *entity.ProjectCreate) (*entity.Project, error) {
	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	exists, err := uc.projectRepo.ExistsByTeamAndSlug(ctx, teamID, input.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugAlreadyExists
	}

	now := time.Now()
	project := &entity.Project{
		ID:           uuid.New(),
		TeamID:       teamID,
		Name:         input.Name,
		Slug:         input.Slug,
		Description:  input.Description,
		GithubRepo:   input.GithubRepo,
		GithubBranch: "main",
		AutoDeploy:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if input.GithubBranch != nil {
		project.GithubBranch = *input.GithubBranch
	}

	if input.AutoDeploy != nil {
		project.AutoDeploy = *input.AutoDeploy
	}

	if input.GithubToken != nil && *input.GithubToken != "" {
		encrypted, err := uc.encryptor.Encrypt([]byte(*input.GithubToken))
		if err != nil {
			return nil, err
		}
		project.GithubTokenEncrypted = encrypted
	}

	webhookSecret, err := crypto.GenerateRandomString(32)
	if err != nil {
		return nil, err
	}
	project.WebhookSecret = &webhookSecret

	if err := uc.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (uc *UseCase) GetByID(ctx context.Context, userID, projectID uuid.UUID) (*entity.Project, error) {
	project, err := uc.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, project.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	return project, nil
}

func (uc *UseCase) ListByTeam(ctx context.Context, userID, teamID uuid.UUID) ([]entity.Project, error) {
	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, teamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	return uc.projectRepo.ListByTeamID(ctx, teamID)
}

func (uc *UseCase) Update(ctx context.Context, userID, projectID uuid.UUID, input *entity.ProjectUpdate) (*entity.Project, error) {
	project, err := uc.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}

	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, project.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}
	if member.Role == entity.TeamRoleMember {
		return nil, ErrNotTeamAdmin
	}

	if input.Name != nil {
		project.Name = *input.Name
	}
	if input.Description != nil {
		project.Description = input.Description
	}
	if input.GithubRepo != nil {
		project.GithubRepo = input.GithubRepo
	}
	if input.GithubBranch != nil {
		project.GithubBranch = *input.GithubBranch
	}
	if input.AutoDeploy != nil {
		project.AutoDeploy = *input.AutoDeploy
	}
	if input.GithubToken != nil {
		if *input.GithubToken == "" {
			project.GithubTokenEncrypted = nil
		} else {
			encrypted, err := uc.encryptor.Encrypt([]byte(*input.GithubToken))
			if err != nil {
				return nil, err
			}
			project.GithubTokenEncrypted = encrypted
		}
	}
	project.UpdatedAt = time.Now()

	if err := uc.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (uc *UseCase) Delete(ctx context.Context, userID, projectID uuid.UUID) error {
	project, err := uc.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return ErrProjectNotFound
	}

	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, project.TeamID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return ErrNotTeamMember
	}
	if member.Role == entity.TeamRoleMember {
		return ErrNotTeamAdmin
	}

	return uc.projectRepo.Delete(ctx, projectID)
}

func (uc *UseCase) GetProjectWithTeamCheck(ctx context.Context, userID, projectID uuid.UUID) (*entity.Project, *entity.TeamMember, error) {
	project, err := uc.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, nil, err
	}
	if project == nil {
		return nil, nil, ErrProjectNotFound
	}

	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, project.TeamID, userID)
	if err != nil {
		return nil, nil, err
	}
	if member == nil {
		return nil, nil, ErrNotTeamMember
	}

	return project, member, nil
}
