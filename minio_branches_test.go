package minio_storage

import (
	"context"
	"testing"

	s3_storage "github.com/storage-lock/go-s3-storage"
	"github.com/stretchr/testify/assert"
)

func TestNewMinIOStorage_InvalidOptions(t *testing.T) {
	_, err := NewMinIOStorage(context.Background(), &s3_storage.S3StorageOptions{})
	assert.Error(t, err)
}

func TestMinIOStorage_GetName(t *testing.T) {
	assert.Equal(t, MinIOStorageName, (&MinIOStorage{}).GetName())
}
