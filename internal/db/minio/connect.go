package minio

import (
	"github.com/SaraNguyen999/setup-server/internal/config"
	mini "github.com/SaraNguyen999/setup-server/pkg/database/minio"
)

var (
	MinioDB *mini.MinioClient
	err     error
)

func init() {
	MinioDB, err = mini.New(&config.CONFIG.MinioConfig.MinioClient)
	if err != nil {
		panic(err)
	}
	MinioDB.EnsureBucket(MinioDB.Bucket)
}
