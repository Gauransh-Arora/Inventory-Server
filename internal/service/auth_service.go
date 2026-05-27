package service

import (
	"context"
	"errors"
	"server/internal/models"
	"server/internal/repository"
	"server/internal/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{
	Repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Register(ctx context.Context, username, password string) (uuid.UUID, error){
	hash,err := utils.HashPassword(password)
	if err != nil{
		return uuid.Nil,err
	}
	return s.Repo.CreateUser(ctx, username, hash)
}

func(s *AuthService) Login (ctx context.Context, username, password string) (*models.TokenResponse, error){
	user, err := s.Repo.GetUserByUsername(ctx, username)
	if err != nil{
		return nil, errors.New("invalid credentials")
	}
	if !user.IsActive{
		return nil, errors.New("account disabled")
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash){
		return nil, errors.New("invalid credentials")
	}
	return s.GenerateTokenPair(ctx, user.ID, user.Username, user.Role)
}

func (s *AuthService) GenerateTokenPair(ctx context.Context, userID uuid.UUID, username, role string) (*models.TokenResponse, error) {
	accessToken, _, exp, err := utils.GenerateAccessToken(userID, username, role)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateRandomToken()
	if err != nil {
		return nil, err
	}
	refHash, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), 10)
	if err = s.Repo.SaveRefreshToken(ctx, userID, string(refHash), time.Now().Add(7*24*time.Hour)); err != nil {
		return nil, err
	}
	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessExpiry: exp,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, oldRefreshToken string) (*models.TokenResponse, error) {
	qCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := s.Repo.DB.Query(qCtx, "SELECT user_id, username, role, token_hash FROM refresh_tokens JOIN users ON users.id = refresh_tokens.user_id WHERE expires_at > NOW()")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userID uuid.UUID
	var username, role string
	found := false

	for rows.Next() {
		var uID uuid.UUID
		var u, r, hash string
		if err := rows.Scan(&uID, &u, &r, &hash); err == nil {
			if bcrypt.CompareHashAndPassword([]byte(hash), []byte(oldRefreshToken)) == nil {
				userID, username, role, found = uID, u, r, true
				break
			}
		}
	}

	if !found {
		return nil, errors.New("invalid refresh token")
	}

	// Token Rotation: Delete old session tokens and issue new ones
	s.Repo.RevokeTokens(ctx, userID)
	return s.GenerateTokenPair(ctx, userID, username, role)
}

func (s *AuthService) Logout(ctx context.Context, jti string, exp int64, userID uuid.UUID) error {
	s.Repo.AddToDenylist(ctx, jti, time.Unix(exp, 0))
	return s.Repo.RevokeTokens(ctx, userID)
}

