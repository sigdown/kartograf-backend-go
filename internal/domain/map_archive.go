package domain

import "time"

const (
	ArchiveStatusUploaded = "UPLOADED"
	ArchiveStatusActive   = "ACTIVE"
	ArchiveStatusReplaced = "REPLACED"
	ArchiveStatusDeleted  = "DELETED"
)

type MapArchive struct {
	ID         string    `json:"id"`
	MapID      string    `json:"map_id"`
	Bucket     string    `json:"bucket"`
	StorageKey string    `json:"storage_key"`
	UploadedBy int64     `json:"uploaded_by"`
	SizeBytes  int64     `json:"size_bytes,omitempty"`
	Checksum   string    `json:"checksum,omitempty"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
