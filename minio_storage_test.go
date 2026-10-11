package minio_storage

import (
	"context"
	"os"
	"testing"

	s3_storage "github.com/storage-lock/go-s3-storage"
	storage_test_helper "github.com/storage-lock/go-storage-test-helper"
	"github.com/stretchr/testify/assert"
)

func TestMinIOStorage(t *testing.T) {
	endpoint := os.Getenv("STORAGE_LOCK_MINIO_ENDPOINT")
	accessKey := os.Getenv("STORAGE_LOCK_MINIO_ACCESS_KEY")
	secretKey := os.Getenv("STORAGE_LOCK_MINIO_SECRET_KEY")
	bucket := os.Getenv("STORAGE_LOCK_MINIO_BUCKET")
	useSSL := os.Getenv("STORAGE_LOCK_MINIO_USE_SSL") == "true"

	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" {
		t.Skip("MinIO credentials not set, skipping test")
	}

	ctx := context.Background()
	cm := s3_storage.NewS3ConnectionManagerFromEndpoint(endpoint, accessKey, secretKey, useSSL)
	s, err := NewMinIOStorage(ctx, &s3_storage.S3StorageOptions{
		ConnectionManager: cm,
		BucketName:        bucket,
		KeyPrefix:         storage_test_helper.TestTableName,
	})
	assert.Nil(t, err)
	storage_test_helper.TestStorage(t, s)
}
