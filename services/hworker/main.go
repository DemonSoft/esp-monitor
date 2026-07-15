// This worker reads topic 'anomalies', and calls webhook when sees new records.
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// AIAnomalyReport represents the expected structured JSON response from OpenAI
type AIAnomalyReport struct {
	DeviceID  string `json:"device_id"`
	Status    string `json:"status"` // NORMAL, WARNING, CRITICAL
	Summary   string `json:"summary"`
	Anomalies []struct {
		Type         string   `json:"type"` // noise, system_shift
		AffectedPins []string `json:"affected_pins"`
		Description  string   `json:"description"`
	} `json:"anomalies"`
}

var webhookURL string

func main() {
	host := GetEnv("KAFKA_HOST", "localhost")
	port := GetEnv("KAFKA_PORT", "9092")
	hook := GetEnv("WEBHOOK_URL", "")

	if len(hook) == 0 {
		fmt.Println("WEB_HOOK URL not found in environment variables.")
		return
	}

	address := host + ":" + port
	webhookURL = hook

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(address),
		kgo.ConsumeTopics("anomalies"),
		kgo.ConsumerGroup("hworker-v1"),
		kgo.DisableAutoCommit(),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),

		// Uncomment this row for diagnostic:
		//kgo.WithLogger(kgo.BasicLogger(os.Stdout, kgo.LogLevelDebug, nil)),
	)
	if err != nil {
		log.Printf("%v", err)
	}

	defer cl.Close()

	ctx := context.Background()

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
			processDeviceData(deviceBatches)

			uncommitted := cl.UncommittedOffsets()

			cl.CommitOffsets(ctx, uncommitted, nil)
		}
	}
}

type Device map[string]int     // Key - PinsString, Value - count.
type Devices map[string]Device // Key - device name, Value - table of pins with column count.

var devices = Devices{}

func processDeviceData(batches map[string][]*kgo.Record) {

	for _, records := range batches {

		for _, record := range records {
			sendToWebhook(record)
		}
	}

}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func sendToWebhook(record *kgo.Record) {

	// defender code
	if webhookURL == "" {
		return
	}

	// var report AIAnomalyReport
	// err := json.Unmarshal([]byte(record.Value), &report)
	// if err != nil {
	// 	log.Printf("[%s] Failed to unmarshal record: %v", record.Key, err)
	// 	return
	// }

	// jsonData, err := json.Marshal(report)
	// if err != nil {
	// 	log.Printf("Error marshaling report for external service: %v", err)
	// 	return
	// }

	log.Printf("%v", string(record.Value))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewBuffer(record.Value))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed to send webhook to service: %v", err)
		return
	}
	defer resp.Body.Close()
}
