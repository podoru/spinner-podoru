package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
	"github.com/podoru/spinner-podoru/internal/infrastructure/config"
	"github.com/podoru/spinner-podoru/internal/mocks"
	"github.com/podoru/spinner-podoru/internal/usecase/auth"
	"github.com/podoru/spinner-podoru/pkg/crypto"
)

func TestRegister_FirstUser_SuperAdmin(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.MockUserRepository{
		CountFunc: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		ExistsByEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, user *entity.User) error {
			if user.Role != entity.UserRoleSuperAdmin {
				t.Errorf("expected role %s, got %s", entity.UserRoleSuperAdmin, user.Role)
			}
			return nil
		},
	}

	refreshTokenRepo := &mocks.MockRefreshTokenRepository{
		CreateFunc: func(ctx context.Context, token *entity.RefreshToken) error {
			return nil
		},
	}

	teamRepo := &mocks.MockTeamRepository{
		CreateFunc: func(ctx context.Context, team *entity.Team) error {
			return nil
		},
	}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		CreateFunc: func(ctx context.Context, member *entity.TeamMember) error {
			return nil
		},
	}

	jwtConfig := &config.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	appConfig := &config.AppConfig{
		RegistrationEnabled: false,
	}

	uc := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	input := &entity.UserCreate{
		Email:    "admin@example.com",
		Password: "password123",
		Name:     "Admin User",
	}

	user, tokens, err := uc.Register(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if tokens == nil {
		t.Fatal("expected tokens, got nil")
	}

	if user.Role != entity.UserRoleSuperAdmin {
		t.Errorf("expected role %s, got %s", entity.UserRoleSuperAdmin, user.Role)
	}
}

func TestRegister_SecondUser_RegistrationDisabled(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.MockUserRepository{
		CountFunc: func(ctx context.Context) (int64, error) {
			return 1, nil
		},
	}

	refreshTokenRepo := &mocks.MockRefreshTokenRepository{}
	teamRepo := &mocks.MockTeamRepository{}
	teamMemberRepo := &mocks.MockTeamMemberRepository{}

	jwtConfig := &config.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	appConfig := &config.AppConfig{
		RegistrationEnabled: false,
	}

	uc := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	input := &entity.UserCreate{
		Email:    "user@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	_, _, err := uc.Register(ctx, input)
	if err != auth.ErrRegistrationDisabled {
		t.Errorf("expected ErrRegistrationDisabled, got %v", err)
	}
}

func TestRegister_SecondUser_RegistrationEnabled(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.MockUserRepository{
		CountFunc: func(ctx context.Context) (int64, error) {
			return 1, nil
		},
		ExistsByEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, user *entity.User) error {
			if user.Role != entity.UserRoleUser {
				t.Errorf("expected role %s, got %s", entity.UserRoleUser, user.Role)
			}
			return nil
		},
	}

	refreshTokenRepo := &mocks.MockRefreshTokenRepository{
		CreateFunc: func(ctx context.Context, token *entity.RefreshToken) error {
			return nil
		},
	}

	teamRepo := &mocks.MockTeamRepository{
		CreateFunc: func(ctx context.Context, team *entity.Team) error {
			return nil
		},
	}

	teamMemberRepo := &mocks.MockTeamMemberRepository{
		CreateFunc: func(ctx context.Context, member *entity.TeamMember) error {
			return nil
		},
	}

	jwtConfig := &config.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	appConfig := &config.AppConfig{
		RegistrationEnabled: true,
	}

	uc := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	input := &entity.UserCreate{
		Email:    "user@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	user, tokens, err := uc.Register(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if tokens == nil {
		t.Fatal("expected tokens, got nil")
	}

	if user.Role != entity.UserRoleUser {
		t.Errorf("expected role %s, got %s", entity.UserRoleUser, user.Role)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.MockUserRepository{
		CountFunc: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		ExistsByEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return true, nil
		},
	}

	refreshTokenRepo := &mocks.MockRefreshTokenRepository{}
	teamRepo := &mocks.MockTeamRepository{}
	teamMemberRepo := &mocks.MockTeamMemberRepository{}

	jwtConfig := &config.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	appConfig := &config.AppConfig{
		RegistrationEnabled: true,
	}

	uc := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	input := &entity.UserCreate{
		Email:    "existing@example.com",
		Password: "password123",
		Name:     "Existing User",
	}

	_, _, err := uc.Register(ctx, input)
	if err != auth.ErrEmailAlreadyExists {
		t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	ctx := context.Background()

	passwordHash, _ := crypto.HashPassword("password123")

	testUser := &entity.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Name:         "Test User",
		Role:         entity.UserRoleUser,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	userRepo := &mocks.MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return testUser, nil
		},
	}

	refreshTokenRepo := &mocks.MockRefreshTokenRepository{
		CreateFunc: func(ctx context.Context, token *entity.RefreshToken) error {
			return nil
		},
	}

	teamRepo := &mocks.MockTeamRepository{}
	teamMemberRepo := &mocks.MockTeamMemberRepository{}

	jwtConfig := &config.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	appConfig := &config.AppConfig{}

	uc := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	input := &entity.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	user, tokens, err := uc.Login(ctx, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if tokens == nil {
		t.Fatal("expected tokens, got nil")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	ctx := context.Background()

	passwordHash, _ := crypto.HashPassword("password123")

	testUser := &entity.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Name:         "Test User",
		Role:         entity.UserRoleUser,
		IsActive:     true,
	}

	userRepo := &mocks.MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return testUser, nil
		},
	}

	refreshTokenRepo := &mocks.MockRefreshTokenRepository{}
	teamRepo := &mocks.MockTeamRepository{}
	teamMemberRepo := &mocks.MockTeamMemberRepository{}

	jwtConfig := &config.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	appConfig := &config.AppConfig{}

	uc := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	input := &entity.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, _, err := uc.Login(ctx, input)
	if err != auth.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := &mocks.MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return nil, nil
		},
	}

	refreshTokenRepo := &mocks.MockRefreshTokenRepository{}
	teamRepo := &mocks.MockTeamRepository{}
	teamMemberRepo := &mocks.MockTeamMemberRepository{}

	jwtConfig := &config.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	appConfig := &config.AppConfig{}

	uc := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	input := &entity.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, _, err := uc.Login(ctx, input)
	if err != auth.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UserInactive(t *testing.T) {
	ctx := context.Background()

	passwordHash, _ := crypto.HashPassword("password123")

	testUser := &entity.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		Name:         "Test User",
		Role:         entity.UserRoleUser,
		IsActive:     false,
	}

	userRepo := &mocks.MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return testUser, nil
		},
	}

	refreshTokenRepo := &mocks.MockRefreshTokenRepository{}
	teamRepo := &mocks.MockTeamRepository{}
	teamMemberRepo := &mocks.MockTeamMemberRepository{}

	jwtConfig := &config.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	appConfig := &config.AppConfig{}

	uc := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	input := &entity.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	_, _, err := uc.Login(ctx, input)
	if err != auth.ErrUserInactive {
		t.Errorf("expected ErrUserInactive, got %v", err)
	}
}

func TestValidateAccessToken_Success(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	refreshTokenRepo := &mocks.MockRefreshTokenRepository{}
	teamRepo := &mocks.MockTeamRepository{}
	teamMemberRepo := &mocks.MockTeamMemberRepository{}

	jwtConfig := &config.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	appConfig := &config.AppConfig{}

	uc := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	testUser := &entity.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	ctx := context.Background()
	passwordHash, _ := crypto.HashPassword("password123")
	testUser.PasswordHash = passwordHash
	testUser.Name = "Test User"
	testUser.Role = entity.UserRoleUser
	testUser.IsActive = true

	userRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entity.User, error) {
		return testUser, nil
	}
	refreshTokenRepo.CreateFunc = func(ctx context.Context, token *entity.RefreshToken) error {
		return nil
	}

	_, tokens, err := uc.Login(ctx, &entity.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	claims, err := uc.ValidateAccessToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if claims.UserID != testUser.ID {
		t.Errorf("expected user ID %s, got %s", testUser.ID, claims.UserID)
	}

	if claims.Email != testUser.Email {
		t.Errorf("expected email %s, got %s", testUser.Email, claims.Email)
	}
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	userRepo := &mocks.MockUserRepository{}
	refreshTokenRepo := &mocks.MockRefreshTokenRepository{}
	teamRepo := &mocks.MockTeamRepository{}
	teamMemberRepo := &mocks.MockTeamMemberRepository{}

	jwtConfig := &config.JWTConfig{
		Secret:        "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	appConfig := &config.AppConfig{}

	uc := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	_, err := uc.ValidateAccessToken("invalid-token")
	if err != auth.ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}
