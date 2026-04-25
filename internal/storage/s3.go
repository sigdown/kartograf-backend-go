package storage

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/sigdown/kartograf-backend-go/internal/usecase"
)

type S3Storage struct {
	client *minio.Client
}

func NewS3Storage(endpoint, region, accessKey, secretKey string, usePathStyle bool) (*S3Storage, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse s3 endpoint: %w", err)
	}

	client, err := minio.New(u.Host, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:       u.Scheme == "https",
		Region:       region,
		BucketLookup: bucketLookup(usePathStyle),
	})
	if err != nil {
		return nil, fmt.Errorf("create s3 client: %w", err)
	}
	return &S3Storage{client: client}, nil
}

var _ usecase.ObjectStorage = (*S3Storage)(nil)

func (s *S3Storage) EnsureBucket(ctx context.Context, bucket string) error {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("check bucket: %w", err)
	}

	if exists {
		return nil
	}

	if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}

	return nil
}

func (s *S3Storage) Delete(ctx context.Context, bucket, objectKey string) error {
	if err := s.client.RemoveObject(ctx, bucket, objectKey, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

func (s *S3Storage) PresignUpload(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
	u, err := s.client.PresignedPutObject(ctx, bucket, objectKey, expiry)
	if err != nil {
		return "", fmt.Errorf("presign upload: %w", err)
	}
	return u.String(), nil
}

func (s *S3Storage) PresignDownload(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
	u, err := s.client.PresignedGetObject(ctx, bucket, objectKey, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("presign download: %w", err)
	}
	return u.String(), nil
}

func (s *S3Storage) StatObject(ctx context.Context, bucket, objectKey string) (usecase.StoredObjectInfo, error) {
	info, err := s.client.StatObject(ctx, bucket, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return usecase.StoredObjectInfo{}, fmt.Errorf("stat object: %w", err)
	}

	return usecase.StoredObjectInfo{
		Size: info.Size,
		ETag: info.ETag,
	}, nil
}

func bucketLookup(usePathStyle bool) minio.BucketLookupType {
	if usePathStyle {
		return minio.BucketLookupPath
	}
	return minio.BucketLookupDNS
}
