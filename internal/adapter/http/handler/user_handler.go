package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/podoru/spinner-podoru/internal/adapter/http/dto"
	"github.com/podoru/spinner-podoru/internal/adapter/http/middleware"
	"github.com/podoru/spinner-podoru/internal/domain/entity"
	"github.com/podoru/spinner-podoru/internal/usecase/user"
	"github.com/podoru/spinner-podoru/pkg/response"
	"github.com/podoru/spinner-podoru/pkg/validator"
)

type UserHandler struct {
	userUseCase *user.UseCase
	validator   *validator.Validator
}

func NewUserHandler(userUseCase *user.UseCase, validator *validator.Validator) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		validator:   validator,
	}
}

// GetMe godoc
// @Summary      Get current user
// @Description  Get the profile of the currently authenticated user
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=dto.UserResponse} "User profile"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      404 {object} response.Response "User not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	u, err := h.userUseCase.GetByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalError(c, "Failed to get user")
		return
	}

	response.Success(c, dto.ToUserResponse(u))
}

// UpdateMe godoc
// @Summary      Update current user
// @Description  Update the profile of the currently authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object{name=string,avatar_url=string} false "User update data"
// @Success      200 {object} response.Response{data=dto.UserResponse} "Updated user profile"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      404 {object} response.Response "User not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req struct {
		Name      *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
		AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	input := &entity.UserUpdate{
		Name:      req.Name,
		AvatarURL: req.AvatarURL,
	}

	u, err := h.userUseCase.Update(c.Request.Context(), userID, input)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalError(c, "Failed to update user")
		return
	}

	response.Success(c, dto.ToUserResponse(u))
}

// UpdatePassword godoc
// @Summary      Update password
// @Description  Change the password of the currently authenticated user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object{current_password=string,new_password=string} true "Password change data"
// @Success      204 "Password updated successfully"
// @Failure      400 {object} response.Response "Invalid request or current password incorrect"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      404 {object} response.Response "User not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /users/me/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	input := &entity.UserPasswordUpdate{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	err := h.userUseCase.UpdatePassword(c.Request.Context(), userID, input)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		if errors.Is(err, user.ErrInvalidPassword) {
			response.BadRequest(c, "Current password is incorrect")
			return
		}
		if errors.Is(err, user.ErrSamePassword) {
			response.BadRequest(c, "New password must be different from current password")
			return
		}
		response.InternalError(c, "Failed to update password")
		return
	}

	response.NoContent(c)
}
