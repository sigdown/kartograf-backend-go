package usecase

import (
	"context"
	"fmt"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

func (s *MapService) DownloadURL(ctx context.Context, mapID string) (string, error) {
	archive, err := s.maps.GetActiveArchive(ctx, mapID)
	if err != nil {
		return "", err
	}

	if archive.StorageKey == "" {
		return "", fmt.Errorf("%w: archive is missing storage key", domain.ErrNotFound)
	}

	return s.storage.PresignDownload(ctx, archive.Bucket, archive.StorageKey, s.downloadTTL)
}
