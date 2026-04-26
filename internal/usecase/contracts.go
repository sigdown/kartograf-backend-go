package usecase

import (
	"context"
	"time"

	"github.com/sigdown/kartograf-backend-go/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByID(ctx context.Context, userID int64) (domain.User, error)
	FindByLogin(ctx context.Context, login string) (domain.User, error)
	Update(ctx context.Context, userID int64, input UpdateAccountInput) (domain.User, error)
	Delete(ctx context.Context, userID int64) error
}

type PointRepository interface {
	ListByOwner(ctx context.Context, ownerID int64) ([]domain.Point, error)
	Create(ctx context.Context, point domain.Point) (domain.Point, error)
	GetByID(ctx context.Context, pointID int64) (domain.Point, error)
	Update(ctx context.Context, pointID int64, input UpdatePointInput) (domain.Point, error)
	Delete(ctx context.Context, pointID int64) error
}

type MapRepository interface {
	List(ctx context.Context) ([]domain.Map, error)
	GetBySlug(ctx context.Context, slug string) (domain.Map, error)
	GetByID(ctx context.Context, mapID string) (domain.Map, error)
	CreateWithArchive(ctx context.Context, m domain.Map, archive domain.MapArchive) (domain.Map, error)
	UpdateMetadata(ctx context.Context, mapID string, input UpdateMapMetadataInput) (domain.Map, error)
	ReplaceArchive(ctx context.Context, mapID string, archive domain.MapArchive) (domain.MapArchive, error)
	GetActiveArchive(ctx context.Context, mapID string) (domain.MapArchive, error)
	ListArchives(ctx context.Context, mapID string) ([]domain.MapArchive, error)
	Delete(ctx context.Context, mapID string) error
}

type ObjectStorage interface {
	EnsureBucket(ctx context.Context, bucket string) error
	Delete(ctx context.Context, bucket, objectKey string) error
	PresignUpload(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error)
	PresignDownload(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error)
	StatObject(ctx context.Context, bucket, objectKey string) (StoredObjectInfo, error)
}

type StoredObjectInfo struct {
	Size int64
	ETag string
}
