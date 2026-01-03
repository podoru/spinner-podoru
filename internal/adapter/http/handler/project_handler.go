package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/adapter/http/dto"
	"github.com/podoru/spinner-podoru/internal/adapter/http/middleware"
	"github.com/podoru/spinner-podoru/internal/domain/entity"
	"github.com/podoru/spinner-podoru/internal/usecase/project"
	"github.com/podoru/spinner-podoru/pkg/response"
	"github.com/podoru/spinner-podoru/pkg/validator"
)

// Ensure dto is used (for swagger)
var _ = dto.ProjectResponse{}

type ProjectHandler struct {
	projectUseCase *project.UseCase
	validator      *validator.Validator
}

func NewProjectHandler(projectUseCase *project.UseCase, validator *validator.Validator) *ProjectHandler {
	return &ProjectHandler{
		projectUseCase: projectUseCase,
		validator:      validator,
	}
}

// ListByTeam godoc
// @Summary      List projects
// @Description  Get all projects in a team
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        teamId path string true "Team ID" format(uuid)
// @Success      200 {object} response.Response{data=[]dto.ProjectResponse} "List of projects"
// @Failure      400 {object} response.Response "Invalid team ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams/{teamId}/projects [get]
func (h *ProjectHandler) ListByTeam(c *gin.Context) {
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

	projects, err := h.projectUseCase.ListByTeam(c.Request.Context(), userID, teamID)
	if err != nil {
		if errors.Is(err, project.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		response.InternalError(c, "Failed to list projects")
		return
	}

	response.Success(c, projects)
}

// Create godoc
// @Summary      Create project
// @Description  Create a new project in a team
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        teamId path string true "Team ID" format(uuid)
// @Param        request body dto.CreateProjectRequest true "Project data"
// @Success      201 {object} response.Response{data=dto.ProjectResponse} "Project created"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      409 {object} response.Response "Slug already exists in this team"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /teams/{teamId}/projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
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

	var req entity.ProjectCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	p, err := h.projectUseCase.Create(c.Request.Context(), userID, teamID, &req)
	if err != nil {
		if errors.Is(err, project.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		if errors.Is(err, project.ErrSlugAlreadyExists) {
			response.Conflict(c, "Slug already exists in this team")
			return
		}
		response.InternalError(c, "Failed to create project")
		return
	}

	response.Created(c, p)
}

// Get godoc
// @Summary      Get project
// @Description  Get project details by ID
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        projectId path string true "Project ID" format(uuid)
// @Success      200 {object} response.Response{data=dto.ProjectResponse} "Project details"
// @Failure      400 {object} response.Response "Invalid project ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Project not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /projects/{projectId} [get]
func (h *ProjectHandler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	p, err := h.projectUseCase.GetByID(c.Request.Context(), userID, projectID)
	if err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
			response.NotFound(c, "Project not found")
			return
		}
		if errors.Is(err, project.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		response.InternalError(c, "Failed to get project")
		return
	}

	response.Success(c, p)
}

// Update godoc
// @Summary      Update project
// @Description  Update project details. Requires admin or owner role.
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        projectId path string true "Project ID" format(uuid)
// @Param        request body dto.UpdateProjectRequest true "Project update data"
// @Success      200 {object} response.Response{data=dto.ProjectResponse} "Updated project"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member or insufficient permissions"
// @Failure      404 {object} response.Response "Project not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /projects/{projectId} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	var req entity.ProjectUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	p, err := h.projectUseCase.Update(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
			response.NotFound(c, "Project not found")
			return
		}
		if errors.Is(err, project.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		if errors.Is(err, project.ErrNotTeamAdmin) {
			response.Forbidden(c, "Requires admin or owner role")
			return
		}
		response.InternalError(c, "Failed to update project")
		return
	}

	response.Success(c, p)
}

// Delete godoc
// @Summary      Delete project
// @Description  Delete a project. Requires admin or owner role.
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        projectId path string true "Project ID" format(uuid)
// @Success      204 "Project deleted"
// @Failure      400 {object} response.Response "Invalid project ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member or insufficient permissions"
// @Failure      404 {object} response.Response "Project not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /projects/{projectId} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	projectID, err := uuid.Parse(c.Param("projectId"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	if err := h.projectUseCase.Delete(c.Request.Context(), userID, projectID); err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
			response.NotFound(c, "Project not found")
			return
		}
		if errors.Is(err, project.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		if errors.Is(err, project.ErrNotTeamAdmin) {
			response.Forbidden(c, "Requires admin or owner role")
			return
		}
		response.InternalError(c, "Failed to delete project")
		return
	}

	response.NoContent(c)
}

func (h *ProjectHandler) Deploy(c *gin.Context) {
	response.Success(c, gin.H{"message": "Deployment triggered"})
}
