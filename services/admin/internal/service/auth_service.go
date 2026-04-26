package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	db "github.com/hoangquan/retail-store-api/db/sqlc"
	"github.com/hoangquan/retail-store-api/pkg/auth"
	"github.com/hoangquan/retail-store-api/pkg/model/request"
	"github.com/hoangquan/retail-store-api/pkg/model/response"
	"github.com/hoangquan/retail-store-api/pkg/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
	repo       *repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthService(queries *db.Queries, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		repo:       repository.NewUserRepository(queries),
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Login(ctx context.Context, req request.LoginRequest) (*response.AuthResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		Token: token,
		User: response.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Time,
		},
	}, nil
}
