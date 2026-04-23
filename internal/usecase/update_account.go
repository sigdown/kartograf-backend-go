package usecase

import (
	"context"
	"fmt"

	"github.com/sigdown/kartograf-backend-go/internal/auth"
	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type UpdateAccountInput struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func (s *AuthService) UpdateAccount(ctx context.Context, userID int64, input UpdateAccountInput) (domain.User, error) {
	input.Username = optionalTrimmed(input.Username)
	input.DisplayName = optionalTrimmed(input.DisplayName)
	input.Email = optionalTrimmed(input.Email)

	if input.Username == "" && input.DisplayName == "" && input.Email == "" && input.Password == "" {
		return domain.User{}, fmt.Errorf("%w: no fields to update", domain.ErrInvalidInput)
	}

	if input.Username != "" {
		if _, err := requiredTrimmed(input.Username, "username"); err != nil {
			return domain.User{}, err
		}
	}

	if input.Email != "" {
		if _, err := requiredTrimmed(input.Email, "email"); err != nil {
			return domain.User{}, err
		}
	}

	if input.Password != "" {
		if err := validatePassword(input.Password); err != nil {
			return domain.User{}, err
		}

		hash, err := auth.HashPassword(input.Password)
		if err != nil {
			return domain.User{}, fmt.Errorf("hash password: %w", err)
		}
		input.Password = hash
	}

	return s.users.Update(ctx, userID, input)
}
