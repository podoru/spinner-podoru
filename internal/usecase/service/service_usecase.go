package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/podoru/podoru/internal/domain/entity"
	"github.com/podoru/podoru/internal/domain/repository"
	"github.com/podoru/podoru/pkg/crypto"
)

var (
	ErrServiceNotFound   = errors.New("service not found")
	ErrProjectNotFound   = errors.New("project not found")
	ErrSlugAlreadyExists = errors.New("slug already exists")
	ErrNotTeamMember     = errors.New("not a team member")
	ErrNotTeamAdmin      = errors.New("requires admin or owner role")
	ErrImageRequired     = errors.New("image is required for image deploy type")
)

type UseCase struct {
	serviceRepo    repository.ServiceRepository
	projectRepo    repository.ProjectRepository
	teamMemberRepo repository.TeamMemberRepository
	encryptor      *crypto.Encryptor
}

func NewUseCase(
	serviceRepo repository.ServiceRepository,
	projectRepo repository.ProjectRepository,
	teamMemberRepo repository.TeamMemberRepository,
	encryptor *crypto.Encryptor,
) *UseCase {
	return &UseCase{
		serviceRepo:    serviceRepo,
		projectRepo:    projectRepo,
		teamMemberRepo: teamMemberRepo,
		encryptor:      encryptor,
	}
}

func (uc *UseCase) Create(ctx context.Context, userID, projectID uuid.UUID, input *entity.ServiceCreate) (*entity.Service, error) {
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

	exists, err := uc.serviceRepo.ExistsByProjectAndSlug(ctx, projectID, input.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSlugAlreadyExists
	}

	if input.DeployType == entity.DeployTypeImage && (input.Image == nil || *input.Image == "") {
		return nil, ErrImageRequired
	}

	now := time.Now()
	service := &entity.Service{
		ID:                  uuid.New(),
		ProjectID:           projectID,
		Name:                input.Name,
		Slug:                input.Slug,
		DeployType:          input.DeployType,
		Image:               input.Image,
		DockerfilePath:      "Dockerfile",
		BuildContext:        ".",
		ComposeFile:         input.ComposeFile,
		Replicas:            1,
		HealthCheckInterval: 30,
		RestartPolicy:       entity.RestartPolicyUnlessStopped,
		Status:              entity.ServiceStatusStopped,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if input.DockerfilePath != nil {
		service.DockerfilePath = *input.DockerfilePath
	}
	if input.BuildContext != nil {
		service.BuildContext = *input.BuildContext
	}
	if input.Replicas != nil {
		service.Replicas = *input.Replicas
	}
	if input.CPULimit != nil {
		service.CPULimit = input.CPULimit
	}
	if input.MemoryLimit != nil {
		service.MemoryLimit = input.MemoryLimit
	}
	if input.HealthCheckPath != nil {
		service.HealthCheckPath = input.HealthCheckPath
	}
	if input.HealthCheckInterval != nil {
		service.HealthCheckInterval = *input.HealthCheckInterval
	}
	if input.RestartPolicy != nil {
		service.RestartPolicy = *input.RestartPolicy
	}

	if input.EnvVars != nil && len(input.EnvVars) > 0 {
		envJSON, err := json.Marshal(input.EnvVars)
		if err != nil {
			return nil, err
		}
		encrypted, err := uc.encryptor.Encrypt(envJSON)
		if err != nil {
			return nil, err
		}
		service.EnvVarsEncrypted = encrypted
	}

	if err := uc.serviceRepo.Create(ctx, service); err != nil {
		return nil, err
	}

	return service, nil
}

func (uc *UseCase) GetByID(ctx context.Context, userID, serviceID uuid.UUID) (*entity.Service, error) {
	service, err := uc.serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrServiceNotFound
	}

	project, err := uc.projectRepo.GetByID(ctx, service.ProjectID)
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

	return service, nil
}

func (uc *UseCase) ListByProject(ctx context.Context, userID, projectID uuid.UUID) ([]entity.Service, error) {
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

	return uc.serviceRepo.ListByProjectID(ctx, projectID)
}

func (uc *UseCase) Update(ctx context.Context, userID, serviceID uuid.UUID, input *entity.ServiceUpdate) (*entity.Service, error) {
	service, err := uc.serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrServiceNotFound
	}

	project, err := uc.projectRepo.GetByID(ctx, service.ProjectID)
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

	if input.Name != nil {
		service.Name = *input.Name
	}
	if input.Image != nil {
		service.Image = input.Image
	}
	if input.DockerfilePath != nil {
		service.DockerfilePath = *input.DockerfilePath
	}
	if input.BuildContext != nil {
		service.BuildContext = *input.BuildContext
	}
	if input.ComposeFile != nil {
		service.ComposeFile = input.ComposeFile
	}
	if input.Replicas != nil {
		service.Replicas = *input.Replicas
	}
	if input.CPULimit != nil {
		service.CPULimit = input.CPULimit
	}
	if input.MemoryLimit != nil {
		service.MemoryLimit = input.MemoryLimit
	}
	if input.HealthCheckPath != nil {
		service.HealthCheckPath = input.HealthCheckPath
	}
	if input.HealthCheckInterval != nil {
		service.HealthCheckInterval = *input.HealthCheckInterval
	}
	if input.RestartPolicy != nil {
		service.RestartPolicy = *input.RestartPolicy
	}

	if input.EnvVars != nil {
		if len(input.EnvVars) > 0 {
			envJSON, err := json.Marshal(input.EnvVars)
			if err != nil {
				return nil, err
			}
			encrypted, err := uc.encryptor.Encrypt(envJSON)
			if err != nil {
				return nil, err
			}
			service.EnvVarsEncrypted = encrypted
		} else {
			service.EnvVarsEncrypted = nil
		}
	}

	service.UpdatedAt = time.Now()

	if err := uc.serviceRepo.Update(ctx, service); err != nil {
		return nil, err
	}

	return service, nil
}

func (uc *UseCase) Delete(ctx context.Context, userID, serviceID uuid.UUID) error {
	service, err := uc.serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return err
	}
	if service == nil {
		return ErrServiceNotFound
	}

	project, err := uc.projectRepo.GetByID(ctx, service.ProjectID)
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

	return uc.serviceRepo.Delete(ctx, serviceID)
}

func (uc *UseCase) Scale(ctx context.Context, userID, serviceID uuid.UUID, replicas int) (*entity.Service, error) {
	service, err := uc.serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrServiceNotFound
	}

	project, err := uc.projectRepo.GetByID(ctx, service.ProjectID)
	if err != nil {
		return nil, err
	}

	member, err := uc.teamMemberRepo.GetByTeamAndUser(ctx, project.TeamID, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrNotTeamMember
	}

	service.Replicas = replicas
	service.UpdatedAt = time.Now()

	if err := uc.serviceRepo.Update(ctx, service); err != nil {
		return nil, err
	}

	return service, nil
}
