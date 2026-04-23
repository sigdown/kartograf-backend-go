package usecase

import (
	"context"
	"fmt"

	"github.com/sigdown/kartograf-backend-go/internal/auth"
	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type AuthTokenManager interface {
	Generate(user domain.User) (string, error)
}

type AuthService struct {
	users  UserRepository
	tokens AuthTokenManager
}

type RegisterUserInput struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type AuthResult struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

func NewAuthService(users UserRepository, tokens AuthTokenManager) *AuthService {
	return &AuthService{
		users:  users,
		tokens: tokens,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterUserInput) (AuthResult, error) {
	username, err := requiredTrimmed(input.Username, "username")
	if err != nil {
		return AuthResult{}, err
	}

	email, err := requiredTrimmed(input.Email, "email")
	if err != nil {
		return AuthResult{}, err
	}

	if err := validatePassword(input.Password); err != nil {
		return AuthResult{}, err
	}

	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		return AuthResult{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.users.Create(ctx, domain.User{
		Username:     username,
		DisplayName:  optionalTrimmed(input.DisplayName),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         domain.RoleUser,
	})
	if err != nil {
		return AuthResult{}, err
	}

	token, err := s.tokens.Generate(user)
	if err != nil {
		return AuthResult{}, fmt.Errorf("generate token: %w", err)
	}

	return AuthResult{
		Token: token,
		User:  user,
	}, nil
}
