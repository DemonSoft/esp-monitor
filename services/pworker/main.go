// This worker reads topic 'photos', and calls N8n when sees new records.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/twmb/franz-go/pkg/kgo"
)

type PhotoReport struct {
	Bucket string `json:"bucket"`
	File   string `json:"file"`
	Size   int    `json:"size"`
}

const (
	expiryDuration = 15 * time.Minute
)

func main() {
	host := GetEnv("KAFKA_HOST", "localhost")
	port := GetEnv("KAFKA_PORT", "9092")
	endpoint := os.Getenv("MINIO_ENDPOINT")    // например: "minio-service:9000"
	accessKey := os.Getenv("MINIO_ACCESS_KEY") // "minioadmin"
	secretKey := os.Getenv("MINIO_SECRET_KEY") // "minioadmin-secret-password"

	address := host + ":" + port

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(address),
		kgo.ConsumeTopics("photos"),
		kgo.ConsumerGroup("pworker-v1"),
		kgo.DisableAutoCommit(),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),

		// Uncomment this row for diagnostic:
		//kgo.WithLogger(kgo.BasicLogger(os.Stdout, kgo.LogLevelDebug, nil)),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer cl.Close()

	ctx := context.Background()

	useSSL := false // Внутри кластера часто используют HTTP без TLS

	// 2. Инициализируем клиент
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("PWorker started SUCCESSFUL.")
	for {
		fetches := cl.PollFetches(ctx)

		if fetches.NumRecords() == 0 {
			continue
		}

		deviceBatches := make(map[string][]*kgo.Record)

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			deviceID := string(record.Key)
			deviceBatches[deviceID] = append(deviceBatches[deviceID], record)
		}

		if len(deviceBatches) > 0 {
			processDeviceData(ctx, deviceBatches, minioClient)

			uncommitted := cl.UncommittedOffsets()

			cl.CommitOffsets(ctx, uncommitted, nil)
		}
	}
}

type Device map[string]int     // Key - PinsString, Value - count.
type Devices map[string]Device // Key - device name, Value - table of pins with column count.

var devices = Devices{}

func processDeviceData(ctx context.Context, batches map[string][]*kgo.Record, minIoClient *minio.Client) {

	for _, records := range batches {

		for _, record := range records {
			extractPhoto(ctx, record, minIoClient)
		}
	}

}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func extractPhoto(ctx context.Context, record *kgo.Record, minIoClient *minio.Client) {

	data := record.Value
	var pr PhotoReport
	err := json.Unmarshal(data, &pr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	path := fmt.Sprintf("%s/%s", string(record.Key), pr.File)
	object, err := minIoClient.GetObject(ctx, pr.Bucket, path, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer object.Close()

	link, err := generatePresignedURL(ctx, minIoClient, pr.Bucket, path, expiryDuration)
	if err != nil {
		fmt.Printf("Link cannot be generated: %v, bucket: %s, file: %s\n", err, pr.Bucket, path)
		return
	}

	if link != "" {
		fmt.Printf("- %s --> %d bytes\n", path, pr.Size)
	}
}

func generatePresignedURL(ctx context.Context, minIoClient *minio.Client, bucket, objectPath string, expiry time.Duration) (string, error) {
	// Указываем дополнительные параметры или заголовки при необходимости (например, Content-Disposition)
	reqParams := make(url.Values)

	// Генерируем presigned URL
	presignedURL, err := minIoClient.PresignedGetObject(ctx, bucket, objectPath, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}

	return presignedURL.String(), nil
}
