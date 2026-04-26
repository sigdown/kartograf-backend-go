package usecase

import (
	"fmt"
	"net/url"
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

func buildObjectKey(slug string) string {
	return fmt.Sprintf("%s.pmtiles", strings.TrimSpace(slug))
}

func validateArchiveName(filename string) (string, error) {
	name := strings.TrimSpace(filename)
	if name == "" {
		return "", fmt.Errorf("%w: archive file name is required", domain.ErrInvalidInput)
	}
	if !strings.HasSuffix(strings.ToLower(name), ".pmtiles") {
		return "", fmt.Errorf("%w: archive file must be .pmtiles", domain.ErrInvalidInput)
	}
	return name, nil
}

func validateStorageKey(slug, storageKey string) error {
	expected := buildObjectKey(slug)
	if storageKey != expected {
		return fmt.Errorf("%w: invalid storage key", domain.ErrInvalidInput)
	}
	return nil
}

func RewriteToRelay(rawURL string, relayBaseURL string) (string, error) {
	direct, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse storage url: %w", err)
	}

	relayBase, err := url.Parse(strings.TrimRight(relayBaseURL, "/") + "/")
	if err != nil {
		return "", fmt.Errorf("parse relay base url: %w", err)
	}

	if relayBase.Scheme == "" || relayBase.Host == "" {
		return "", fmt.Errorf("invalid relay base url: %q", relayBaseURL)
	}

	out := *direct
	out.Scheme = relayBase.Scheme
	out.Host = relayBase.Host
	out.Path = strings.TrimRight(relayBase.Path, "/") + direct.EscapedPath()

	return out.String(), nil
}
