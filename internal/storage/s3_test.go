package storage

import (
	"testing"

	"github.com/minio/minio-go/v7"
)

func TestBucketLookupUsesDNSForFirstVDS(t *testing.T) {
	if got := bucketLookup("s3.firstvds.ru", true); got != minio.BucketLookupDNS {
		t.Fatalf("expected DNS bucket lookup for FirstVDS, got %v", got)
	}
}

func TestBucketLookupUsesConfiguredStyleForOtherHosts(t *testing.T) {
	if got := bucketLookup("localhost", true); got != minio.BucketLookupPath {
		t.Fatalf("expected path-style bucket lookup for localhost, got %v", got)
	}

	if got := bucketLookup("storage.example.com", false); got != minio.BucketLookupDNS {
		t.Fatalf("expected DNS bucket lookup for non-path-style host, got %v", got)
	}
}
