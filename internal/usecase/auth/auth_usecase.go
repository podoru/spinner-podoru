package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/podoru/podoru/internal/domain/entity"
	"github.com/podoru/podoru/internal/domain/repository"
	"github.com/podoru/podoru/internal/infrastructure/config"
	"github.com/podoru/podoru/pkg/crypto"
)

var (
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrInvalidToken          = errors.New("invalid token")
	ErrTokenExpired          = errors.New("token expired")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserInactive          = errors.New("user account is inactive")
	ErrRegistrationDisabled  = errors.New("registration is disabled")
)

type UseCase struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	teamRepo         repository.TeamRepository
	teamMemberRepo   repository.TeamMemberRepository
	jwtConfig        *config.JWTConfig
	appConfig        *config.AppConfig
}

func NewUseCase(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	teamRepo repository.TeamRepository,
	teamMemberRepo repository.TeamMemberRepository,
	jwtConfig *config.JWTConfig,
	appConfig *config.AppConfig,
) *UseCase {
	return &UseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		teamRepo:         teamRepo,
		teamMemberRepo:   teamMemberRepo,
		jwtConfig:        jwtConfig,
		appConfig:        appConfig,
	}
}

func (uc *UseCase) Register(ctx context.Context, input *entity.UserCreate) (*entity.User, *entity.TokenPair, error) {
	userCount, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	isFirstUser := userCount == 0

	if !isFirstUser && !uc.appConfig.RegistrationEnabled {
		return nil, nil, ErrRegistrationDisabled
	}

	exists, err := uc.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, nil, err
	}
	if exists {
		return nil, nil, ErrEmailAlreadyExists
	}

	passwordHash, err := crypto.HashPassword(input.Password)
	if err != nil {
		return nil, nil, err
	}

	role := entity.UserRoleUser
	if isFirstUser {
		role = entity.UserRoleSuperAdmin
	}

	now := time.Now()
	user := &entity.User{
		ID:           uuid.New(),
		Email:        input.Email,
		PasswordHash: passwordHash,
		Name:         input.Name,
		Role:         role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	teamSlug, _ := crypto.GenerateRandomString(8)
	team := &entity.Team{
		ID:          uuid.New(),
		Name:        user.Name + "'s Team",
		Slug:        "team-" + teamSlug,
		OwnerID:     user.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.teamRepo.Create(ctx, team); err != nil {
		return nil, nil, err
	}

	teamMember := &entity.TeamMember{
		ID:        uuid.New(),
		TeamID:    team.ID,
		UserID:    user.ID,
		Role:      entity.TeamRoleOwner,
		CreatedAt: now,
	}

	if err := uc.teamMemberRepo.Create(ctx, teamMember); err != nil {
		return nil, nil, err
	}

	tokens, err := uc.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (uc *UseCase) Login(ctx context.Context, input *entity.LoginRequest) (*entity.User, *entity.TokenPair, error) {
	user, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, nil, ErrUserInactive
	}

	if !crypto.CheckPassword(input.Password, user.PasswordHash) {
		return nil, nil, ErrInvalidCredentials
	}

	tokens, err := uc.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (uc *UseCase) RefreshToken(ctx context.Context, refreshToken string) (*entity.TokenPair, error) {
	tokenHash := crypto.HashToken(refreshToken)

	storedToken, err := uc.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if storedToken == nil {
		return nil, ErrInvalidToken
	}

	if time.Now().After(storedToken.ExpiresAt) {
		uc.refreshTokenRepo.DeleteByTokenHash(ctx, tokenHash)
		return nil, ErrTokenExpired
	}

	user, err := uc.userRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if err := uc.refreshTokenRepo.DeleteByTokenHash(ctx, tokenHash); err != nil {
		return nil, err
	}

	return uc.generateTokenPair(ctx, user)
}

func (uc *UseCase) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := crypto.HashToken(refreshToken)
	return uc.refreshTokenRepo.DeleteByTokenHash(ctx, tokenHash)
}

func (uc *UseCase) ValidateAccessToken(tokenString string) (*entity.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(uc.jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	email, _ := claims["email"].(string)

	return &entity.JWTClaims{
		UserID: userID,
		Email:  email,
	}, nil
}

func (uc *UseCase) generateTokenPair(ctx context.Context, user *entity.User) (*entity.TokenPair, error) {
	accessToken, err := uc.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := crypto.GenerateRandomString(64)
	if err != nil {
		return nil, err
	}

	tokenHash := crypto.HashToken(refreshToken)
	expiresAt := time.Now().Add(uc.jwtConfig.RefreshExpiry)

	refreshTokenEntity := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := uc.refreshTokenRepo.Create(ctx, refreshTokenEntity); err != nil {
		return nil, err
	}

	return &entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(uc.jwtConfig.AccessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (uc *UseCase) generateAccessToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(uc.jwtConfig.AccessExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtConfig.Secret))
}
