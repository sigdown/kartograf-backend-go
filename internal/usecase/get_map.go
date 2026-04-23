package usecase

import (
	"context"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

func (s *MapService) GetBySlug(ctx context.Context, slug string) (domain.Map, error) {
	return s.maps.GetBySlug(ctx, slug)
}
