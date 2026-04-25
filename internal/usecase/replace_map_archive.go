package usecase

import (
	"context"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type ReplaceMapArchiveUploadInput struct {
	ArchiveName     string `json:"archive_name"`
	ArchiveMimeType string `json:"archive_mime_type"`
}

type ReplaceMapArchiveInput struct {
	ArchiveID  string `json:"archive_id"`
	StorageKey string `json:"storage_key"`
}

func (s *MapService) StartReplaceArchiveUpload(ctx context.Context, mapID string, input ReplaceMapArchiveUploadInput) (PresignedUploadResult, error) {
	normalizedMapID, err := parseUUIDValue(mapID, "map_id")
	if err != nil {
		return PresignedUploadResult{}, err
	}

	m, err := s.maps.GetByID(ctx, normalizedMapID)
	if err != nil {
		return PresignedUploadResult{}, err
	}

	if _, err := validateArchiveName(input.ArchiveName); err != nil {
		return PresignedUploadResult{}, err
	}

	archiveID := newUUID()
	objectKey := buildObjectKey(m.Slug)
	uploadURL, err := s.storage.PresignUpload(ctx, s.bucket, objectKey, s.uploadTTL)
	if err != nil {
		return PresignedUploadResult{}, err
	}

	return PresignedUploadResult{
		ArchiveID:        archiveID,
		StorageKey:       objectKey,
		UploadURL:        uploadURL,
		ArchiveMimeType:  input.ArchiveMimeType,
		ExpiresInSeconds: int64(s.uploadTTL.Seconds()),
	}, nil
}

func (s *MapService) ReplaceArchive(ctx context.Context, actorID int64, mapID string, input ReplaceMapArchiveInput) (domain.MapArchive, error) {
	normalizedMapID, err := parseUUIDValue(mapID, "map_id")
	if err != nil {
		return domain.MapArchive{}, err
	}

	m, err := s.maps.GetByID(ctx, normalizedMapID)
	if err != nil {
		return domain.MapArchive{}, err
	}

	archiveID, err := parseUUIDValue(input.ArchiveID, "archive_id")
	if err != nil {
		return domain.MapArchive{}, err
	}

	if err := validateStorageKey(m.Slug, input.StorageKey); err != nil {
		return domain.MapArchive{}, err
	}

	objectKey := buildObjectKey(m.Slug)
	objectInfo, err := s.storage.StatObject(ctx, s.bucket, objectKey)
	if err != nil {
		return domain.MapArchive{}, err
	}

	archive, err := s.maps.ReplaceArchive(ctx, normalizedMapID, domain.MapArchive{
		ID:         archiveID,
		MapID:      normalizedMapID,
		Bucket:     s.bucket,
		StorageKey: objectKey,
		UploadedBy: actorID,
		SizeBytes:  objectInfo.Size,
		Checksum:   "",
		Status:     domain.ArchiveStatusActive,
	})
	if err != nil {
		return domain.MapArchive{}, err
	}

	return archive, nil
}
