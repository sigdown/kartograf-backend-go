package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/sigdown/kartograf-backend-go/internal/auth"
	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type LoginUserInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *AuthService) Login(ctx context.Context, input LoginUserInput) (AuthResult, error) {
	login, err := requiredTrimmed(input.Login, "login")
	if err != nil {
		return AuthResult{}, err
	}

	user, err := s.users.FindByLogin(ctx, strings.ToLower(login))
	if err != nil {
		return AuthResult{}, err
	}

	if err := auth.CheckPassword(user.PasswordHash, input.Password); err != nil {
		return AuthResult{}, fmt.Errorf("%w: bad credentials", domain.ErrUnauthorized)
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
