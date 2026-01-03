package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/adapter/http/dto"
	"github.com/podoru/spinner-podoru/internal/adapter/http/middleware"
	"github.com/podoru/spinner-podoru/internal/domain/entity"
	"github.com/podoru/spinner-podoru/internal/usecase/deployment"
	"github.com/podoru/spinner-podoru/internal/usecase/service"
	"github.com/podoru/spinner-podoru/pkg/response"
	"github.com/podoru/spinner-podoru/pkg/validator"
)

// Ensure dto is used (for swagger)
var _ = dto.ServiceResponse{}

type ServiceHandler struct {
	serviceUseCase    *service.UseCase
	deploymentUseCase *deployment.UseCase
	validator         *validator.Validator
}

func NewServiceHandler(serviceUseCase *service.UseCase, deploymentUseCase *deployment.UseCase, validator *validator.Validator) *ServiceHandler {
	return &ServiceHandler{
		serviceUseCase:    serviceUseCase,
		deploymentUseCase: deploymentUseCase,
		validator:         validator,
	}
}

// ListByProject godoc
// @Summary      List services
// @Description  Get all services in a project
// @Tags         services
// @Produce      json
// @Security     BearerAuth
// @Param        projectId path string true "Project ID" format(uuid)
// @Success      200 {object} response.Response{data=[]dto.ServiceResponse} "List of services"
// @Failure      400 {object} response.Response "Invalid project ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Project not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /projects/{projectId}/services [get]
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

// Create godoc
// @Summary      Create service
// @Description  Create a new service in a project
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        projectId path string true "Project ID" format(uuid)
// @Param        request body dto.CreateServiceRequest true "Service data"
// @Success      201 {object} response.Response{data=dto.ServiceResponse} "Service created"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Project not found"
// @Failure      409 {object} response.Response "Slug already exists in this project"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /projects/{projectId}/services [post]
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

// Get godoc
// @Summary      Get service
// @Description  Get service details by ID
// @Tags         services
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Success      200 {object} response.Response{data=dto.ServiceResponse} "Service details"
// @Failure      400 {object} response.Response "Invalid service ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId} [get]
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

// Update godoc
// @Summary      Update service
// @Description  Update service details
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Param        request body dto.UpdateServiceRequest true "Service update data"
// @Success      200 {object} response.Response{data=dto.ServiceResponse} "Updated service"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId} [put]
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

// Delete godoc
// @Summary      Delete service
// @Description  Delete a service. Requires admin or owner role.
// @Tags         services
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Success      204 "Service deleted"
// @Failure      400 {object} response.Response "Invalid service ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member or insufficient permissions"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId} [delete]
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

	// First, destroy the container if it exists
	if err := h.deploymentUseCase.Destroy(c.Request.Context(), userID, serviceID); err != nil {
		switch {
		case errors.Is(err, deployment.ErrServiceNotFound):
			response.NotFound(c, "Service not found")
		case errors.Is(err, deployment.ErrNotTeamMember):
			response.Forbidden(c, "Not a team member")
		default:
			response.InternalError(c, "Failed to destroy container")
		}
		return
	}

	// Then delete from database
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

// Scale godoc
// @Summary      Scale service
// @Description  Scale a service to a specified number of replicas
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Param        request body dto.ScaleServiceRequest true "Scale data"
// @Success      200 {object} response.Response{data=dto.ServiceResponse} "Scaled service"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId}/scale [post]
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

// Deploy godoc
// @Summary      Deploy service
// @Description  Trigger a deployment for the service
// @Tags         services
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Success      200 {object} response.Response{data=dto.DeploymentResponse} "Deployment triggered"
// @Failure      400 {object} response.Response "Invalid service ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      409 {object} response.Response "Deployment already in progress"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId}/deploy [post]
func (h *ServiceHandler) Deploy(c *gin.Context) {
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

	dep, err := h.deploymentUseCase.Deploy(c.Request.Context(), userID, serviceID)
	if err != nil {
		switch {
		case errors.Is(err, deployment.ErrServiceNotFound):
			response.NotFound(c, "Service not found")
		case errors.Is(err, deployment.ErrNotTeamMember):
			response.Forbidden(c, "Not a team member")
		case errors.Is(err, deployment.ErrAlreadyDeploying):
			response.Conflict(c, "Deployment already in progress")
		case errors.Is(err, deployment.ErrNoImageSpecified):
			response.BadRequest(c, "No image specified for deployment")
		default:
			response.InternalError(c, "Failed to deploy service")
		}
		return
	}

	response.Success(c, dto.DeploymentResponse{
		ID:        dep.ID,
		ServiceID: dep.ServiceID,
		Status:    string(dep.Status),
		StartedAt: dep.StartedAt,
	})
}

// Start godoc
// @Summary      Start service
// @Description  Start a stopped service
// @Tags         services
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Success      200 {object} response.Response{data=dto.MessageResponse} "Service started"
// @Failure      400 {object} response.Response "Invalid service ID or service not deployed"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId}/start [post]
func (h *ServiceHandler) Start(c *gin.Context) {
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

	if err := h.deploymentUseCase.Start(c.Request.Context(), userID, serviceID); err != nil {
		switch {
		case errors.Is(err, deployment.ErrServiceNotFound):
			response.NotFound(c, "Service not found")
		case errors.Is(err, deployment.ErrNotTeamMember):
			response.Forbidden(c, "Not a team member")
		case errors.Is(err, deployment.ErrServiceNotDeployed):
			response.BadRequest(c, "Service not deployed yet")
		default:
			response.InternalError(c, "Failed to start service")
		}
		return
	}

	response.Success(c, dto.MessageResponse{Message: "Service started"})
}

// Stop godoc
// @Summary      Stop service
// @Description  Stop a running service
// @Tags         services
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Success      200 {object} response.Response{data=dto.MessageResponse} "Service stopped"
// @Failure      400 {object} response.Response "Invalid service ID or service not deployed"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId}/stop [post]
func (h *ServiceHandler) Stop(c *gin.Context) {
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

	if err := h.deploymentUseCase.Stop(c.Request.Context(), userID, serviceID); err != nil {
		switch {
		case errors.Is(err, deployment.ErrServiceNotFound):
			response.NotFound(c, "Service not found")
		case errors.Is(err, deployment.ErrNotTeamMember):
			response.Forbidden(c, "Not a team member")
		case errors.Is(err, deployment.ErrServiceNotDeployed):
			response.BadRequest(c, "Service not deployed yet")
		default:
			response.InternalError(c, "Failed to stop service")
		}
		return
	}

	response.Success(c, dto.MessageResponse{Message: "Service stopped"})
}

// Restart godoc
// @Summary      Restart service
// @Description  Restart a running service
// @Tags         services
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Success      200 {object} response.Response{data=dto.MessageResponse} "Service restarted"
// @Failure      400 {object} response.Response "Invalid service ID or service not deployed"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId}/restart [post]
func (h *ServiceHandler) Restart(c *gin.Context) {
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

	if err := h.deploymentUseCase.Restart(c.Request.Context(), userID, serviceID); err != nil {
		switch {
		case errors.Is(err, deployment.ErrServiceNotFound):
			response.NotFound(c, "Service not found")
		case errors.Is(err, deployment.ErrNotTeamMember):
			response.Forbidden(c, "Not a team member")
		case errors.Is(err, deployment.ErrServiceNotDeployed):
			response.BadRequest(c, "Service not deployed yet")
		default:
			response.InternalError(c, "Failed to restart service")
		}
		return
	}

	response.Success(c, dto.MessageResponse{Message: "Service restarted"})
}

// Logs godoc
// @Summary      Get service logs
// @Description  Get logs from a service
// @Tags         services
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Param        tail query int false "Number of lines to return" default(100)
// @Param        since query string false "Return logs since this timestamp (RFC3339)" example("2024-01-15T10:30:00Z")
// @Success      200 {object} response.Response{data=dto.ServiceLogsResponse} "Service logs"
// @Failure      400 {object} response.Response "Invalid service ID or service not deployed"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId}/logs [get]
func (h *ServiceHandler) Logs(c *gin.Context) {
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

	tail := c.DefaultQuery("tail", "100")
	since := c.Query("since")

	logs, err := h.deploymentUseCase.GetLogs(c.Request.Context(), userID, serviceID, tail, since)
	if err != nil {
		switch {
		case errors.Is(err, deployment.ErrServiceNotFound):
			response.NotFound(c, "Service not found")
		case errors.Is(err, deployment.ErrNotTeamMember):
			response.Forbidden(c, "Not a team member")
		case errors.Is(err, deployment.ErrServiceNotDeployed):
			response.BadRequest(c, "Service not deployed yet")
		default:
			response.InternalError(c, "Failed to get logs")
		}
		return
	}

	response.Success(c, dto.ServiceLogsResponse{
		ServiceID: serviceID,
		Logs:      logs,
		Timestamp: time.Now(),
	})
}

// ListDomains godoc
// @Summary      List service domains
// @Description  Get all domains for a service
// @Tags         services
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Success      200 {object} response.Response{data=[]dto.DomainResponse} "List of domains"
// @Failure      400 {object} response.Response "Invalid service ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId}/domains [get]
func (h *ServiceHandler) ListDomains(c *gin.Context) {
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

	domains, err := h.serviceUseCase.ListDomains(c.Request.Context(), userID, serviceID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrServiceNotFound):
			response.NotFound(c, "Service not found")
		case errors.Is(err, service.ErrNotTeamMember):
			response.Forbidden(c, "Not a team member")
		default:
			response.InternalError(c, "Failed to list domains")
		}
		return
	}

	response.Success(c, dto.ToDomainsResponse(domains))
}

// AddDomain godoc
// @Summary      Add domain to service
// @Description  Add a domain to a service for Traefik routing
// @Tags         services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Param        request body dto.CreateDomainRequest true "Domain data"
// @Success      201 {object} response.Response{data=dto.DomainResponse} "Domain added"
// @Failure      400 {object} response.Response "Invalid request body or validation error"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service not found"
// @Failure      409 {object} response.Response "Domain already in use"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId}/domains [post]
func (h *ServiceHandler) AddDomain(c *gin.Context) {
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

	var req entity.DomainCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		response.ValidationError(c, validator.FormatValidationErrors(err))
		return
	}

	domain, err := h.serviceUseCase.AddDomain(c.Request.Context(), userID, serviceID, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrServiceNotFound):
			response.NotFound(c, "Service not found")
		case errors.Is(err, service.ErrNotTeamMember):
			response.Forbidden(c, "Not a team member")
		case errors.Is(err, service.ErrDomainAlreadyInUse):
			response.Conflict(c, "Domain already in use")
		default:
			response.InternalError(c, "Failed to add domain")
		}
		return
	}

	response.Created(c, dto.ToDomainResponse(domain))
}

// DeleteDomain godoc
// @Summary      Delete domain from service
// @Description  Remove a domain from a service
// @Tags         services
// @Produce      json
// @Security     BearerAuth
// @Param        serviceId path string true "Service ID" format(uuid)
// @Param        domainId path string true "Domain ID" format(uuid)
// @Success      204 "Domain deleted"
// @Failure      400 {object} response.Response "Invalid ID"
// @Failure      401 {object} response.Response "User not authenticated"
// @Failure      403 {object} response.Response "Not a team member"
// @Failure      404 {object} response.Response "Service or domain not found"
// @Failure      500 {object} response.Response "Internal server error"
// @Router       /services/{serviceId}/domains/{domainId} [delete]
func (h *ServiceHandler) DeleteDomain(c *gin.Context) {
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

	domainID, err := uuid.Parse(c.Param("domainId"))
	if err != nil {
		response.BadRequest(c, "Invalid domain ID")
		return
	}

	if err := h.serviceUseCase.DeleteDomain(c.Request.Context(), userID, serviceID, domainID); err != nil {
		switch {
		case errors.Is(err, service.ErrServiceNotFound):
			response.NotFound(c, "Service not found")
		case errors.Is(err, service.ErrNotTeamMember):
			response.Forbidden(c, "Not a team member")
		case errors.Is(err, service.ErrDomainNotFound):
			response.NotFound(c, "Domain not found")
		default:
			response.InternalError(c, "Failed to delete domain")
		}
		return
	}

	response.NoContent(c)
}
