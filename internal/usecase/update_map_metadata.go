package usecase

import (
	"context"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type UpdateMapMetadataInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}

func (s *MapService) UpdateMetadata(ctx context.Context, mapID string, input UpdateMapMetadataInput) (domain.Map, error) {
	title, err := requiredTrimmed(input.Title, "title")
	if err != nil {
		return domain.Map{}, err
	}

	if err := validateYear(input.Year); err != nil {
		return domain.Map{}, err
	}

	input.Title = title
	input.Description = optionalTrimmed(input.Description)
	return s.maps.UpdateMetadata(ctx, mapID, input)
}
