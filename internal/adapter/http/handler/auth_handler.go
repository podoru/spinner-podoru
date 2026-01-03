package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/podoru/spinner-podoru/internal/adapter/http/dto"
	"github.com/podoru/spinner-podoru/internal/usecase/auth"
	"github.com/podoru/spinner-podoru/pkg/response"
	"github.com/podoru/spinner-podoru/pkg/validator"
)

type AuthHandler struct {
	authUseCase *auth.UseCase
	validator   *validator.Validator
}

func NewAuthHandler(authUseCase *auth.UseCase, validator *validator.Validator) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		validator:   validator,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user account. First user becomes superadmin.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Registration details"
// @Success      201 {object} response.Response{data=dto.AuthResponse} "User registered successfully"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      403 {object} response.Response "Registration is disabled"
// @Failure      409 {object} response.Response "Email already exists"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	user, tokens, err := h.authUseCase.Register(c.Request.Context(), req.ToEntity())
	if err != nil {
		if errors.Is(err, auth.ErrEmailAlreadyExists) {
			response.Conflict(c, "Email already exists")
			return
		}
		if errors.Is(err, auth.ErrRegistrationDisabled) {
			response.Forbidden(c, "Registration is disabled")
			return
		}
		response.InternalError(c, "Failed to register user")
		return
	}

	response.Created(c, dto.ToAuthResponse(user, tokens))
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200 {object} response.Response{data=dto.AuthResponse} "Login successful"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "Invalid email or password"
// @Failure      403 {object} response.Response "Account is inactive"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	user, tokens, err := h.authUseCase.Login(c.Request.Context(), req.ToEntity())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			response.Unauthorized(c, "Invalid email or password")
			return
		}
		if errors.Is(err, auth.ErrUserInactive) {
			response.Forbidden(c, "Account is inactive")
			return
		}
		response.InternalError(c, "Failed to login")
		return
	}

	response.Success(c, dto.ToAuthResponse(user, tokens))
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Get new access token using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshRequest true "Refresh token"
// @Success      200 {object} response.Response{data=dto.TokenResponse} "Token refreshed successfully"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "Invalid or expired refresh token"
// @Failure      403 {object} response.Response "Account is inactive"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	tokens, err := h.authUseCase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrTokenExpired) {
			response.Unauthorized(c, "Invalid or expired refresh token")
			return
		}
		if errors.Is(err, auth.ErrUserInactive) {
			response.Forbidden(c, "Account is inactive")
			return
		}
		response.InternalError(c, "Failed to refresh token")
		return
	}

	response.Success(c, dto.ToTokenResponse(tokens))
}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshRequest true "Refresh token to invalidate"
// @Success      204 "Logout successful"
// @Failure      400 {object} response.Response "Invalid request body"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.authUseCase.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.InternalError(c, "Failed to logout")
		return
	}

	response.NoContent(c)
}
