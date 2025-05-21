package db

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/your-org/lyss-chat-2.0/backend/pkg/config"
)

// MinIO 表示 MinIO 客户端
type MinIO struct {
	Client *minio.Client
	Bucket string
}

// NewMinIO 创建一个新的 MinIO 客户端
func NewMinIO(cfg config.MinIOConfig) (*MinIO, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// 检查存储桶是否存在，如果不存在则创建
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &MinIO{
		Client: client,
		Bucket: cfg.Bucket,
	}, nil
}
