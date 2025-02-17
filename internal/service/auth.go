package service

import (
	"context"
	"database/sql"
	"errors"

	"expenseapi/internal/auth"
	"expenseapi/internal/model"
	"expenseapi/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUserExists         = errors.New("usuário já existe")
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtService *auth.JWTService
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *auth.JWTService) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (s *AuthService) Register(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	// Verifica se o usuário já existe
	_, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err == nil {
		return nil, ErrUserExists
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	// Hash da senha
	passwordHash, err := model.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Cria o usuário
	userID, err := s.userRepo.Create(ctx, input.Email, passwordHash)
	if err != nil {
		return nil, err
	}

	// Busca o usuário criado
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, input model.LoginInput) (*model.LoginResponse, error) {
	// Busca o usuário pelo email
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Verifica a senha
	if err := user.ComparePassword(input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Gera os tokens
	accessToken, refreshToken, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}
