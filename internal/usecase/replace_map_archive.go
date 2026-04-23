package usecase

import (
	"bytes"
	"context"
	"fmt"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type ReplaceMapArchiveInput struct {
	ArchiveName     string
	ArchiveData     []byte
	ArchiveMimeType string
}

func (s *MapService) ReplaceArchive(ctx context.Context, actorID int64, mapID string, input ReplaceMapArchiveInput) (domain.MapArchive, error) {
	if len(input.ArchiveData) == 0 {
		return domain.MapArchive{}, fmt.Errorf("%w: archive file is required", domain.ErrInvalidInput)
	}

	archiveID := newUUID()
	objectKey := buildObjectKey(mapID, archiveID, input.ArchiveName)

	if err := s.storage.Upload(
		ctx,
		s.bucket,
		objectKey,
		bytes.NewReader(input.ArchiveData),
		int64(len(input.ArchiveData)),
		input.ArchiveMimeType,
	); err != nil {
		return domain.MapArchive{}, err
	}

	archive, err := s.maps.ReplaceArchive(ctx, mapID, domain.MapArchive{
		ID:         archiveID,
		MapID:      mapID,
		Bucket:     s.bucket,
		StorageKey: objectKey,
		UploadedBy: actorID,
		SizeBytes:  int64(len(input.ArchiveData)),
		Checksum:   checksumSHA256(input.ArchiveData),
		Status:     domain.ArchiveStatusActive,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, s.bucket, objectKey)
		return domain.MapArchive{}, err
	}

	return archive, nil
}
