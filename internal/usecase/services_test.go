package usecase

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/sigdown/kartograf-backend-go/internal/auth"
	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type fakeUserRepo struct {
	createFn    func(ctx context.Context, user domain.User) (domain.User, error)
	findFn      func(ctx context.Context, login string) (domain.User, error)
	updateFn    func(ctx context.Context, userID int64, input UpdateAccountInput) (domain.User, error)
	deleteFn    func(ctx context.Context, userID int64) error
	lastCreated domain.User
}

func (f *fakeUserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	f.lastCreated = user
	if f.createFn != nil {
		return f.createFn(ctx, user)
	}
	user.ID = 1
	return user, nil
}

func (f *fakeUserRepo) FindByLogin(ctx context.Context, login string) (domain.User, error) {
	if f.findFn != nil {
		return f.findFn(ctx, login)
	}
	return domain.User{}, domain.ErrNotFound
}

func (f *fakeUserRepo) Update(ctx context.Context, userID int64, input UpdateAccountInput) (domain.User, error) {
	if f.updateFn != nil {
		return f.updateFn(ctx, userID, input)
	}
	return domain.User{}, nil
}

func (f *fakeUserRepo) Delete(ctx context.Context, userID int64) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, userID)
	}
	return nil
}

type fakePointRepo struct {
	getFn     func(ctx context.Context, pointID int64) (domain.Point, error)
	updateFn  func(ctx context.Context, pointID int64, input UpdatePointInput) (domain.Point, error)
	createFn  func(ctx context.Context, point domain.Point) (domain.Point, error)
	deleteFn  func(ctx context.Context, pointID int64) error
	listFn    func(ctx context.Context, ownerID int64) ([]domain.Point, error)
	lastPoint domain.Point
}

func (f *fakePointRepo) ListByOwner(ctx context.Context, ownerID int64) ([]domain.Point, error) {
	if f.listFn != nil {
		return f.listFn(ctx, ownerID)
	}
	return nil, nil
}

func (f *fakePointRepo) Create(ctx context.Context, point domain.Point) (domain.Point, error) {
	f.lastPoint = point
	if f.createFn != nil {
		return f.createFn(ctx, point)
	}
	point.ID = 1
	return point, nil
}

func (f *fakePointRepo) GetByID(ctx context.Context, pointID int64) (domain.Point, error) {
	if f.getFn != nil {
		return f.getFn(ctx, pointID)
	}
	return domain.Point{}, domain.ErrNotFound
}

func (f *fakePointRepo) Update(ctx context.Context, pointID int64, input UpdatePointInput) (domain.Point, error) {
	if f.updateFn != nil {
		return f.updateFn(ctx, pointID, input)
	}
	return domain.Point{ID: pointID, Name: input.Name, Description: input.Description, Lat: input.Lat, Lon: input.Lon}, nil
}

func (f *fakePointRepo) Delete(ctx context.Context, pointID int64) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, pointID)
	}
	return nil
}

type fakeMapRepo struct {
	createFn       func(ctx context.Context, m domain.Map, archive domain.MapArchive) (domain.Map, error)
	replaceFn      func(ctx context.Context, mapID string, archive domain.MapArchive) (domain.MapArchive, error)
	getArchiveFn   func(ctx context.Context, mapID string) (domain.MapArchive, error)
	listFn         func(ctx context.Context) ([]domain.Map, error)
	getBySlugFn    func(ctx context.Context, slug string) (domain.Map, error)
	getByIDFn      func(ctx context.Context, mapID string) (domain.Map, error)
	updateMetaFn   func(ctx context.Context, mapID string, input UpdateMapMetadataInput) (domain.Map, error)
	listArchivesFn func(ctx context.Context, mapID string) ([]domain.MapArchive, error)
	deleteFn       func(ctx context.Context, mapID string) error
	lastMap        domain.Map
	lastArchive    domain.MapArchive
}

func (f *fakeMapRepo) List(ctx context.Context) ([]domain.Map, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return nil, nil
}

func (f *fakeMapRepo) GetBySlug(ctx context.Context, slug string) (domain.Map, error) {
	if f.getBySlugFn != nil {
		return f.getBySlugFn(ctx, slug)
	}
	return domain.Map{}, nil
}

func (f *fakeMapRepo) GetByID(ctx context.Context, mapID string) (domain.Map, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, mapID)
	}
	return domain.Map{}, nil
}

func (f *fakeMapRepo) CreateWithArchive(ctx context.Context, m domain.Map, archive domain.MapArchive) (domain.Map, error) {
	f.lastMap = m
	f.lastArchive = archive
	if f.createFn != nil {
		return f.createFn(ctx, m, archive)
	}
	m.CurrentArchiveID = archive.ID
	return m, nil
}

func (f *fakeMapRepo) UpdateMetadata(ctx context.Context, mapID string, input UpdateMapMetadataInput) (domain.Map, error) {
	if f.updateMetaFn != nil {
		return f.updateMetaFn(ctx, mapID, input)
	}
	return domain.Map{}, nil
}

func (f *fakeMapRepo) ReplaceArchive(ctx context.Context, mapID string, archive domain.MapArchive) (domain.MapArchive, error) {
	f.lastArchive = archive
	if f.replaceFn != nil {
		return f.replaceFn(ctx, mapID, archive)
	}
	return archive, nil
}

func (f *fakeMapRepo) GetActiveArchive(ctx context.Context, mapID string) (domain.MapArchive, error) {
	if f.getArchiveFn != nil {
		return f.getArchiveFn(ctx, mapID)
	}
	return domain.MapArchive{}, nil
}

func (f *fakeMapRepo) ListArchives(ctx context.Context, mapID string) ([]domain.MapArchive, error) {
	if f.listArchivesFn != nil {
		return f.listArchivesFn(ctx, mapID)
	}
	return nil, nil
}

func (f *fakeMapRepo) Delete(ctx context.Context, mapID string) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, mapID)
	}
	return nil
}

type fakeStorage struct {
	uploadFn    func(ctx context.Context, bucket, objectKey string, body io.Reader, size int64, contentType string) error
	deleteFn    func(ctx context.Context, bucket, objectKey string) error
	presignFn   func(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error)
	uploadedKey string
	deletedKey  string
	uploadedBuf []byte
}

func (f *fakeStorage) EnsureBucket(ctx context.Context, bucket string) error {
	return nil
}

func (f *fakeStorage) Upload(ctx context.Context, bucket, objectKey string, body io.Reader, size int64, contentType string) error {
	f.uploadedKey = objectKey
	f.uploadedBuf, _ = io.ReadAll(body)
	if f.uploadFn != nil {
		return f.uploadFn(ctx, bucket, objectKey, bytes.NewReader(f.uploadedBuf), size, contentType)
	}
	return nil
}

func (f *fakeStorage) Delete(ctx context.Context, bucket, objectKey string) error {
	f.deletedKey = objectKey
	if f.deleteFn != nil {
		return f.deleteFn(ctx, bucket, objectKey)
	}
	return nil
}

func (f *fakeStorage) PresignDownload(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
	if f.presignFn != nil {
		return f.presignFn(ctx, bucket, objectKey, expiry)
	}
	return "http://example.com/archive", nil
}

func TestAuthServiceRegisterHashesPasswordAndReturnsToken(t *testing.T) {
	userRepo := &fakeUserRepo{}
	tokenManager := auth.NewTokenManager("secret", time.Hour)
	service := NewAuthService(userRepo, tokenManager)

	result, err := service.Register(context.Background(), RegisterUserInput{
		Username:    "alice",
		DisplayName: "Alice",
		Email:       "alice@example.com",
		Password:    "password123",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	if result.Token == "" {
		t.Fatal("expected token")
	}
	if userRepo.lastCreated.Role != domain.RoleUser {
		t.Fatalf("expected role %s, got %s", domain.RoleUser, userRepo.lastCreated.Role)
	}
	if userRepo.lastCreated.PasswordHash == "password123" || userRepo.lastCreated.PasswordHash == "" {
		t.Fatal("expected hashed password")
	}
}

func TestAuthServiceLoginRejectsBadPassword(t *testing.T) {
	hash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	userRepo := &fakeUserRepo{
		findFn: func(ctx context.Context, login string) (domain.User, error) {
			return domain.User{
				ID:           7,
				Username:     "alice",
				Email:        "alice@example.com",
				PasswordHash: hash,
				Role:         domain.RoleUser,
			}, nil
		},
	}

	service := NewAuthService(userRepo, auth.NewTokenManager("secret", time.Hour))
	_, err = service.Login(context.Background(), LoginUserInput{
		Login:    "alice@example.com",
		Password: "wrong-password",
	})
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}

func TestPointServiceUpdateRejectsForeignPoint(t *testing.T) {
	points := &fakePointRepo{
		getFn: func(ctx context.Context, pointID int64) (domain.Point, error) {
			return domain.Point{ID: pointID, OwnerID: 99}, nil
		},
	}
	service := NewPointService(points)

	_, err := service.Update(context.Background(), 10, 5, UpdatePointInput{
		Name: "Point",
		Lat:  10,
		Lon:  20,
	})
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
}

func TestMapServiceCreateDeletesUploadedObjectOnRepositoryFailure(t *testing.T) {
	repo := &fakeMapRepo{
		createFn: func(ctx context.Context, m domain.Map, archive domain.MapArchive) (domain.Map, error) {
			return domain.Map{}, domain.ErrConflict
		},
	}
	storage := &fakeStorage{}
	service := NewMapService(repo, storage, "maps", time.Minute)

	_, err := service.Create(context.Background(), 1, CreateMapInput{
		Slug:            "old-map",
		Title:           "Old Map",
		Year:            1901,
		ArchiveName:     "map.zip",
		ArchiveMimeType: "application/zip",
		ArchiveData:     []byte("archive-bytes"),
	})
	if !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
	if storage.uploadedKey == "" {
		t.Fatal("expected upload to happen")
	}
	if storage.deletedKey != storage.uploadedKey {
		t.Fatalf("expected cleanup delete for %s, got %s", storage.uploadedKey, storage.deletedKey)
	}
}

func TestMapServiceDownloadURLUsesActiveArchive(t *testing.T) {
	repo := &fakeMapRepo{
		getArchiveFn: func(ctx context.Context, mapID string) (domain.MapArchive, error) {
			return domain.MapArchive{
				ID:         "archive-id",
				MapID:      mapID,
				Bucket:     "maps",
				StorageKey: "maps/1/archive.zip",
			}, nil
		},
	}
	storage := &fakeStorage{
		presignFn: func(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
			if bucket != "maps" {
				t.Fatalf("unexpected bucket: %s", bucket)
			}
			if objectKey != "maps/1/archive.zip" {
				t.Fatalf("unexpected key: %s", objectKey)
			}
			return "http://example.com/download", nil
		},
	}
	service := NewMapService(repo, storage, "maps", time.Minute)

	url, err := service.DownloadURL(context.Background(), "map-id")
	if err != nil {
		t.Fatalf("download url: %v", err)
	}
	if url != "http://example.com/download" {
		t.Fatalf("unexpected url: %s", url)
	}
}
