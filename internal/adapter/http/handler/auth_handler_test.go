package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/podoru/podoru/internal/adapter/http/handler"
	"github.com/podoru/podoru/internal/domain/entity"
	"github.com/podoru/podoru/internal/infrastructure/config"
	"github.com/podoru/podoru/internal/mocks"
	"github.com/podoru/podoru/internal/usecase/auth"
	"github.com/podoru/podoru/pkg/validator"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthHandler(t *testing.T) (*handler.AuthHandler, *gin.Engine) {
	userRepo := &mocks.MockUserRepository{
		CountFunc: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
		ExistsByEmailFunc: func(ctx context.Context, email string) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, user *entity.User) error {
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

	authUseCase := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	v, err := validator.New()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	h := handler.NewAuthHandler(authUseCase, v)

	r := gin.New()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	return h, r
}

func TestAuthHandler_Register_Success(t *testing.T) {
	_, router := setupAuthHandler(t)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if response["success"] != true {
		t.Error("expected success to be true")
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data in response")
	}

	if _, ok := data["user"]; !ok {
		t.Error("expected user in response data")
	}

	if _, ok := data["tokens"]; !ok {
		t.Error("expected tokens in response data")
	}
}

func TestAuthHandler_Register_InvalidBody(t *testing.T) {
	_, router := setupAuthHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_Register_ValidationError(t *testing.T) {
	_, router := setupAuthHandler(t)

	body := map[string]string{
		"email":    "invalid-email",
		"password": "short",
		"name":     "",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestAuthHandler_Login_InvalidBody(t *testing.T) {
	_, router := setupAuthHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_Login_ValidationError(t *testing.T) {
	_, router := setupAuthHandler(t)

	body := map[string]string{
		"email":    "invalid",
		"password": "",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func setupLoginHandler(t *testing.T, testUser *entity.User) (*handler.AuthHandler, *gin.Engine) {
	userRepo := &mocks.MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			if email == testUser.Email {
				return testUser, nil
			}
			return nil, nil
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

	authUseCase := auth.NewUseCase(userRepo, refreshTokenRepo, teamRepo, teamMemberRepo, jwtConfig, appConfig)

	v, err := validator.New()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	h := handler.NewAuthHandler(authUseCase, v)

	r := gin.New()
	r.POST("/login", h.Login)

	return h, r
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	testUser := &entity.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		IsActive: true,
	}

	_, router := setupLoginHandler(t, testUser)

	body := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d: %s", http.StatusUnauthorized, w.Code, w.Body.String())
	}
}
