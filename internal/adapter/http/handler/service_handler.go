package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/podoru/podoru/internal/adapter/http/middleware"
	"github.com/podoru/podoru/internal/domain/entity"
	"github.com/podoru/podoru/internal/usecase/service"
	"github.com/podoru/podoru/pkg/response"
	"github.com/podoru/podoru/pkg/validator"
)

type ServiceHandler struct {
	serviceUseCase *service.UseCase
	validator      *validator.Validator
}

func NewServiceHandler(serviceUseCase *service.UseCase, validator *validator.Validator) *ServiceHandler {
	return &ServiceHandler{
		serviceUseCase: serviceUseCase,
		validator:      validator,
	}
}

func (h *ServiceHandler) ListByProject(c *gin.Context) {
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

	services, err := h.serviceUseCase.ListByProject(c.Request.Context(), userID, projectID)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			response.NotFound(c, "Project not found")
			return
		}
		if errors.Is(err, service.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		response.InternalError(c, "Failed to list services")
		return
	}

	response.Success(c, services)
}

func (h *ServiceHandler) Create(c *gin.Context) {
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

	var req entity.ServiceCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	s, err := h.serviceUseCase.Create(c.Request.Context(), userID, projectID, &req)
	if err != nil {
		if errors.Is(err, service.ErrProjectNotFound) {
			response.NotFound(c, "Project not found")
			return
		}
		if errors.Is(err, service.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		if errors.Is(err, service.ErrSlugAlreadyExists) {
			response.Conflict(c, "Slug already exists in this project")
			return
		}
		if errors.Is(err, service.ErrImageRequired) {
			response.BadRequest(c, "Image is required for image deploy type")
			return
		}
		response.InternalError(c, "Failed to create service")
		return
	}

	response.Created(c, s)
}

func (h *ServiceHandler) Get(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	serviceID, err := uuid.Parse(c.Param("serviceId"))
	if err != nil {
		response.BadRequest(c, "Invalid service ID")
		return
	}

	s, err := h.serviceUseCase.GetByID(c.Request.Context(), userID, serviceID)
	if err != nil {
		if errors.Is(err, service.ErrServiceNotFound) {
			response.NotFound(c, "Service not found")
			return
		}
		if errors.Is(err, service.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		response.InternalError(c, "Failed to get service")
		return
	}

	response.Success(c, s)
}

func (h *ServiceHandler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	serviceID, err := uuid.Parse(c.Param("serviceId"))
	if err != nil {
		response.BadRequest(c, "Invalid service ID")
		return
	}

	var req entity.ServiceUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	s, err := h.serviceUseCase.Update(c.Request.Context(), userID, serviceID, &req)
	if err != nil {
		if errors.Is(err, service.ErrServiceNotFound) {
			response.NotFound(c, "Service not found")
			return
		}
		if errors.Is(err, service.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		response.InternalError(c, "Failed to update service")
		return
	}

	response.Success(c, s)
}

func (h *ServiceHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	serviceID, err := uuid.Parse(c.Param("serviceId"))
	if err != nil {
		response.BadRequest(c, "Invalid service ID")
		return
	}

	if err := h.serviceUseCase.Delete(c.Request.Context(), userID, serviceID); err != nil {
		if errors.Is(err, service.ErrServiceNotFound) {
			response.NotFound(c, "Service not found")
			return
		}
		if errors.Is(err, service.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		if errors.Is(err, service.ErrNotTeamAdmin) {
			response.Forbidden(c, "Requires admin or owner role")
			return
		}
		response.InternalError(c, "Failed to delete service")
		return
	}

	response.NoContent(c)
}

func (h *ServiceHandler) Scale(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	serviceID, err := uuid.Parse(c.Param("serviceId"))
	if err != nil {
		response.BadRequest(c, "Invalid service ID")
		return
	}

	var req entity.ServiceScale
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	s, err := h.serviceUseCase.Scale(c.Request.Context(), userID, serviceID, req.Replicas)
	if err != nil {
		if errors.Is(err, service.ErrServiceNotFound) {
			response.NotFound(c, "Service not found")
			return
		}
		if errors.Is(err, service.ErrNotTeamMember) {
			response.Forbidden(c, "Not a team member")
			return
		}
		response.InternalError(c, "Failed to scale service")
		return
	}

	response.Success(c, s)
}

func (h *ServiceHandler) Deploy(c *gin.Context) {
	response.Success(c, gin.H{"message": "Deployment triggered"})
}

func (h *ServiceHandler) Start(c *gin.Context) {
	response.Success(c, gin.H{"message": "Service started"})
}

func (h *ServiceHandler) Stop(c *gin.Context) {
	response.Success(c, gin.H{"message": "Service stopped"})
}

func (h *ServiceHandler) Restart(c *gin.Context) {
	response.Success(c, gin.H{"message": "Service restarted"})
}

func (h *ServiceHandler) Logs(c *gin.Context) {
	response.Success(c, gin.H{"logs": []string{}})
}
