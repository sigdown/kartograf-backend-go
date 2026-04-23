package domain

import "time"

type Map struct {
	ID               string      `json:"id"`
	CreatedBy        int64       `json:"created_by"`
	CurrentArchiveID string      `json:"current_archive_id,omitempty"`
	Slug             string      `json:"slug"`
	Title            string      `json:"title"`
	Description      string      `json:"description,omitempty"`
	Year             int         `json:"year,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	CurrentArchive   *MapArchive `json:"current_archive,omitempty"`
}
