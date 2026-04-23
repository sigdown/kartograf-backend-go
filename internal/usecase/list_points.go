package usecase

import (
	"context"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

func (s *PointService) List(ctx context.Context, ownerID int64) ([]domain.Point, error) {
	return s.points.ListByOwner(ctx, ownerID)
}
