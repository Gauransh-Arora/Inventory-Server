package repository

import (
	"context"
	"server/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	DB *pgxpool.Pool
}

func NewAuthRepository(conn *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{DB: conn}
}

func (r *AuthRepository) CreateUser(ctx context.Context, username, hash string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var id uuid.UUID
	query := `insert into users (username, password_hash, role, is_active) values($1, $2, 'user', true) returning id`
	err := r.DB.QueryRow(ctx, query, username, hash).Scan(&id)
	return id, err
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var u models.User
	query := `select id, username, password_hash, role, is_active, created_at from users where username=$1`
	err := r.DB.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.IsActive, &u.CreatedAt)
	return &u, err
}

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, hash string, exp time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.DB.Exec(ctx, "insert into refresh_tokens (user_id, token_hash, expires_at) values($1,$2,$3)", userID, hash, exp)
	return err
}

func (r *AuthRepository) RevokeTokens(ctx context.Context, userId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.DB.Exec(ctx, "delete from refresh_tokens where user_id = $1", userId)
	return err
}

func (r *AuthRepository) AddToDenylist(ctx context.Context, jti string, exp time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err := r.DB.Exec(ctx, "insert into jwt_denylist (jti, expires_at) values($1, $2)", jti, exp)
	return err
}

func (r *AuthRepository) IsDenyListed(ctx context.Context, jti string) bool {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var exists bool
	r.DB.QueryRow(ctx, "select exists(select 1 from jwt_denylist where jti=$1)", jti).Scan(&exists)
	return exists
}