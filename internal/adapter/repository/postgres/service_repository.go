package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

type ServiceRepository struct {
	pool *pgxpool.Pool
}

func NewServiceRepository(pool *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{pool: pool}
}

func (r *ServiceRepository) Create(ctx context.Context, service *entity.Service) error {
	query := `
		INSERT INTO services (id, project_id, name, slug, deploy_type, image, dockerfile_path,
			build_context, compose_file, env_vars_encrypted, replicas, cpu_limit, memory_limit,
			health_check_path, health_check_interval, restart_policy, status, container_id,
			swarm_service_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`
	_, err := r.pool.Exec(ctx, query,
		service.ID, service.ProjectID, service.Name, service.Slug, service.DeployType,
		service.Image, service.DockerfilePath, service.BuildContext, service.ComposeFile,
		service.EnvVarsEncrypted, service.Replicas, service.CPULimit, service.MemoryLimit,
		service.HealthCheckPath, service.HealthCheckInterval, service.RestartPolicy,
		service.Status, service.ContainerID, service.SwarmServiceID, service.CreatedAt, service.UpdatedAt,
	)
	return err
}

func (r *ServiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	query := `
		SELECT id, project_id, name, slug, deploy_type, image, dockerfile_path, build_context,
			compose_file, env_vars_encrypted, replicas, cpu_limit, memory_limit, health_check_path,
			health_check_interval, restart_policy, status, container_id, swarm_service_id,
			created_at, updated_at
		FROM services WHERE id = $1
	`
	service := &entity.Service{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&service.ID, &service.ProjectID, &service.Name, &service.Slug, &service.DeployType,
		&service.Image, &service.DockerfilePath, &service.BuildContext, &service.ComposeFile,
		&service.EnvVarsEncrypted, &service.Replicas, &service.CPULimit, &service.MemoryLimit,
		&service.HealthCheckPath, &service.HealthCheckInterval, &service.RestartPolicy,
		&service.Status, &service.ContainerID, &service.SwarmServiceID, &service.CreatedAt, &service.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (r *ServiceRepository) GetByProjectAndSlug(ctx context.Context, projectID uuid.UUID, slug string) (*entity.Service, error) {
	query := `
		SELECT id, project_id, name, slug, deploy_type, image, dockerfile_path, build_context,
			compose_file, env_vars_encrypted, replicas, cpu_limit, memory_limit, health_check_path,
			health_check_interval, restart_policy, status, container_id, swarm_service_id,
			created_at, updated_at
		FROM services WHERE project_id = $1 AND slug = $2
	`
	service := &entity.Service{}
	err := r.pool.QueryRow(ctx, query, projectID, slug).Scan(
		&service.ID, &service.ProjectID, &service.Name, &service.Slug, &service.DeployType,
		&service.Image, &service.DockerfilePath, &service.BuildContext, &service.ComposeFile,
		&service.EnvVarsEncrypted, &service.Replicas, &service.CPULimit, &service.MemoryLimit,
		&service.HealthCheckPath, &service.HealthCheckInterval, &service.RestartPolicy,
		&service.Status, &service.ContainerID, &service.SwarmServiceID, &service.CreatedAt, &service.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (r *ServiceRepository) Update(ctx context.Context, service *entity.Service) error {
	query := `
		UPDATE services SET name = $1, image = $2, dockerfile_path = $3, build_context = $4,
			compose_file = $5, env_vars_encrypted = $6, replicas = $7, cpu_limit = $8,
			memory_limit = $9, health_check_path = $10, health_check_interval = $11,
			restart_policy = $12, updated_at = $13
		WHERE id = $14
	`
	_, err := r.pool.Exec(ctx, query,
		service.Name, service.Image, service.DockerfilePath, service.BuildContext,
		service.ComposeFile, service.EnvVarsEncrypted, service.Replicas, service.CPULimit,
		service.MemoryLimit, service.HealthCheckPath, service.HealthCheckInterval,
		service.RestartPolicy, service.UpdatedAt, service.ID,
	)
	return err
}

func (r *ServiceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM services WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *ServiceRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]entity.Service, error) {
	query := `
		SELECT id, project_id, name, slug, deploy_type, image, dockerfile_path, build_context,
			compose_file, env_vars_encrypted, replicas, cpu_limit, memory_limit, health_check_path,
			health_check_interval, restart_policy, status, container_id, swarm_service_id,
			created_at, updated_at
		FROM services WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []entity.Service
	for rows.Next() {
		var s entity.Service
		err := rows.Scan(
			&s.ID, &s.ProjectID, &s.Name, &s.Slug, &s.DeployType,
			&s.Image, &s.DockerfilePath, &s.BuildContext, &s.ComposeFile,
			&s.EnvVarsEncrypted, &s.Replicas, &s.CPULimit, &s.MemoryLimit,
			&s.HealthCheckPath, &s.HealthCheckInterval, &s.RestartPolicy,
			&s.Status, &s.ContainerID, &s.SwarmServiceID, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, rows.Err()
}

func (r *ServiceRepository) ExistsByProjectAndSlug(ctx context.Context, projectID uuid.UUID, slug string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM services WHERE project_id = $1 AND slug = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, projectID, slug).Scan(&exists)
	return exists, err
}

func (r *ServiceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.ServiceStatus) error {
	query := `UPDATE services SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func (r *ServiceRepository) UpdateContainerID(ctx context.Context, id uuid.UUID, containerID *string) error {
	query := `UPDATE services SET container_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, containerID, id)
	return err
}

func (r *ServiceRepository) UpdateSwarmServiceID(ctx context.Context, id uuid.UUID, swarmServiceID *string) error {
	query := `UPDATE services SET swarm_service_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, swarmServiceID, id)
	return err
}

type DomainRepository struct {
	pool *pgxpool.Pool
}

func NewDomainRepository(pool *pgxpool.Pool) *DomainRepository {
	return &DomainRepository{pool: pool}
}

func (r *DomainRepository) Create(ctx context.Context, domain *entity.Domain) error {
	query := `
		INSERT INTO domains (id, service_id, domain, ssl_enabled, ssl_auto, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		domain.ID, domain.ServiceID, domain.Domain, domain.SSLEnabled, domain.SSLAuto, domain.CreatedAt,
	)
	return err
}

func (r *DomainRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Domain, error) {
	query := `SELECT id, service_id, domain, ssl_enabled, ssl_auto, created_at FROM domains WHERE id = $1`
	domain := &entity.Domain{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&domain.ID, &domain.ServiceID, &domain.Domain, &domain.SSLEnabled, &domain.SSLAuto, &domain.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (r *DomainRepository) GetByDomain(ctx context.Context, domainName string) (*entity.Domain, error) {
	query := `SELECT id, service_id, domain, ssl_enabled, ssl_auto, created_at FROM domains WHERE domain = $1`
	domain := &entity.Domain{}
	err := r.pool.QueryRow(ctx, query, domainName).Scan(
		&domain.ID, &domain.ServiceID, &domain.Domain, &domain.SSLEnabled, &domain.SSLAuto, &domain.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (r *DomainRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM domains WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *DomainRepository) ListByServiceID(ctx context.Context, serviceID uuid.UUID) ([]entity.Domain, error) {
	query := `SELECT id, service_id, domain, ssl_enabled, ssl_auto, created_at FROM domains WHERE service_id = $1`
	rows, err := r.pool.Query(ctx, query, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []entity.Domain
	for rows.Next() {
		var d entity.Domain
		err := rows.Scan(&d.ID, &d.ServiceID, &d.Domain, &d.SSLEnabled, &d.SSLAuto, &d.CreatedAt)
		if err != nil {
			return nil, err
		}
		domains = append(domains, d)
	}
	return domains, rows.Err()
}

func (r *DomainRepository) ExistsByDomain(ctx context.Context, domain string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM domains WHERE domain = $1)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, domain).Scan(&exists)
	return exists, err
}

type DeploymentRepository struct {
	pool *pgxpool.Pool
}

func NewDeploymentRepository(pool *pgxpool.Pool) *DeploymentRepository {
	return &DeploymentRepository{pool: pool}
}

func (r *DeploymentRepository) Create(ctx context.Context, deployment *entity.Deployment) error {
	query := `
		INSERT INTO deployments (id, service_id, triggered_by, commit_sha, commit_message, status, logs, started_at, finished_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query,
		deployment.ID, deployment.ServiceID, deployment.TriggeredBy, deployment.CommitSHA,
		deployment.CommitMessage, deployment.Status, deployment.Logs, deployment.StartedAt, deployment.FinishedAt,
	)
	return err
}

func (r *DeploymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Deployment, error) {
	query := `
		SELECT id, service_id, triggered_by, commit_sha, commit_message, status, logs, started_at, finished_at
		FROM deployments WHERE id = $1
	`
	deployment := &entity.Deployment{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&deployment.ID, &deployment.ServiceID, &deployment.TriggeredBy, &deployment.CommitSHA,
		&deployment.CommitMessage, &deployment.Status, &deployment.Logs, &deployment.StartedAt, &deployment.FinishedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func (r *DeploymentRepository) Update(ctx context.Context, deployment *entity.Deployment) error {
	query := `UPDATE deployments SET status = $1, logs = $2, finished_at = $3 WHERE id = $4`
	_, err := r.pool.Exec(ctx, query, deployment.Status, deployment.Logs, deployment.FinishedAt, deployment.ID)
	return err
}

func (r *DeploymentRepository) ListByServiceID(ctx context.Context, serviceID uuid.UUID, limit, offset int) ([]entity.Deployment, error) {
	query := `
		SELECT id, service_id, triggered_by, commit_sha, commit_message, status, logs, started_at, finished_at
		FROM deployments WHERE service_id = $1
		ORDER BY started_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, serviceID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deployments []entity.Deployment
	for rows.Next() {
		var d entity.Deployment
		err := rows.Scan(
			&d.ID, &d.ServiceID, &d.TriggeredBy, &d.CommitSHA,
			&d.CommitMessage, &d.Status, &d.Logs, &d.StartedAt, &d.FinishedAt,
		)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, d)
	}
	return deployments, rows.Err()
}

func (r *DeploymentRepository) GetLatestByServiceID(ctx context.Context, serviceID uuid.UUID) (*entity.Deployment, error) {
	query := `
		SELECT id, service_id, triggered_by, commit_sha, commit_message, status, logs, started_at, finished_at
		FROM deployments WHERE service_id = $1
		ORDER BY started_at DESC LIMIT 1
	`
	deployment := &entity.Deployment{}
	err := r.pool.QueryRow(ctx, query, serviceID).Scan(
		&deployment.ID, &deployment.ServiceID, &deployment.TriggeredBy, &deployment.CommitSHA,
		&deployment.CommitMessage, &deployment.Status, &deployment.Logs, &deployment.StartedAt, &deployment.FinishedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return deployment, nil
}
