package team_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/podoru/podoru/internal/domain/entity"
	"github.com/podoru/podoru/internal/mocks"
	"github.com/podoru/podoru/internal/usecase/team"
)

func TestCreate_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	teamRepo := &mocks.MockTeamRepository{
		ExistsBySlugFunc: func(ctx context.Context, slug string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, tm *entity.Team) error {
			return nil
		},
	}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		CreateFunc: func(ctx context.Context, member *entity.TeamMember) error {
			if member.Role != entity.TeamRoleOwner {
				t.Errorf("expected role %s, got %s", entity.TeamRoleOwner, member.Role)
			}
			return nil
		},
	}

	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	input := &entity.TeamCreate{
		Name: "Test Team",
		Slug: "test-team",
	}

	result, err := uc.Create(ctx, userID, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected team, got nil")
	}

	if result.Name != input.Name {
		t.Errorf("expected name %s, got %s", input.Name, result.Name)
	}

	if result.OwnerID != userID {
		t.Errorf("expected owner ID %s, got %s", userID, result.OwnerID)
	}
}

func TestCreate_SlugAlreadyExists(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	teamRepo := &mocks.MockTeamRepository{
		ExistsBySlugFunc: func(ctx context.Context, slug string) (bool, error) {
			return true, nil
		},
	}

	teamMemberRepo := &mocks.MockTeamMemberRepository{}
	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	input := &entity.TeamCreate{
		Name: "Test Team",
		Slug: "existing-slug",
	}

	_, err := uc.Create(ctx, userID, input)
	if err != team.ErrSlugAlreadyExists {
		t.Errorf("expected ErrSlugAlreadyExists, got %v", err)
	}
}

func TestGetByID_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	teamID := uuid.New()

	testTeam := &entity.Team{
		ID:        teamID,
		Name:      "Test Team",
		Slug:      "test-team",
		OwnerID:   userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	teamRepo := &mocks.MockTeamRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
			if id == teamID {
				return testTeam, nil
			}
			return nil, nil
		},
	}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			if tmID == teamID && uID == userID {
				return &entity.TeamMember{
					TeamID: tmID,
					UserID: uID,
					Role:   entity.TeamRoleOwner,
				}, nil
			}
			return nil, nil
		},
	}

	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	result, err := uc.GetByID(ctx, userID, teamID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected team, got nil")
	}

	if result.ID != teamID {
		t.Errorf("expected ID %s, got %s", teamID, result.ID)
	}
}

func TestGetByID_NotTeamMember(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	teamID := uuid.New()

	teamRepo := &mocks.MockTeamRepository{}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			return nil, nil
		},
	}

	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	_, err := uc.GetByID(ctx, userID, teamID)
	if err != team.ErrNotTeamMember {
		t.Errorf("expected ErrNotTeamMember, got %v", err)
	}
}

func TestUpdate_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	teamID := uuid.New()

	testTeam := &entity.Team{
		ID:        teamID,
		Name:      "Old Name",
		Slug:      "test-team",
		OwnerID:   userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	teamRepo := &mocks.MockTeamRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
			if id == teamID {
				return testTeam, nil
			}
			return nil, nil
		},
		UpdateFunc: func(ctx context.Context, tm *entity.Team) error {
			return nil
		},
	}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			if tmID == teamID && uID == userID {
				return &entity.TeamMember{
					TeamID: tmID,
					UserID: uID,
					Role:   entity.TeamRoleOwner,
				}, nil
			}
			return nil, nil
		},
	}

	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	newName := "New Name"
	input := &entity.TeamUpdate{
		Name: &newName,
	}

	result, err := uc.Update(ctx, userID, teamID, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != newName {
		t.Errorf("expected name %s, got %s", newName, result.Name)
	}
}

func TestUpdate_NotAdmin(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	teamID := uuid.New()

	teamRepo := &mocks.MockTeamRepository{}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			return &entity.TeamMember{
				TeamID: tmID,
				UserID: uID,
				Role:   entity.TeamRoleMember,
			}, nil
		},
	}

	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	newName := "New Name"
	input := &entity.TeamUpdate{
		Name: &newName,
	}

	_, err := uc.Update(ctx, userID, teamID, input)
	if err != team.ErrNotTeamAdmin {
		t.Errorf("expected ErrNotTeamAdmin, got %v", err)
	}
}

func TestDelete_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	teamID := uuid.New()

	teamRepo := &mocks.MockTeamRepository{
		DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
			return nil
		},
	}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			if tmID == teamID && uID == userID {
				return &entity.TeamMember{
					TeamID: tmID,
					UserID: uID,
					Role:   entity.TeamRoleOwner,
				}, nil
			}
			return nil, nil
		},
	}

	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	err := uc.Delete(ctx, userID, teamID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete_NotOwner(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	teamID := uuid.New()

	teamRepo := &mocks.MockTeamRepository{}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			return &entity.TeamMember{
				TeamID: tmID,
				UserID: uID,
				Role:   entity.TeamRoleAdmin,
			}, nil
		},
	}

	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	err := uc.Delete(ctx, userID, teamID)
	if err != team.ErrNotTeamOwner {
		t.Errorf("expected ErrNotTeamOwner, got %v", err)
	}
}

func TestAddMember_Success(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()
	newUserID := uuid.New()
	teamID := uuid.New()

	newUser := &entity.User{
		ID:    newUserID,
		Email: "newuser@example.com",
		Name:  "New User",
	}

	teamRepo := &mocks.MockTeamRepository{}

	callCount := 0
	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			callCount++
			if callCount == 1 && uID == ownerID {
				return &entity.TeamMember{
					TeamID: tmID,
					UserID: uID,
					Role:   entity.TeamRoleOwner,
				}, nil
			}
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, member *entity.TeamMember) error {
			if member.Role != entity.TeamRoleMember {
				t.Errorf("expected role %s, got %s", entity.TeamRoleMember, member.Role)
			}
			return nil
		},
	}

	userRepo := &mocks.MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			if email == "newuser@example.com" {
				return newUser, nil
			}
			return nil, nil
		},
	}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	input := &entity.TeamMemberCreate{
		Email: "newuser@example.com",
		Role:  entity.TeamRoleMember,
	}

	result, err := uc.AddMember(ctx, ownerID, teamID, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected team member, got nil")
	}

	if result.UserID != newUserID {
		t.Errorf("expected user ID %s, got %s", newUserID, result.UserID)
	}
}

func TestAddMember_UserNotFound(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()
	teamID := uuid.New()

	teamRepo := &mocks.MockTeamRepository{}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			return &entity.TeamMember{
				TeamID: tmID,
				UserID: uID,
				Role:   entity.TeamRoleOwner,
			}, nil
		},
	}

	userRepo := &mocks.MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return nil, nil
		},
	}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	input := &entity.TeamMemberCreate{
		Email: "nonexistent@example.com",
		Role:  entity.TeamRoleMember,
	}

	_, err := uc.AddMember(ctx, ownerID, teamID, input)
	if err != team.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestAddMember_AlreadyMember(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()
	existingUserID := uuid.New()
	teamID := uuid.New()

	existingUser := &entity.User{
		ID:    existingUserID,
		Email: "existing@example.com",
		Name:  "Existing User",
	}

	teamRepo := &mocks.MockTeamRepository{}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			return &entity.TeamMember{
				TeamID: tmID,
				UserID: uID,
				Role:   entity.TeamRoleOwner,
			}, nil
		},
	}

	userRepo := &mocks.MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return existingUser, nil
		},
	}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	input := &entity.TeamMemberCreate{
		Email: "existing@example.com",
		Role:  entity.TeamRoleMember,
	}

	_, err := uc.AddMember(ctx, ownerID, teamID, input)
	if err != team.ErrAlreadyMember {
		t.Errorf("expected ErrAlreadyMember, got %v", err)
	}
}

func TestRemoveMember_Success(t *testing.T) {
	ctx := context.Background()
	ownerID := uuid.New()
	memberID := uuid.New()
	teamID := uuid.New()

	teamRepo := &mocks.MockTeamRepository{}

	callCount := 0
	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			callCount++
			if callCount == 1 && uID == ownerID {
				return &entity.TeamMember{
					TeamID: tmID,
					UserID: uID,
					Role:   entity.TeamRoleOwner,
				}, nil
			}
			if uID == memberID {
				return &entity.TeamMember{
					TeamID: tmID,
					UserID: uID,
					Role:   entity.TeamRoleMember,
				}, nil
			}
			return nil, nil
		},
		DeleteFunc: func(ctx context.Context, tmID, uID uuid.UUID) error {
			return nil
		},
	}

	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	err := uc.RemoveMember(ctx, ownerID, teamID, memberID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoveMember_CannotRemoveOwner(t *testing.T) {
	ctx := context.Background()
	adminID := uuid.New()
	ownerID := uuid.New()
	teamID := uuid.New()

	teamRepo := &mocks.MockTeamRepository{}

	callCount := 0
	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			callCount++
			if callCount == 1 && uID == adminID {
				return &entity.TeamMember{
					TeamID: tmID,
					UserID: uID,
					Role:   entity.TeamRoleAdmin,
				}, nil
			}
			if uID == ownerID {
				return &entity.TeamMember{
					TeamID: tmID,
					UserID: uID,
					Role:   entity.TeamRoleOwner,
				}, nil
			}
			return nil, nil
		},
	}

	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	err := uc.RemoveMember(ctx, adminID, teamID, ownerID)
	if err != team.ErrCannotRemoveOwner {
		t.Errorf("expected ErrCannotRemoveOwner, got %v", err)
	}
}

func TestListMembers_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	teamID := uuid.New()

	members := []entity.TeamMember{
		{ID: uuid.New(), TeamID: teamID, UserID: uuid.New(), Role: entity.TeamRoleOwner},
		{ID: uuid.New(), TeamID: teamID, UserID: uuid.New(), Role: entity.TeamRoleMember},
	}

	teamRepo := &mocks.MockTeamRepository{}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		GetByTeamAndUserFunc: func(ctx context.Context, tmID, uID uuid.UUID) (*entity.TeamMember, error) {
			return &entity.TeamMember{
				TeamID: tmID,
				UserID: uID,
				Role:   entity.TeamRoleMember,
			}, nil
		},
		ListByTeamIDFunc: func(ctx context.Context, tmID uuid.UUID) ([]entity.TeamMember, error) {
			return members, nil
		},
	}

	userRepo := &mocks.MockUserRepository{}

	uc := team.NewUseCase(teamRepo, teamMemberRepo, userRepo)

	result, err := uc.ListMembers(ctx, userID, teamID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 members, got %d", len(result))
	}
}
