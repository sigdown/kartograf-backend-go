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

	if _, err := s.maps.GetByID(ctx, normalizedMapID); err != nil {
		return PresignedUploadResult{}, err
	}

	archiveName, err := validateArchiveName(input.ArchiveName)
	if err != nil {
		return PresignedUploadResult{}, err
	}

	archiveID := newUUID()
	objectKey := buildObjectKey(normalizedMapID, archiveID, archiveName)
	uploadURL, err := s.storage.PresignUpload(ctx, s.bucket, objectKey, s.uploadTTL, input.ArchiveMimeType)
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

	archiveID, err := parseUUIDValue(input.ArchiveID, "archive_id")
	if err != nil {
		return domain.MapArchive{}, err
	}

	if err := validateStorageKey(normalizedMapID, archiveID, input.StorageKey); err != nil {
		return domain.MapArchive{}, err
	}

	objectInfo, err := s.storage.StatObject(ctx, s.bucket, input.StorageKey)
	if err != nil {
		return domain.MapArchive{}, err
	}

	archive, err := s.maps.ReplaceArchive(ctx, normalizedMapID, domain.MapArchive{
		ID:         archiveID,
		MapID:      normalizedMapID,
		Bucket:     s.bucket,
		StorageKey: input.StorageKey,
		UploadedBy: actorID,
		SizeBytes:  objectInfo.Size,
		Checksum:   "",
		Status:     domain.ArchiveStatusActive,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, s.bucket, input.StorageKey)
		return domain.MapArchive{}, err
	}

	return archive, nil
}
