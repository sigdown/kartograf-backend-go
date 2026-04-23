package usecase

import (
	"context"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type PointService struct {
	points PointRepository
}

type CreatePointInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

func NewPointService(points PointRepository) *PointService {
	return &PointService{points: points}
}

func (s *PointService) Create(ctx context.Context, ownerID int64, input CreatePointInput) (domain.Point, error) {
	name, err := requiredTrimmed(input.Name, "name")
	if err != nil {
		return domain.Point{}, err
	}

	if err := validateCoordinates(input.Lat, input.Lon); err != nil {
		return domain.Point{}, err
	}

	return s.points.Create(ctx, domain.Point{
		OwnerID:     ownerID,
		Name:        name,
		Description: optionalTrimmed(input.Description),
		Lat:         input.Lat,
		Lon:         input.Lon,
	})
}
