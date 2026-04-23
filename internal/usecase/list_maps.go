package usecase

import (
	"context"
	"time"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type MapService struct {
	maps        MapRepository
	storage     ObjectStorage
	bucket      string
	downloadTTL time.Duration
}

func NewMapService(maps MapRepository, storage ObjectStorage, bucket string, downloadTTL time.Duration) *MapService {
	return &MapService{
		maps:        maps,
		storage:     storage,
		bucket:      bucket,
		downloadTTL: downloadTTL,
	}
}

func (s *MapService) List(ctx context.Context) ([]domain.Map, error) {
	return s.maps.List(ctx)
}
