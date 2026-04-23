package usecase

import (
	"bytes"
	"context"
	"fmt"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type CreateMapInput struct {
	Slug            string
	Title           string
	Description     string
	Year            int
	ArchiveName     string
	ArchiveData     []byte
	ArchiveMimeType string
}

func (s *MapService) Create(ctx context.Context, actorID int64, input CreateMapInput) (domain.Map, error) {
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

	if len(input.ArchiveData) == 0 {
		return domain.Map{}, fmt.Errorf("%w: archive file is required", domain.ErrInvalidInput)
	}

	mapID := newUUID()
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
		SizeBytes:  int64(len(input.ArchiveData)),
		Checksum:   checksumSHA256(input.ArchiveData),
		Status:     domain.ArchiveStatusActive,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, s.bucket, objectKey)
		return domain.Map{}, err
	}

	return created, nil
}
