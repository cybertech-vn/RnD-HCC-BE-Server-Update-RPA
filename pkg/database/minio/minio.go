package minio

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	AccessKey string `json:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	UseSSL    bool   `json:"use_ssl" yaml:"use_ssl"`
	Bucket    string `json:"bucket" yaml:"bucket"`
	Client    *minio.Client
}

func New(m *MinioClient) (*MinioClient, error) {
	var err error
	m.Client, err = m.Connect()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MinioClient) Connect() (*minio.Client, error) {

	client, err := minio.New(m.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessKey, m.SecretKey, ""),
		Secure: m.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (m *MinioClient) EnsureBucket(bucket string) error {
	ctx := context.Background()

	exists, err := m.Client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	if !exists {
		return m.Client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
	}
	return nil
}

func (m *MinioClient) UploadFile(bucket, objectName, filePath, contentType string) error {
	if bucket == "" {
		bucket = m.Bucket
	}
	ctx := context.Background()

	_, err := m.Client.FPutObject(ctx, bucket, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

func (m *MinioClient) UploadStream(bucket, objectName string, reader io.Reader, size int64, contentType string) error {
	if bucket == "" {
		bucket = m.Bucket
	}
	fmt.Println(bucket)
	ctx := context.Background()

	_, err := m.Client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

func (m *MinioClient) DownloadFile(bucket, objectName, filePath string) error {
	if bucket == "" {
		bucket = m.Bucket
	}
	ctx := context.Background()

	return m.Client.FGetObject(ctx, bucket, objectName, filePath, minio.GetObjectOptions{})
}

func (m *MinioClient) GetObject(bucket, objectName string) (*minio.Object, error) {
	if bucket == "" {
		bucket = m.Bucket
	}
	ctx := context.Background()

	obj, err := m.Client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (m *MinioClient) PresignedGet(bucket, objectName string, expiry time.Duration) (string, error) {
	if bucket == "" {
		bucket = m.Bucket
	}
	ctx := context.Background()

	url, err := m.Client.PresignedGetObject(ctx, bucket, objectName, expiry, nil)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}
