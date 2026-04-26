package usecase

import (
	"context"
	"time"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type MapService struct {
	maps                 MapRepository
	storage              ObjectStorage
	bucket               string
	uploadTTL            time.Duration
	downloadTTL          time.Duration
	proxyEnabled         bool
	uploadBaseProxyURL   string
	downloadBaseProxyURL string
}

func NewMapService(
	maps MapRepository,
	storage ObjectStorage,
	bucket string,
	uploadTTL, downloadTTL time.Duration,
	proxyEnabled bool,
	uploadBaseProxyURL, downloadBaseProxyURL string,
) *MapService {
	return &MapService{
		maps:                 maps,
		storage:              storage,
		bucket:               bucket,
		uploadTTL:            uploadTTL,
		downloadTTL:          downloadTTL,
		proxyEnabled:         proxyEnabled,
		uploadBaseProxyURL:   uploadBaseProxyURL,
		downloadBaseProxyURL: downloadBaseProxyURL,
	}
}

func (s *MapService) List(ctx context.Context) ([]domain.Map, error) {
	return s.maps.List(ctx)
}
