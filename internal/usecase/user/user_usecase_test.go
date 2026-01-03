package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/podoru/podoru/internal/domain/entity"
	"github.com/podoru/podoru/internal/mocks"
	"github.com/podoru/podoru/internal/usecase/user"
	"github.com/podoru/podoru/pkg/crypto"
)

func TestGetByID_Success(t *testing.T) {
	ctx := context.Background()
	testID := uuid.New()

	testUser := &entity.User{
		ID:        testID,
		Email:     "test@example.com",
		Name:      "Test User",
		Role:      entity.UserRoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userRepo := &mocks.MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			if id == testID {
				return testUser, nil
			}
			return nil, nil
		},
	}

	uc := user.NewUseCase(userRepo)

	result, err := uc.GetByID(ctx, testID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected user, got nil")
	}

	if result.ID != testID {
		t.Errorf("expected ID %s, got %s", testID, result.ID)
	}
}

func TestGetByID_UserNotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			return nil, nil
		},
	}

	uc := user.NewUseCase(userRepo)

	_, err := uc.GetByID(ctx, uuid.New())
	if err != user.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdate_Success(t *testing.T) {
	ctx := context.Background()
	testID := uuid.New()

	testUser := &entity.User{
		ID:        testID,
		Email:     "test@example.com",
		Name:      "Test User",
		Role:      entity.UserRoleUser,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userRepo := &mocks.MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			if id == testID {
				return testUser, nil
			}
			return nil, nil
		},
		UpdateFunc: func(ctx context.Context, u *entity.User) error {
			return nil
		},
	}

	uc := user.NewUseCase(userRepo)

	newName := "Updated Name"
	avatarURL := "https://example.com/avatar.png"
	input := &entity.UserUpdate{
		Name:      &newName,
		AvatarURL: &avatarURL,
	}

	result, err := uc.Update(ctx, testID, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != newName {
		t.Errorf("expected name %s, got %s", newName, result.Name)
	}

	if result.AvatarURL == nil || *result.AvatarURL != avatarURL {
		t.Errorf("expected avatar URL %s, got %v", avatarURL, result.AvatarURL)
	}
}

func TestUpdate_UserNotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			return nil, nil
		},
	}

	uc := user.NewUseCase(userRepo)

	newName := "Updated Name"
	input := &entity.UserUpdate{
		Name: &newName,
	}

	_, err := uc.Update(ctx, uuid.New(), input)
	if err != user.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdatePassword_Success(t *testing.T) {
	ctx := context.Background()
	testID := uuid.New()

	passwordHash, _ := crypto.HashPassword("currentpassword")

	testUser := &entity.User{
		ID:           testID,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Name:         "Test User",
		Role:         entity.UserRoleUser,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo := &mocks.MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			if id == testID {
				return testUser, nil
			}
			return nil, nil
		},
		UpdateFunc: func(ctx context.Context, u *entity.User) error {
			return nil
		},
	}

	uc := user.NewUseCase(userRepo)

	input := &entity.UserPasswordUpdate{
		CurrentPassword: "currentpassword",
		NewPassword:     "newpassword123",
	}

	err := uc.UpdatePassword(ctx, testID, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdatePassword_UserNotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			return nil, nil
		},
	}

	uc := user.NewUseCase(userRepo)

	input := &entity.UserPasswordUpdate{
		CurrentPassword: "currentpassword",
		NewPassword:     "newpassword123",
	}

	err := uc.UpdatePassword(ctx, uuid.New(), input)
	if err != user.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdatePassword_InvalidCurrentPassword(t *testing.T) {
	ctx := context.Background()
	testID := uuid.New()

	passwordHash, _ := crypto.HashPassword("currentpassword")

	testUser := &entity.User{
		ID:           testID,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Name:         "Test User",
		Role:         entity.UserRoleUser,
		IsActive:     true,
	}

	userRepo := &mocks.MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			if id == testID {
				return testUser, nil
			}
			return nil, nil
		},
	}

	uc := user.NewUseCase(userRepo)

	input := &entity.UserPasswordUpdate{
		CurrentPassword: "wrongpassword",
		NewPassword:     "newpassword123",
	}

	err := uc.UpdatePassword(ctx, testID, input)
	if err != user.ErrInvalidPassword {
		t.Errorf("expected ErrInvalidPassword, got %v", err)
	}
}

func TestUpdatePassword_SamePassword(t *testing.T) {
	ctx := context.Background()
	testID := uuid.New()

	passwordHash, _ := crypto.HashPassword("samepassword")

	testUser := &entity.User{
		ID:           testID,
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Name:         "Test User",
		Role:         entity.UserRoleUser,
		IsActive:     true,
	}

	userRepo := &mocks.MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entity.User, error) {
			if id == testID {
				return testUser, nil
			}
			return nil, nil
		},
	}

	uc := user.NewUseCase(userRepo)

	input := &entity.UserPasswordUpdate{
		CurrentPassword: "samepassword",
		NewPassword:     "samepassword",
	}

	err := uc.UpdatePassword(ctx, testID, input)
	if err != user.ErrSamePassword {
		t.Errorf("expected ErrSamePassword, got %v", err)
	}
}
