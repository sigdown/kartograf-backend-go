package repository

import "github.com/sigdown/kartograf-backend-go/internal/domain"

func scanArchive(row interface{ Scan(dest ...any) error }) (domain.MapArchive, error) {
	var archive domain.MapArchive
	err := row.Scan(
		&archive.ID,
		&archive.MapID,
		&archive.Bucket,
		&archive.StorageKey,
		&archive.UploadedBy,
		&archive.SizeBytes,
		&archive.Checksum,
		&archive.Status,
		&archive.CreatedAt,
		&archive.UpdatedAt,
	)
	return archive, err
}
