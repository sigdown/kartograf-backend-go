package usecase

import "context"

func (s *AuthService) DeleteAccount(ctx context.Context, userID int64) error {
	return s.users.Delete(ctx, userID)
}
