package minio_storage

// MinIO Storage 实际上是 S3 兼容存储
// 由于 go-s3-storage 已经使用 minio-go 客户端实现，
// MinIO 存储完全复用 S3 实现，仅更改名称标识

import (
	"context"

	s3_storage "github.com/storage-lock/go-s3-storage"
	"github.com/storage-lock/go-storage"
)

// MinIOStorage MinIO 对象存储实现
// MinIO 是 S3 兼容对象存储，完全复用 S3 存储实现
type MinIOStorage struct {
	*s3_storage.S3Storage
}

var _ storage.Storage = &MinIOStorage{}

// NewMinIOStorage 创建 MinIO 存储引擎
// 参数与 S3 存储完全相同，MinIO 兼容 S3 协议
func NewMinIOStorage(ctx context.Context, options *s3_storage.S3StorageOptions) (*MinIOStorage, error) {
	s3Storage, err := s3_storage.NewS3Storage(ctx, options)
	if err != nil {
		return nil, err
	}
	return &MinIOStorage{
		S3Storage: s3Storage,
	}, nil
}

const MinIOStorageName = "minio-storage"

func (x *MinIOStorage) GetName() string {
	return MinIOStorageName
}

// MinIO 与 S3 能力完全相同：
// - CAS: 通过条件 PUT（If-None-Match / If-Match ETag）实现
// - ReliableTime: 通过外部注入的 TimeProvider 满足
// - 不支持 AtomicDelete，释放走墓碑降级路径
func (x *MinIOStorage) Capabilities() []storage.StorageCapability {
	return x.S3Storage.Capabilities()
}