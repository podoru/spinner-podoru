package mocks

import (
	"context"

	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	CreateFunc        func(ctx context.Context, user *entity.User) error
	GetByIDFunc       func(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmailFunc    func(ctx context.Context, email string) (*entity.User, error)
	UpdateFunc        func(ctx context.Context, user *entity.User) error
	DeleteFunc        func(ctx context.Context, id uuid.UUID) error
	ExistsByEmailFunc func(ctx context.Context, email string) (bool, error)
	CountFunc         func(ctx context.Context) (int64, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.ExistsByEmailFunc != nil {
		return m.ExistsByEmailFunc(ctx, email)
	}
	return false, nil
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx)
	}
	return 0, nil
}

// MockRefreshTokenRepository is a mock implementation of RefreshTokenRepository
type MockRefreshTokenRepository struct {
	CreateFunc            func(ctx context.Context, token *entity.RefreshToken) error
	GetByTokenHashFunc    func(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	DeleteByUserIDFunc    func(ctx context.Context, userID uuid.UUID) error
	DeleteByTokenHashFunc func(ctx context.Context, tokenHash string) error
	DeleteExpiredFunc     func(ctx context.Context) error
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, token)
	}
	return nil
}

func (m *MockRefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	if m.GetByTokenHashFunc != nil {
		return m.GetByTokenHashFunc(ctx, tokenHash)
	}
	return nil, nil
}

func (m *MockRefreshTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	if m.DeleteByUserIDFunc != nil {
		return m.DeleteByUserIDFunc(ctx, userID)
	}
	return nil
}

func (m *MockRefreshTokenRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	if m.DeleteByTokenHashFunc != nil {
		return m.DeleteByTokenHashFunc(ctx, tokenHash)
	}
	return nil
}

func (m *MockRefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	if m.DeleteExpiredFunc != nil {
		return m.DeleteExpiredFunc(ctx)
	}
	return nil
}

// MockTeamRepository is a mock implementation of TeamRepository
type MockTeamRepository struct {
	CreateFunc       func(ctx context.Context, team *entity.Team) error
	GetByIDFunc      func(ctx context.Context, id uuid.UUID) (*entity.Team, error)
	GetBySlugFunc    func(ctx context.Context, slug string) (*entity.Team, error)
	UpdateFunc       func(ctx context.Context, team *entity.Team) error
	DeleteFunc       func(ctx context.Context, id uuid.UUID) error
	ListByUserIDFunc func(ctx context.Context, userID uuid.UUID) ([]entity.TeamWithRole, error)
	ExistsBySlugFunc func(ctx context.Context, slug string) (bool, error)
}

func (m *MockTeamRepository) Create(ctx context.Context, team *entity.Team) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, team)
	}
	return nil
}

func (m *MockTeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTeamRepository) GetBySlug(ctx context.Context, slug string) (*entity.Team, error) {
	if m.GetBySlugFunc != nil {
		return m.GetBySlugFunc(ctx, slug)
	}
	return nil, nil
}

func (m *MockTeamRepository) Update(ctx context.Context, team *entity.Team) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, team)
	}
	return nil
}

func (m *MockTeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockTeamRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.TeamWithRole, error) {
	if m.ListByUserIDFunc != nil {
		return m.ListByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockTeamRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	if m.ExistsBySlugFunc != nil {
		return m.ExistsBySlugFunc(ctx, slug)
	}
	return false, nil
}

// MockTeamMemberRepository is a mock implementation of TeamMemberRepository
type MockTeamMemberRepository struct {
	CreateFunc          func(ctx context.Context, member *entity.TeamMember) error
	GetByTeamAndUserFunc func(ctx context.Context, teamID, userID uuid.UUID) (*entity.TeamMember, error)
	UpdateFunc          func(ctx context.Context, member *entity.TeamMember) error
	DeleteFunc          func(ctx context.Context, teamID, userID uuid.UUID) error
	ListByTeamIDFunc    func(ctx context.Context, teamID uuid.UUID) ([]entity.TeamMember, error)
	CountByTeamIDFunc   func(ctx context.Context, teamID uuid.UUID) (int, error)
}

func (m *MockTeamMemberRepository) Create(ctx context.Context, member *entity.TeamMember) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, member)
	}
	return nil
}

func (m *MockTeamMemberRepository) GetByTeamAndUser(ctx context.Context, teamID, userID uuid.UUID) (*entity.TeamMember, error) {
	if m.GetByTeamAndUserFunc != nil {
		return m.GetByTeamAndUserFunc(ctx, teamID, userID)
	}
	return nil, nil
}

func (m *MockTeamMemberRepository) Update(ctx context.Context, member *entity.TeamMember) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, member)
	}
	return nil
}

func (m *MockTeamMemberRepository) Delete(ctx context.Context, teamID, userID uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, teamID, userID)
	}
	return nil
}

func (m *MockTeamMemberRepository) ListByTeamID(ctx context.Context, teamID uuid.UUID) ([]entity.TeamMember, error) {
	if m.ListByTeamIDFunc != nil {
		return m.ListByTeamIDFunc(ctx, teamID)
	}
	return nil, nil
}

func (m *MockTeamMemberRepository) CountByTeamID(ctx context.Context, teamID uuid.UUID) (int, error) {
	if m.CountByTeamIDFunc != nil {
		return m.CountByTeamIDFunc(ctx, teamID)
	}
	return 0, nil
}

// MockProjectRepository is a mock implementation of ProjectRepository
type MockProjectRepository struct {
	CreateFunc              func(ctx context.Context, project *entity.Project) error
	GetByIDFunc             func(ctx context.Context, id uuid.UUID) (*entity.Project, error)
	GetByTeamAndSlugFunc    func(ctx context.Context, teamID uuid.UUID, slug string) (*entity.Project, error)
	UpdateFunc              func(ctx context.Context, project *entity.Project) error
	DeleteFunc              func(ctx context.Context, id uuid.UUID) error
	ListByTeamIDFunc        func(ctx context.Context, teamID uuid.UUID) ([]entity.Project, error)
	ExistsByTeamAndSlugFunc func(ctx context.Context, teamID uuid.UUID, slug string) (bool, error)
}

func (m *MockProjectRepository) Create(ctx context.Context, project *entity.Project) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, project)
	}
	return nil
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Project, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockProjectRepository) GetByTeamAndSlug(ctx context.Context, teamID uuid.UUID, slug string) (*entity.Project, error) {
	if m.GetByTeamAndSlugFunc != nil {
		return m.GetByTeamAndSlugFunc(ctx, teamID, slug)
	}
	return nil, nil
}

func (m *MockProjectRepository) Update(ctx context.Context, project *entity.Project) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, project)
	}
	return nil
}

func (m *MockProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockProjectRepository) ListByTeamID(ctx context.Context, teamID uuid.UUID) ([]entity.Project, error) {
	if m.ListByTeamIDFunc != nil {
		return m.ListByTeamIDFunc(ctx, teamID)
	}
	return nil, nil
}

func (m *MockProjectRepository) ExistsByTeamAndSlug(ctx context.Context, teamID uuid.UUID, slug string) (bool, error) {
	if m.ExistsByTeamAndSlugFunc != nil {
		return m.ExistsByTeamAndSlugFunc(ctx, teamID, slug)
	}
	return false, nil
}

// MockServiceRepository is a mock implementation of ServiceRepository
type MockServiceRepository struct {
	CreateFunc                 func(ctx context.Context, service *entity.Service) error
	GetByIDFunc                func(ctx context.Context, id uuid.UUID) (*entity.Service, error)
	GetByProjectAndSlugFunc    func(ctx context.Context, projectID uuid.UUID, slug string) (*entity.Service, error)
	UpdateFunc                 func(ctx context.Context, service *entity.Service) error
	DeleteFunc                 func(ctx context.Context, id uuid.UUID) error
	ListByProjectIDFunc        func(ctx context.Context, projectID uuid.UUID) ([]entity.Service, error)
	ExistsByProjectAndSlugFunc func(ctx context.Context, projectID uuid.UUID, slug string) (bool, error)
	UpdateStatusFunc           func(ctx context.Context, id uuid.UUID, status entity.ServiceStatus) error
	UpdateContainerIDFunc      func(ctx context.Context, id uuid.UUID, containerID *string) error
	UpdateSwarmServiceIDFunc   func(ctx context.Context, id uuid.UUID, swarmServiceID *string) error
}

func (m *MockServiceRepository) Create(ctx context.Context, service *entity.Service) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, service)
	}
	return nil
}

func (m *MockServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockServiceRepository) GetByProjectAndSlug(ctx context.Context, projectID uuid.UUID, slug string) (*entity.Service, error) {
	if m.GetByProjectAndSlugFunc != nil {
		return m.GetByProjectAndSlugFunc(ctx, projectID, slug)
	}
	return nil, nil
}

func (m *MockServiceRepository) Update(ctx context.Context, service *entity.Service) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, service)
	}
	return nil
}

func (m *MockServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockServiceRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]entity.Service, error) {
	if m.ListByProjectIDFunc != nil {
		return m.ListByProjectIDFunc(ctx, projectID)
	}
	return nil, nil
}

func (m *MockServiceRepository) ExistsByProjectAndSlug(ctx context.Context, projectID uuid.UUID, slug string) (bool, error) {
	if m.ExistsByProjectAndSlugFunc != nil {
		return m.ExistsByProjectAndSlugFunc(ctx, projectID, slug)
	}
	return false, nil
}

func (m *MockServiceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ServiceStatus) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *MockServiceRepository) UpdateContainerID(ctx context.Context, id uuid.UUID, containerID *string) error {
	if m.UpdateContainerIDFunc != nil {
		return m.UpdateContainerIDFunc(ctx, id, containerID)
	}
	return nil
}

func (m *MockServiceRepository) UpdateSwarmServiceID(ctx context.Context, id uuid.UUID, swarmServiceID *string) error {
	if m.UpdateSwarmServiceIDFunc != nil {
		return m.UpdateSwarmServiceIDFunc(ctx, id, swarmServiceID)
	}
	return nil
}
