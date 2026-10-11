package minio_storage

import (
	"context"
	"os"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	s3_storage "github.com/storage-lock/go-s3-storage"
	"github.com/storage-lock/go-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
)

// TestMinIOStorageSemanticsSmoke 验证 MinIO 瘦包装继承 S3 基类不断裂。
// 需要 STORAGE_LOCK_MINIO_ENDPOINT 等环境变量（默认连本地 sl-minio）。
func TestMinIOStorageSemanticsSmoke(t *testing.T) {
	endpoint := os.Getenv("STORAGE_LOCK_MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "127.0.0.1:9000"
	}
	accessKey := os.Getenv("STORAGE_LOCK_MINIO_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}
	secretKey := os.Getenv("STORAGE_LOCK_MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}
	bucket := os.Getenv("STORAGE_LOCK_MINIO_BUCKET")
	if bucket == "" {
		bucket = "storage-lock-test"
	}
	ctx := context.Background()
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	assert.Nil(t, err)
	exists, err := client.BucketExists(ctx, bucket)
	assert.Nil(t, err)
	if !exists {
		assert.Nil(t, client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}))
	}
	cm := s3_storage.NewS3ConnectionManagerFromClient(client)
	s, err := NewMinIOStorage(ctx, &s3_storage.S3StorageOptions{
		ConnectionManager: cm,
		BucketName:        bucket,
		KeyPrefix:         "minio_smoke_test",
		TimeProvider:      &s3_storage.LocalTimeProvider{},
	})
	assert.Nil(t, err)

	// 能力继承不断裂
	caps := s.Capabilities()
	assert.Contains(t, caps, storage.CapabilityCAS)
	assert.Contains(t, caps, storage.CapabilityWatch)

	// 基础 CRUD 走通
	info := &storage.LockInformation{LockId: "smoke-lock", OwnerId: "owner-1", Version: 1}
	assert.Nil(t, s.CreateWithVersion(ctx, "smoke-lock", 1, info))
	got, err := s.Get(ctx, "smoke-lock")
	assert.Nil(t, err)
	assert.NotEmpty(t, got)
	info2 := &storage.LockInformation{LockId: "smoke-lock", OwnerId: "owner-1", Version: 2}
	assert.Nil(t, s.UpdateWithVersion(ctx, "smoke-lock", 1, 2, info2))
	_, err = s.Get(ctx, "smoke-lock")
	assert.Nil(t, err)

	// 核心锁走通
	lock, err := storage_lock.NewStorageLock(s, "smoke-lock-2")
	assert.Nil(t, err)
	assert.Nil(t, lock.Lock(ctx, "smoke-owner"))
	assert.Nil(t, lock.UnLock(ctx, "smoke-owner"))

	// 清理
	_ = s.DeleteWithVersion(ctx, "smoke-lock", 2, info2)
}
