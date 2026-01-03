package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/podoru/podoru/internal/domain/entity"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, avatar_url, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name,
		user.AvatarURL, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, name, avatar_url, role, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &entity.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name,
		&user.AvatarURL, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, name, avatar_url, role, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &entity.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name,
		&user.AvatarURL, &user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users SET name = $1, avatar_url = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.pool.Exec(ctx, query, user.Name, user.AvatarURL, user.UpdatedAt, user.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`
	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{pool: pool}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.CreatedAt,
	)
	return err
}

func (r *RefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens WHERE token_hash = $1
	`
	token := &entity.RefreshToken{}
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.pool.Exec(ctx, query, userID)
	return err
}

func (r *RefreshTokenRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = $1`
	_, err := r.pool.Exec(ctx, query, tokenHash)
	return err
}

func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.pool.Exec(ctx, query)
	return err
}
