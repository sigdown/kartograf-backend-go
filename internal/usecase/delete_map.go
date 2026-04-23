package usecase

import "context"

func (s *MapService) Delete(ctx context.Context, mapID string) error {
	archives, err := s.maps.ListArchives(ctx, mapID)
	if err != nil {
		return err
	}

	for _, archive := range archives {
		if archive.StorageKey == "" {
			continue
		}
		if err := s.storage.Delete(ctx, archive.Bucket, archive.StorageKey); err != nil {
			return err
		}
	}

	return s.maps.Delete(ctx, mapID)
}
