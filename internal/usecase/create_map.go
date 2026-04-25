package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type CreateMapUploadInput struct {
	Slug            string `json:"slug"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Year            int    `json:"year"`
	ArchiveName     string `json:"archive_name"`
	ArchiveMimeType string `json:"archive_mime_type"`
}

type CreateMapInput struct {
	MapID       string `json:"map_id"`
	ArchiveID   string `json:"archive_id"`
	StorageKey  string `json:"storage_key"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        int    `json:"year"`
}

type PresignedUploadResult struct {
	MapID            string `json:"map_id,omitempty"`
	ArchiveID        string `json:"archive_id"`
	StorageKey       string `json:"storage_key"`
	UploadURL        string `json:"upload_url"`
	ArchiveMimeType  string `json:"archive_mime_type,omitempty"`
	ExpiresInSeconds int64  `json:"expires_in_seconds"`
}

func (s *MapService) StartCreateUpload(ctx context.Context, input CreateMapUploadInput) (PresignedUploadResult, error) {
	slug, err := requiredTrimmed(input.Slug, "slug")
	if err != nil {
		return PresignedUploadResult{}, err
	}

	if _, err := requiredTrimmed(input.Title, "title"); err != nil {
		return PresignedUploadResult{}, err
	}

	if err := validateYear(input.Year); err != nil {
		return PresignedUploadResult{}, err
	}

	if _, err := validateArchiveName(input.ArchiveName); err != nil {
		return PresignedUploadResult{}, err
	}

	if _, err := s.maps.GetBySlug(ctx, slug); err == nil {
		return PresignedUploadResult{}, fmt.Errorf("%w: map with slug already exists", domain.ErrConflict)
	} else if !errors.Is(err, domain.ErrNotFound) {
		return PresignedUploadResult{}, err
	}

	mapID := newUUID()
	archiveID := newUUID()
	objectKey := buildObjectKey(slug)
	uploadURL, err := s.storage.PresignUpload(ctx, s.bucket, objectKey, s.uploadTTL)
	if err != nil {
		return PresignedUploadResult{}, err
	}

	return PresignedUploadResult{
		MapID:            mapID,
		ArchiveID:        archiveID,
		StorageKey:       objectKey,
		UploadURL:        uploadURL,
		ArchiveMimeType:  input.ArchiveMimeType,
		ExpiresInSeconds: int64(s.uploadTTL.Seconds()),
	}, nil
}

func (s *MapService) Create(ctx context.Context, actorID int64, input CreateMapInput) (domain.Map, error) {
	mapID, err := parseUUIDValue(input.MapID, "map_id")
	if err != nil {
		return domain.Map{}, err
	}

	archiveID, err := parseUUIDValue(input.ArchiveID, "archive_id")
	if err != nil {
		return domain.Map{}, err
	}

	slug, err := requiredTrimmed(input.Slug, "slug")
	if err != nil {
		return domain.Map{}, err
	}

	title, err := requiredTrimmed(input.Title, "title")
	if err != nil {
		return domain.Map{}, err
	}

	if err := validateYear(input.Year); err != nil {
		return domain.Map{}, err
	}

	if err := validateStorageKey(slug, input.StorageKey); err != nil {
		return domain.Map{}, err
	}

	objectKey := buildObjectKey(slug)
	objectInfo, err := s.storage.StatObject(ctx, s.bucket, objectKey)
	if err != nil {
		return domain.Map{}, err
	}

	created, err := s.maps.CreateWithArchive(ctx, domain.Map{
		ID:          mapID,
		CreatedBy:   actorID,
		Slug:        slug,
		Title:       title,
		Description: optionalTrimmed(input.Description),
		Year:        input.Year,
	}, domain.MapArchive{
		ID:         archiveID,
		MapID:      mapID,
		Bucket:     s.bucket,
		StorageKey: objectKey,
		UploadedBy: actorID,
		SizeBytes:  objectInfo.Size,
		Checksum:   "",
		Status:     domain.ArchiveStatusActive,
	})
	if err != nil {
		return domain.Map{}, err
	}

	return created, nil
}

func parseUUIDValue(value, field string) (string, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return "", fmt.Errorf("%w: invalid %s", domain.ErrInvalidInput, field)
	}
	return id.String(), nil
}
