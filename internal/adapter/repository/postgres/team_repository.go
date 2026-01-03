package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/podoru/spinner-podoru/internal/domain/entity"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (r *TeamRepository) Create(ctx context.Context, team *entity.Team) error {
	query := `
		INSERT INTO teams (id, name, slug, description, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(ctx, query,
		team.ID, team.Name, team.Slug, team.Description,
		team.OwnerID, team.CreatedAt, team.UpdatedAt,
	)
	return err
}

func (r *TeamRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Team, error) {
	query := `
		SELECT id, name, slug, description, owner_id, created_at, updated_at
		FROM teams WHERE id = $1
	`
	team := &entity.Team{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&team.ID, &team.Name, &team.Slug, &team.Description,
		&team.OwnerID, &team.CreatedAt, &team.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (r *TeamRepository) GetBySlug(ctx context.Context, slug string) (*entity.Team, error) {
	query := `
		SELECT id, name, slug, description, owner_id, created_at, updated_at
		FROM teams WHERE slug = $1
	`
	team := &entity.Team{}
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&team.ID, &team.Name, &team.Slug, &team.Description,
		&team.OwnerID, &team.CreatedAt, &team.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (r *TeamRepository) Update(ctx context.Context, team *entity.Team) error {
	query := `
		UPDATE teams SET name = $1, description = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.pool.Exec(ctx, query, team.Name, team.Description, team.UpdatedAt, team.ID)
	return err
}

func (r *TeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM teams WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *TeamRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.TeamWithRole, error) {
	query := `
		SELECT t.id, t.name, t.slug, t.description, t.owner_id, t.created_at, t.updated_at, tm.role
		FROM teams t
		INNER JOIN team_members tm ON t.id = tm.team_id
		WHERE tm.user_id = $1
		ORDER BY t.created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []entity.TeamWithRole
	for rows.Next() {
		var t entity.TeamWithRole
		err := rows.Scan(
			&t.ID, &t.Name, &t.Slug, &t.Description,
			&t.OwnerID, &t.CreatedAt, &t.UpdatedAt, &t.Role,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

func (r *TeamRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM teams WHERE slug = $1)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, slug).Scan(&exists)
	return exists, err
}

type TeamMemberRepository struct {
	pool *pgxpool.Pool
}

func NewTeamMemberRepository(pool *pgxpool.Pool) *TeamMemberRepository {
	return &TeamMemberRepository{pool: pool}
}

func (r *TeamMemberRepository) Create(ctx context.Context, member *entity.TeamMember) error {
	query := `
		INSERT INTO team_members (id, team_id, user_id, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		member.ID, member.TeamID, member.UserID, member.Role, member.CreatedAt,
	)
	return err
}

func (r *TeamMemberRepository) GetByTeamAndUser(ctx context.Context, teamID, userID uuid.UUID) (*entity.TeamMember, error) {
	query := `
		SELECT tm.id, tm.team_id, tm.user_id, tm.role, tm.created_at,
			   u.id, u.email, u.name, u.avatar_url, u.is_active, u.created_at, u.updated_at
		FROM team_members tm
		INNER JOIN users u ON tm.user_id = u.id
		WHERE tm.team_id = $1 AND tm.user_id = $2
	`
	member := &entity.TeamMember{User: &entity.User{}}
	err := r.pool.QueryRow(ctx, query, teamID, userID).Scan(
		&member.ID, &member.TeamID, &member.UserID, &member.Role, &member.CreatedAt,
		&member.User.ID, &member.User.Email, &member.User.Name, &member.User.AvatarURL,
		&member.User.IsActive, &member.User.CreatedAt, &member.User.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (r *TeamMemberRepository) Update(ctx context.Context, member *entity.TeamMember) error {
	query := `UPDATE team_members SET role = $1 WHERE team_id = $2 AND user_id = $3`
	_, err := r.pool.Exec(ctx, query, member.Role, member.TeamID, member.UserID)
	return err
}

func (r *TeamMemberRepository) Delete(ctx context.Context, teamID, userID uuid.UUID) error {
	query := `DELETE FROM team_members WHERE team_id = $1 AND user_id = $2`
	_, err := r.pool.Exec(ctx, query, teamID, userID)
	return err
}

func (r *TeamMemberRepository) ListByTeamID(ctx context.Context, teamID uuid.UUID) ([]entity.TeamMember, error) {
	query := `
		SELECT tm.id, tm.team_id, tm.user_id, tm.role, tm.created_at,
			   u.id, u.email, u.name, u.avatar_url, u.is_active, u.created_at, u.updated_at
		FROM team_members tm
		INNER JOIN users u ON tm.user_id = u.id
		WHERE tm.team_id = $1
		ORDER BY tm.created_at ASC
	`
	rows, err := r.pool.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []entity.TeamMember
	for rows.Next() {
		var m entity.TeamMember
		m.User = &entity.User{}
		err := rows.Scan(
			&m.ID, &m.TeamID, &m.UserID, &m.Role, &m.CreatedAt,
			&m.User.ID, &m.User.Email, &m.User.Name, &m.User.AvatarURL,
			&m.User.IsActive, &m.User.CreatedAt, &m.User.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *TeamMemberRepository) CountByTeamID(ctx context.Context, teamID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM team_members WHERE team_id = $1`
	var count int
	err := r.pool.QueryRow(ctx, query, teamID).Scan(&count)
	return count, err
}
