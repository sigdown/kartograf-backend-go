package usecase

import "context"

func (s *MapService) Delete(ctx context.Context, mapID string) error {
	archive, err := s.maps.GetActiveArchive(ctx, mapID)
	if err != nil {
		return err
	}

	if archive.StorageKey != "" {
		if err := s.storage.Delete(ctx, archive.Bucket, archive.StorageKey); err != nil {
			return err
		}
	}

	return s.maps.Delete(ctx, mapID)
}
