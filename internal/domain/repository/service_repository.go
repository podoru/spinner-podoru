package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *entity.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error)
	GetByProjectAndSlug(ctx context.Context, projectID uuid.UUID, slug string) (*entity.Service, error)
	Update(ctx context.Context, service *entity.Service) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]entity.Service, error)
	ExistsByProjectAndSlug(ctx context.Context, projectID uuid.UUID, slug string) (bool, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ServiceStatus) error
	UpdateContainerID(ctx context.Context, id uuid.UUID, containerID *string) error
	UpdateSwarmServiceID(ctx context.Context, id uuid.UUID, swarmServiceID *string) error
}

type DomainRepository interface {
	Create(ctx context.Context, domain *entity.Domain) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Domain, error)
	GetByDomain(ctx context.Context, domain string) (*entity.Domain, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByServiceID(ctx context.Context, serviceID uuid.UUID) ([]entity.Domain, error)
	ExistsByDomain(ctx context.Context, domain string) (bool, error)
}

type PortMappingRepository interface {
	Create(ctx context.Context, portMapping *entity.PortMapping) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByServiceID(ctx context.Context, serviceID uuid.UUID) ([]entity.PortMapping, error)
	DeleteByServiceID(ctx context.Context, serviceID uuid.UUID) error
}

type VolumeRepository interface {
	Create(ctx context.Context, volume *entity.Volume) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Volume, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByServiceID(ctx context.Context, serviceID uuid.UUID) ([]entity.Volume, error)
	DeleteByServiceID(ctx context.Context, serviceID uuid.UUID) error
}

type NetworkRepository interface {
	Create(ctx context.Context, network *entity.Network) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Network, error)
	GetByName(ctx context.Context, name string) (*entity.Network, error)
	Update(ctx context.Context, network *entity.Network) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]entity.Network, error)
	ListAll(ctx context.Context) ([]entity.Network, error)
}

type DeploymentRepository interface {
	Create(ctx context.Context, deployment *entity.Deployment) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Deployment, error)
	Update(ctx context.Context, deployment *entity.Deployment) error
	ListByServiceID(ctx context.Context, serviceID uuid.UUID, limit, offset int) ([]entity.Deployment, error)
	GetLatestByServiceID(ctx context.Context, serviceID uuid.UUID) (*entity.Deployment, error)
}
