package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sigdown/kartograf-backend-go/internal/auth"
	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type fakeUserRepo struct {
	createFn    func(ctx context.Context, user domain.User) (domain.User, error)
	getByIDFn   func(ctx context.Context, userID int64) (domain.User, error)
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

func (f *fakeUserRepo) GetByID(ctx context.Context, userID int64) (domain.User, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, userID)
	}
	return domain.User{}, domain.ErrNotFound
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
	return domain.Map{}, domain.ErrNotFound
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
	deleteFn     func(ctx context.Context, bucket, objectKey string) error
	presignPutFn func(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error)
	presignGetFn func(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error)
	statFn       func(ctx context.Context, bucket, objectKey string) (StoredObjectInfo, error)
	presignedKey string
	deletedKey   string
	uploadExpiry time.Duration
	downloadTTL  time.Duration
}

func (f *fakeStorage) EnsureBucket(ctx context.Context, bucket string) error {
	return nil
}

func (f *fakeStorage) Delete(ctx context.Context, bucket, objectKey string) error {
	f.deletedKey = objectKey
	if f.deleteFn != nil {
		return f.deleteFn(ctx, bucket, objectKey)
	}
	return nil
}

func (f *fakeStorage) PresignUpload(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
	f.presignedKey = objectKey
	f.uploadExpiry = expiry
	if f.presignPutFn != nil {
		return f.presignPutFn(ctx, bucket, objectKey, expiry)
	}
	return "http://example.com/upload", nil
}

func (f *fakeStorage) PresignDownload(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
	f.downloadTTL = expiry
	if f.presignGetFn != nil {
		return f.presignGetFn(ctx, bucket, objectKey, expiry)
	}
	return "http://example.com/archive", nil
}

func (f *fakeStorage) StatObject(ctx context.Context, bucket, objectKey string) (StoredObjectInfo, error) {
	if f.statFn != nil {
		return f.statFn(ctx, bucket, objectKey)
	}
	return StoredObjectInfo{Size: 1024, ETag: "etag"}, nil
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

func TestAuthServiceMeReturnsCurrentUser(t *testing.T) {
	userRepo := &fakeUserRepo{
		getByIDFn: func(ctx context.Context, userID int64) (domain.User, error) {
			return domain.User{
				ID:          userID,
				Username:    "alice",
				DisplayName: "Alice",
				Email:       "alice@example.com",
				Role:        domain.RoleUser,
			}, nil
		},
	}

	service := NewAuthService(userRepo, auth.NewTokenManager("secret", time.Hour))
	user, err := service.Me(context.Background(), 7)
	if err != nil {
		t.Fatalf("me: %v", err)
	}

	if user.ID != 7 {
		t.Fatalf("unexpected user id: %d", user.ID)
	}
	if user.Email != "alice@example.com" {
		t.Fatalf("unexpected email: %s", user.Email)
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

func TestMapServiceStartCreateUploadReturnsPresignedURL(t *testing.T) {
	storage := &fakeStorage{
		presignPutFn: func(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
			if bucket != "maps" {
				t.Fatalf("unexpected bucket: %s", bucket)
			}
			return "http://example.com/upload", nil
		},
	}

	service := NewMapService(&fakeMapRepo{}, storage, "maps", time.Minute, 2*time.Minute)
	result, err := service.StartCreateUpload(context.Background(), CreateMapUploadInput{
		Slug:            "old-map",
		Title:           "Old Map",
		Year:            1901,
		ArchiveName:     "map.pmtiles",
		ArchiveMimeType: "application/zip",
	})
	if err != nil {
		t.Fatalf("start create upload: %v", err)
	}

	if result.UploadURL != "http://example.com/upload" {
		t.Fatalf("unexpected upload url: %s", result.UploadURL)
	}
	if result.MapID == "" || result.ArchiveID == "" {
		t.Fatal("expected generated ids")
	}
	if storage.presignedKey == "" {
		t.Fatal("expected storage key")
	}
	if storage.presignedKey != "old-map.pmtiles" {
		t.Fatalf("unexpected storage key: %s", storage.presignedKey)
	}
	if storage.uploadExpiry != time.Minute {
		t.Fatalf("unexpected upload ttl: %s", storage.uploadExpiry)
	}
	if result.ExpiresInSeconds != int64(time.Minute.Seconds()) {
		t.Fatalf("unexpected expires_in_seconds: %d", result.ExpiresInSeconds)
	}
}

func TestMapServiceStartCreateUploadRejectsExistingSlug(t *testing.T) {
	repo := &fakeMapRepo{
		getBySlugFn: func(ctx context.Context, slug string) (domain.Map, error) {
			return domain.Map{ID: "map-id", Slug: slug}, nil
		},
	}

	service := NewMapService(repo, &fakeStorage{}, "maps", time.Minute, time.Minute)
	_, err := service.StartCreateUpload(context.Background(), CreateMapUploadInput{
		Slug:            "old-map",
		Title:           "Old Map",
		Year:            1901,
		ArchiveName:     "map.pmtiles",
		ArchiveMimeType: "application/zip",
	})
	if !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestMapServiceCreateDeletesUploadedObjectOnRepositoryFailure(t *testing.T) {
	repo := &fakeMapRepo{
		createFn: func(ctx context.Context, m domain.Map, archive domain.MapArchive) (domain.Map, error) {
			return domain.Map{}, domain.ErrConflict
		},
	}
	storage := &fakeStorage{}
	service := NewMapService(repo, storage, "maps", time.Minute, time.Minute)

	_, err := service.Create(context.Background(), 1, CreateMapInput{
		MapID:      "550e8400-e29b-41d4-a716-446655440000",
		ArchiveID:  "3d6f0a8b-1a2b-4c5d-9e7f-123456789abc",
		StorageKey: "old-map.pmtiles",
		Slug:       "old-map",
		Title:      "Old Map",
		Year:       1901,
	})
	if !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("expected conflict error, got %v", err)
	}
	if storage.deletedKey != "" {
		t.Fatalf("did not expect delete for stable storage key, got %s", storage.deletedKey)
	}
}

func TestMapServiceDownloadURLUsesActiveArchive(t *testing.T) {
	repo := &fakeMapRepo{
		getArchiveFn: func(ctx context.Context, mapID string) (domain.MapArchive, error) {
			return domain.MapArchive{
				ID:         "archive-id",
				MapID:      mapID,
				Bucket:     "maps",
				StorageKey: "old-map.pmtiles",
			}, nil
		},
	}
	storage := &fakeStorage{
		presignGetFn: func(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
			if bucket != "maps" {
				t.Fatalf("unexpected bucket: %s", bucket)
			}
			if objectKey != "old-map.pmtiles" {
				t.Fatalf("unexpected key: %s", objectKey)
			}
			return "http://example.com/download", nil
		},
	}
	service := NewMapService(repo, storage, "maps", time.Minute, 2*time.Minute)

	url, err := service.DownloadURL(context.Background(), "map-id")
	if err != nil {
		t.Fatalf("download url: %v", err)
	}
	if url != "http://example.com/download" {
		t.Fatalf("unexpected url: %s", url)
	}
	if storage.downloadTTL != 2*time.Minute {
		t.Fatalf("unexpected download ttl: %s", storage.downloadTTL)
	}
}

func TestMapServiceReplaceArchiveDeletesUploadedObjectOnRepositoryFailure(t *testing.T) {
	repo := &fakeMapRepo{
		getByIDFn: func(ctx context.Context, mapID string) (domain.Map, error) {
			return domain.Map{ID: mapID, Slug: "old-map"}, nil
		},
		replaceFn: func(ctx context.Context, mapID string, archive domain.MapArchive) (domain.MapArchive, error) {
			return domain.MapArchive{}, domain.ErrNotFound
		},
	}
	storage := &fakeStorage{}
	service := NewMapService(repo, storage, "maps", time.Minute, time.Minute)

	_, err := service.ReplaceArchive(context.Background(), 1, "550e8400-e29b-41d4-a716-446655440000", ReplaceMapArchiveInput{
		ArchiveID:  "3d6f0a8b-1a2b-4c5d-9e7f-123456789abc",
		StorageKey: "old-map.pmtiles",
	})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
	if storage.deletedKey != "" {
		t.Fatalf("did not expect delete for stable storage key, got %s", storage.deletedKey)
	}
}

func TestMapServiceStartReplaceArchiveUploadUsesStableSlugKey(t *testing.T) {
	repo := &fakeMapRepo{
		getByIDFn: func(ctx context.Context, mapID string) (domain.Map, error) {
			return domain.Map{ID: mapID, Slug: "old-map"}, nil
		},
	}
	storage := &fakeStorage{}
	service := NewMapService(repo, storage, "maps", 3*time.Minute, time.Minute)

	result, err := service.StartReplaceArchiveUpload(context.Background(), "550e8400-e29b-41d4-a716-446655440000", ReplaceMapArchiveUploadInput{
		ArchiveName:     "map.pmtiles",
		ArchiveMimeType: "application/pmtiles",
	})
	if err != nil {
		t.Fatalf("start replace upload: %v", err)
	}

	if result.StorageKey != "old-map.pmtiles" {
		t.Fatalf("unexpected storage key: %s", result.StorageKey)
	}
	if storage.uploadExpiry != 3*time.Minute {
		t.Fatalf("unexpected upload ttl: %s", storage.uploadExpiry)
	}
	if result.ExpiresInSeconds != int64((3 * time.Minute).Seconds()) {
		t.Fatalf("unexpected expires_in_seconds: %d", result.ExpiresInSeconds)
	}
}
