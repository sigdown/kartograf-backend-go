package usecase

import (
	"context"
	"fmt"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

func (s *PointService) Delete(ctx context.Context, ownerID, pointID int64) error {
	point, err := s.points.GetByID(ctx, pointID)
	if err != nil {
		return err
	}

	if point.OwnerID != ownerID {
		return fmt.Errorf("%w: point does not belong to user", domain.ErrForbidden)
	}

	return s.points.Delete(ctx, pointID)
}
