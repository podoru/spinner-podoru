package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/podoru/podoru/internal/domain/entity"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{pool: pool}
}

func (r *ProjectRepository) Create(ctx context.Context, project *entity.Project) error {
	query := `
		INSERT INTO projects (id, team_id, name, slug, description, github_repo, github_branch,
			github_token_encrypted, auto_deploy, webhook_secret, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.pool.Exec(ctx, query,
		project.ID, project.TeamID, project.Name, project.Slug, project.Description,
		project.GithubRepo, project.GithubBranch, project.GithubTokenEncrypted,
		project.AutoDeploy, project.WebhookSecret, project.CreatedAt, project.UpdatedAt,
	)
	return err
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Project, error) {
	query := `
		SELECT id, team_id, name, slug, description, github_repo, github_branch,
			github_token_encrypted, auto_deploy, webhook_secret, created_at, updated_at
		FROM projects WHERE id = $1
	`
	project := &entity.Project{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&project.ID, &project.TeamID, &project.Name, &project.Slug, &project.Description,
		&project.GithubRepo, &project.GithubBranch, &project.GithubTokenEncrypted,
		&project.AutoDeploy, &project.WebhookSecret, &project.CreatedAt, &project.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *ProjectRepository) GetByTeamAndSlug(ctx context.Context, teamID uuid.UUID, slug string) (*entity.Project, error) {
	query := `
		SELECT id, team_id, name, slug, description, github_repo, github_branch,
			github_token_encrypted, auto_deploy, webhook_secret, created_at, updated_at
		FROM projects WHERE team_id = $1 AND slug = $2
	`
	project := &entity.Project{}
	err := r.pool.QueryRow(ctx, query, teamID, slug).Scan(
		&project.ID, &project.TeamID, &project.Name, &project.Slug, &project.Description,
		&project.GithubRepo, &project.GithubBranch, &project.GithubTokenEncrypted,
		&project.AutoDeploy, &project.WebhookSecret, &project.CreatedAt, &project.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *entity.Project) error {
	query := `
		UPDATE projects SET name = $1, description = $2, github_repo = $3, github_branch = $4,
			github_token_encrypted = $5, auto_deploy = $6, updated_at = $7
		WHERE id = $8
	`
	_, err := r.pool.Exec(ctx, query,
		project.Name, project.Description, project.GithubRepo, project.GithubBranch,
		project.GithubTokenEncrypted, project.AutoDeploy, project.UpdatedAt, project.ID,
	)
	return err
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *ProjectRepository) ListByTeamID(ctx context.Context, teamID uuid.UUID) ([]entity.Project, error) {
	query := `
		SELECT id, team_id, name, slug, description, github_repo, github_branch,
			github_token_encrypted, auto_deploy, webhook_secret, created_at, updated_at
		FROM projects WHERE team_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []entity.Project
	for rows.Next() {
		var p entity.Project
		err := rows.Scan(
			&p.ID, &p.TeamID, &p.Name, &p.Slug, &p.Description,
			&p.GithubRepo, &p.GithubBranch, &p.GithubTokenEncrypted,
			&p.AutoDeploy, &p.WebhookSecret, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) ExistsByTeamAndSlug(ctx context.Context, teamID uuid.UUID, slug string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM projects WHERE team_id = $1 AND slug = $2)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, teamID, slug).Scan(&exists)
	return exists, err
}
