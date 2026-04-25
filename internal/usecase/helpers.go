package usecase

import (
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

func requiredTrimmed(value, field string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%w: %s is required", domain.ErrInvalidInput, field)
	}
	return trimmed, nil
}

func optionalTrimmed(value string) string {
	return strings.TrimSpace(value)
}

func validateYear(year int) error {
	if year == 0 {
		return nil
	}
	if year < 1 || year > 2100 {
		return fmt.Errorf("%w: year out of range", domain.ErrInvalidInput)
	}
	return nil
}

func validateCoordinates(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("%w: latitude out of range", domain.ErrInvalidInput)
	}
	if lon < -180 || lon > 180 {
		return fmt.Errorf("%w: longitude out of range", domain.ErrInvalidInput)
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", domain.ErrInvalidInput)
	}
	return nil
}

func newUUID() string {
	return uuid.NewString()
}

func buildObjectKey(mapID, archiveID, filename string) string {
	name := strings.ReplaceAll(strings.TrimSpace(filename), " ", "_")
	if name == "" {
		name = "archive.bin"
	}
	return fmt.Sprintf("maps/%s/%s-%s", mapID, archiveID, name)
}

func validateArchiveName(filename string) (string, error) {
	name := strings.TrimSpace(filename)
	if name == "" {
		return "", fmt.Errorf("%w: archive file name is required", domain.ErrInvalidInput)
	}
	return name, nil
}

func validateStorageKey(mapID, archiveID, storageKey string) error {
	expectedPrefix := fmt.Sprintf("maps/%s/%s-", mapID, archiveID)
	if !strings.HasPrefix(storageKey, expectedPrefix) {
		return fmt.Errorf("%w: invalid storage key", domain.ErrInvalidInput)
	}
	return nil
}
