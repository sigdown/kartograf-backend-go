package usecase

import (
	"context"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

func (s *AuthService) Me(ctx context.Context, userID int64) (domain.User, error) {
	return s.users.GetByID(ctx, userID)
}
