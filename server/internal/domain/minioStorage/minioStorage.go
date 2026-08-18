package minioStorage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client     *minio.Client
	BucketName string
}

func NewMinIOClient() (*MinIOClient, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")    // например: "minio-service:9000"
	accessKey := os.Getenv("MINIO_ACCESS_KEY") // "minioadmin"
	secretKey := os.Getenv("MINIO_SECRET_KEY") // "minioadmin-secret-password"
	bucketName := os.Getenv("MINIO_BUCKET")    // "esp32-photos"

	if endpoint == "" {
		endpoint = "minio-service:9000"
	}
	if bucketName == "" {
		bucketName = "esp32-photos"
	}

	// Инициализация S3 клиента (для локального MinIO useSSL = false)
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init minio client: %w", err)
	}

	ctx := context.Background()

	// Автоматически создаем бакет, если его еще нет
	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	if errBucketExists == nil && !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
		log.Printf("Bucket %s successfully created", bucketName)
	} else if errBucketExists != nil {
		return nil, fmt.Errorf("error checking bucket: %w", errBucketExists)
	}

	return &MinIOClient{
		client:     minioClient,
		BucketName: bucketName,
	}, nil
}

// UploadPhoto сохраняет массив байт (картинку) в MinIO
func (m *MinIOClient) UploadPhoto(ctx context.Context, deviceID string, filename string, photoBytes []byte) (string, error) {
	// Формируем путь внутри бакета: deviceID/filename.jpg (например: cam-01/2026-07-30_20-15-00.jpg)
	objectName := fmt.Sprintf("%s/%s", deviceID, filename)
	contentType := "image/jpeg"

	reader := bytes.NewReader(photoBytes)
	objectSize := int64(len(photoBytes))

	info, err := m.client.PutObject(ctx, m.BucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	log.Printf("Successfully uploaded %s of size %d bytes to MinIO", objectName, info.Size)
	return objectName, nil
}

func (m *MinIOClient) GetBucket() string {
	return m.BucketName
}
