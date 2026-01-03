package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/adapter/http/dto"
	"github.com/podoru/spinner-podoru/internal/adapter/http/middleware"
	"github.com/podoru/spinner-podoru/internal/domain/entity"
	"github.com/podoru/spinner-podoru/internal/usecase/team"
	"github.com/podoru/spinner-podoru/pkg/response"
	"github.com/podoru/spinner-podoru/pkg/validator"
)

// Ensure dto is used (for swagger)
var _ = dto.TeamResponse{}

type TeamHandler struct {
	teamUseCase *team.UseCase
	validator   *validator.Validator
}

func NewTeamHandler(teamUseCase *team.UseCase, validator *validator.Validator) *TeamHandler {
	return &TeamHandler{
		teamUseCase: teamUseCase,
		validator:   validator,
	}
}

// List godoc
// @Summary      List teams
// @Description  Get all teams the current user is a member of
// @Tags         teams
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=[]dto.TeamWithRoleResponse} "List of teams"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams [get]
func (h *TeamHandler) List(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	teams, err := h.teamUseCase.List(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, "Failed to list teams")
		return
	}

	response.Success(c, teams)
}

// Create godoc
// @Summary      Create team
// @Description  Create a new team. The creator becomes the owner.
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateTeamRequest true "Team data"
// @Success      201 {object} response.Response{data=dto.TeamResponse} "Team created"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      409 {object} response.Response "Slug already exists"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams [post]
func (h *TeamHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req entity.TeamCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	t, err := h.teamUseCase.Create(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, team.ErrSlugAlreadyExists) {
			response.Conflict(c, "Slug already exists")
			return
		}
		response.InternalError(c, "Failed to create team")
		return
	}

	response.Created(c, t)
}

// Get godoc
// @Summary      Get team
// @Description  Get team details by ID
// @Tags         teams
// @Produce      json
// @Security     BearerAuth
// @Param        teamId path string true "Team ID" format(uuid)
// @Success      200 {object} response.Response{data=dto.TeamResponse} "Team details"
// @Failure      400 {object} response.Response "Invalid team ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Team not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams/{teamId} [get]
func (h *TeamHandler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		response.BadRequest(c, "Invalid team ID")
		return
	}

	t, err := h.teamUseCase.GetByID(c.Request.Context(), userID, teamID)
	if err != nil {
		if errors.Is(err, team.ErrTeamNotFound) {
			response.NotFound(c, "Team not found")
			return
		}
		if errors.Is(err, team.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		response.InternalError(c, "Failed to get team")
		return
	}

	response.Success(c, t)
}

// Update godoc
// @Summary      Update team
// @Description  Update team details. Requires admin or owner role.
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        teamId path string true "Team ID" format(uuid)
// @Param        request body dto.UpdateTeamRequest true "Team update data"
// @Success      200 {object} response.Response{data=dto.TeamResponse} "Updated team"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member or insufficient permissions"
// @Failure      404 {object} response.Response "Team not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams/{teamId} [put]
func (h *TeamHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		response.BadRequest(c, "Invalid team ID")
		return
	}

	var req entity.TeamUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	t, err := h.teamUseCase.Update(c.Request.Context(), userID, teamID, &req)
	if err != nil {
		if errors.Is(err, team.ErrTeamNotFound) {
			response.NotFound(c, "Team not found")
			return
		}
		if errors.Is(err, team.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		if errors.Is(err, team.ErrNotTeamAdmin) {
			response.Forbidden(c, "Requires admin or owner role")
			return
		}
		response.InternalError(c, "Failed to update team")
		return
	}

	response.Success(c, t)
}

// Delete godoc
// @Summary      Delete team
// @Description  Delete a team. Only the owner can delete a team.
// @Tags         teams
// @Produce      json
// @Security     BearerAuth
// @Param        teamId path string true "Team ID" format(uuid)
// @Success      204 "Team deleted"
// @Failure      400 {object} response.Response "Invalid team ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not team owner"
// @Failure      404 {object} response.Response "Team not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams/{teamId} [delete]
func (h *TeamHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		response.BadRequest(c, "Invalid team ID")
		return
	}

	if err := h.teamUseCase.Delete(c.Request.Context(), userID, teamID); err != nil {
		if errors.Is(err, team.ErrTeamNotFound) {
			response.NotFound(c, "Team not found")
			return
		}
		if errors.Is(err, team.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		if errors.Is(err, team.ErrNotTeamOwner) {
			response.Forbidden(c, "Only team owner can delete team")
			return
		}
		response.InternalError(c, "Failed to delete team")
		return
	}

	response.NoContent(c)
}

// ListMembers godoc
// @Summary      List team members
// @Description  Get all members of a team
// @Tags         teams
// @Produce      json
// @Security     BearerAuth
// @Param        teamId path string true "Team ID" format(uuid)
// @Success      200 {object} response.Response{data=[]dto.TeamMemberResponse} "List of team members"
// @Failure      400 {object} response.Response "Invalid team ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams/{teamId}/members [get]
func (h *TeamHandler) ListMembers(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		response.BadRequest(c, "Invalid team ID")
		return
	}

	members, err := h.teamUseCase.ListMembers(c.Request.Context(), userID, teamID)
	if err != nil {
		if errors.Is(err, team.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		response.InternalError(c, "Failed to list members")
		return
	}

	response.Success(c, members)
}

// AddMember godoc
// @Summary      Add team member
// @Description  Add a new member to the team. Requires admin or owner role.
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        teamId path string true "Team ID" format(uuid)
// @Param        request body dto.AddTeamMemberRequest true "Member data"
// @Success      201 {object} response.Response{data=dto.TeamMemberResponse} "Member added"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member or insufficient permissions"
// @Failure      404 {object} response.Response "User not found"
// @Failure      409 {object} response.Response "User is already a team member"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams/{teamId}/members [post]
func (h *TeamHandler) AddMember(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		response.BadRequest(c, "Invalid team ID")
		return
	}

	var req entity.TeamMemberCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	member, err := h.teamUseCase.AddMember(c.Request.Context(), userID, teamID, &req)
	if err != nil {
		if errors.Is(err, team.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		if errors.Is(err, team.ErrNotTeamAdmin) {
			response.Forbidden(c, "Requires admin or owner role")
			return
		}
		if errors.Is(err, team.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		if errors.Is(err, team.ErrAlreadyMember) {
			response.Conflict(c, "User is already a team member")
			return
		}
		response.InternalError(c, "Failed to add member")
		return
	}

	response.Created(c, member)
}

// UpdateMember godoc
// @Summary      Update team member
// @Description  Update a team member's role. Requires admin or owner role.
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        teamId path string true "Team ID" format(uuid)
// @Param        userId path string true "User ID" format(uuid)
// @Param        request body dto.UpdateTeamMemberRequest true "Member update data"
// @Success      200 {object} response.Response{data=dto.TeamMemberResponse} "Updated member"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member, insufficient permissions, or cannot modify owner"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams/{teamId}/members/{userId} [put]
func (h *TeamHandler) UpdateMember(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		response.BadRequest(c, "Invalid team ID")
		return
	}

	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req entity.TeamMemberUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	member, err := h.teamUseCase.UpdateMember(c.Request.Context(), userID, teamID, targetUserID, &req)
	if err != nil {
		if errors.Is(err, team.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		if errors.Is(err, team.ErrNotTeamAdmin) {
			response.Forbidden(c, "Requires admin or owner role")
			return
		}
		if errors.Is(err, team.ErrCannotRemoveOwner) {
			response.Forbidden(c, "Cannot modify team owner")
			return
		}
		response.InternalError(c, "Failed to update member")
		return
	}

	response.Success(c, member)
}

// RemoveMember godoc
// @Summary      Remove team member
// @Description  Remove a member from the team. Requires admin or owner role.
// @Tags         teams
// @Produce      json
// @Security     BearerAuth
// @Param        teamId path string true "Team ID" format(uuid)
// @Param        userId path string true "User ID" format(uuid)
// @Success      204 "Member removed"
// @Failure      400 {object} response.Response "Invalid team or user ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member, insufficient permissions, or cannot remove owner"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams/{teamId}/members/{userId} [delete]
func (h *TeamHandler) RemoveMember(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		response.BadRequest(c, "Invalid team ID")
		return
	}

	targetUserID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	if err := h.teamUseCase.RemoveMember(c.Request.Context(), userID, teamID, targetUserID); err != nil {
		if errors.Is(err, team.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		if errors.Is(err, team.ErrNotTeamAdmin) {
			response.Forbidden(c, "Requires admin or owner role")
			return
		}
		if errors.Is(err, team.ErrCannotRemoveOwner) {
			response.Forbidden(c, "Cannot remove team owner")
			return
		}
		response.InternalError(c, "Failed to remove member")
		return
	}

	response.NoContent(c)
}
