package usecase

import (
	"context"
	"fmt"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type UpdatePointInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

func (s *PointService) Update(ctx context.Context, ownerID, pointID int64, input UpdatePointInput) (domain.Point, error) {
	point, err := s.points.GetByID(ctx, pointID)
	if err != nil {
		return domain.Point{}, err
	}

	if point.OwnerID != ownerID {
		return domain.Point{}, fmt.Errorf("%w: point does not belong to user", domain.ErrForbidden)
	}

	name, err := requiredTrimmed(input.Name, "name")
	if err != nil {
		return domain.Point{}, err
	}

	if err := validateCoordinates(input.Lat, input.Lon); err != nil {
		return domain.Point{}, err
	}

	input.Name = name
	input.Description = optionalTrimmed(input.Description)
	return s.points.Update(ctx, pointID, input)
}
