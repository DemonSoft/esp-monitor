package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/twmb/franz-go/pkg/kgo"
)

// WindowSize defines how many messages from a specific device trigger an AI analysis
const WindowSize = 500
const AIModelName = "meta-llama/llama-3.1-8b-instruct"

type Device map[string]int     // Key - PinsString, Value - count.
type Devices map[string]Device // Key - device name, Value - table of pins with column count.

var devices = Devices{}

// Track the window size for each device independently
var deviceWindowCounters = make(map[string]int)

var kafkaClient *kgo.Client

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

type AIJob struct {
	DeviceID       string
	FormattedTable string
}

// Create queue for ai tasks
var aiQueue = make(chan AIJob, 100)

func main() {
	// Loading variables from .env file to process environtment
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env variables")
	}

	host := GetEnv("KAFKA_HOST", "localhost")
	port := GetEnv("KAFKA_PORT", "9092")

	address := host + ":" + port

	// Initialize OpenRouter Client using OpenAI SDK
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: OPENROUTER_API_KEY env variable is not set.")
	}

	// OpenRouter is fully compatible with OpenAI API structure
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"

	aiClient := openai.NewClientWithConfig(config)

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(address),
		kgo.ConsumeTopics("pins"),
		kgo.ConsumerGroup("ai-worker-v1"),
		kgo.DisableAutoCommit(),
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer cl.Close()

	kafkaClient = cl
	ctx := context.Background()

	// Start One thread for handling ai.
	// all request will pass step by step (Sequential)
	go startAIWorker(ctx, aiClient)

	log.Println("AI worker successfully started and listening to Kafka...")
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

func processDeviceData(batches map[string][]*kgo.Record) {
	for ssdp, records := range batches {
		for _, record := range records {
			pins := string(record.Value)

			if devices[ssdp] == nil {
				devices[ssdp] = Device{}
			}

			device := devices[ssdp]
			device[pins]++
			deviceWindowCounters[ssdp]++

			// TRIGGER: The window is full! Time to call the AI
			if deviceWindowCounters[ssdp] >= WindowSize {
				tableString := buildSingleDeviceTable(ssdp, device, deviceWindowCounters[ssdp])

				// Add tasks to queue
				aiQueue <- AIJob{
					DeviceID:       ssdp,
					FormattedTable: tableString,
				}

				// Reset aggregation for THIS specific device
				devices[ssdp] = Device{}
				deviceWindowCounters[ssdp] = 0
			}
		}
	}
}

func buildSingleDeviceTable(ssdp string, device Device, total int) string {
	keys := make([]string, 0, len(device))
	maxKeyLen := 0
	for k := range device {
		keys = append(keys, k)
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}
	sort.Strings(keys)

	if maxKeyLen < 20 {
		maxKeyLen = 20
	}

	lineLen := maxKeyLen + 13
	separatorLine := strings.Repeat("-", lineLen)

	var sb strings.Builder
	top := fmt.Sprintf("%*s | Count\n", maxKeyLen, ssdp)
	sb.WriteString(top)
	sb.WriteString(separatorLine)
	sb.WriteString("\n")

	for _, key := range keys {
		body := fmt.Sprintf("%*s | %d\n", maxKeyLen, key, device[key])
		sb.WriteString(body)
	}

	footer := fmt.Sprintf("%*s | %d\n", maxKeyLen, "Total", total)
	sb.WriteString(separatorLine)
	sb.WriteString("\n")
	sb.WriteString(footer)
	sb.WriteString(separatorLine)
	sb.WriteString("\n")

	return sb.String()
}

func analyzeWithAI(ctx context.Context, aiClient *openai.Client, deviceID string, formattedTable string) {
	// Set a reasonable timeout for the network request to OpenAI
	apiCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	log.Printf("[%s] Launching AI analysis for the current sliding window...", deviceID)

	prompt := fmt.Sprintf(`You are a specialized AI assistant for monitoring IoT devices (ESP32 and ESP8266 microcontrollers).
You are given a text table of aggregated pin states collected over a sliding time window.

Your tasks:
1. Compare the distribution of states. Identify baseline (stable) states and rare outliers.
2. Group anomalies into two categories: "noise" or "system_shift".
3. Evaluate the criticality (Status): "NORMAL", "WARNING", or "CRITICAL".

YOU MUST RETURN THE RESPONSE STRICTLY AS A VALID JSON OBJECT.
Return ONLY a valid JSON object matching the requested schema.
Do not include any explanations, markdown formatting, or code blocks like json or python.
Your response must start with '{' and end with '}'.
Follow this exact structure:

{
  "device_id": "%s",
  "status": "NORMAL | WARNING | CRITICAL",
  "summary": "Brief description of the state.",
  "anomalies": [
    {
      "type": "noise | system_shift",
      "affected_pins": ["pin_name1", "pin_name2"],
      "description": "Detailed description of what is wrong with these pins"
    }
  ]
}

EXAMPLES OF VALID RESPONSES:

Example 1 (No anomalies):
{
  "device_id": "esp8266-01",
  "status": "NORMAL",
  "summary": "All pins are operating within normal baseline parameters.",
  "anomalies": []
}

Example 2 (Anomalies found):
{
  "device_id": "esp32-main",
  "status": "WARNING",
  "summary": "System is mostly stable but minor noise detected on analog pins.",
  "anomalies": [
    {
      "type": "noise",
      "affected_pins": ["A0", "A1"],
      "description": "Frequent unstable state switches indicating electrical noise."
    }
  ]
}

Data for analysis:
%s`, deviceID, formattedTable)

	// Call OpenAI API using JSON mode to enforce structured results
	resp, err := aiClient.CreateChatCompletion(
		apiCtx,
		openai.ChatCompletionRequest{
			Model: AIModelName,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)

	if err != nil {
		log.Printf("[%s] OpenAI API request failed: %v", deviceID, err)
		return
	}

	rawContent := resp.Choices[0].Message.Content

	// Clear markdown
	rawContent = regexp.MustCompile("^?s*```[a-zA-Z]*").ReplaceAllString(rawContent, "")
	rawContent = regexp.MustCompile("```?s*$").ReplaceAllString(rawContent, "")
	rawContent = strings.TrimSpace(rawContent)

	// Parse the structured JSON response
	var report AIAnomalyReport
	err = json.Unmarshal([]byte(rawContent), &report)
	if err != nil {
		log.Printf("[%s] Failed to unmarshal AI response JSON: %v. Raw text: %s", deviceID, err, rawContent)
		return
	}

	// Route results based on critical levels
	processAIReport(report)
}

func processAIReport(report AIAnomalyReport) {
	// Filter alerts. You can change this to push to Telegram, Slack, or a local database.
	switch report.Status {
	case "CRITICAL":
		log.Printf("🚨 [CRITICAL ALERT] Device: %s -> %s", report.DeviceID, report.Summary)
		for _, anomaly := range report.Anomalies {
			log.Printf("   -> [%s] Affected: %v | %s", anomaly.Type, anomaly.AffectedPins, anomaly.Description)
		}
		anomalyMessage(report)
	case "WARNING":
		log.Printf("⚠️ [WARNING LOG] Device: %s -> %s", report.DeviceID, report.Summary)
		anomalyMessage(report)
	default:
		log.Printf("🟢 [STATUS NORMAL] Device: %s is running smoothly.", report.DeviceID)
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func startAIWorker(ctx context.Context, aiClient *openai.Client) {
	for job := range aiQueue {
		analyzeWithAI(ctx, aiClient, job.DeviceID, job.FormattedTable)
		time.Sleep(3 * time.Second)
	}
}

func anomalyMessage(report AIAnomalyReport) {
	jsonData, err := json.Marshal(report)
	if err != nil {
		return
	}

	record := &kgo.Record{
		Key:   []byte(report.DeviceID),
		Value: []byte(jsonData),
	}

	record.Topic = "anomalies"

	kafkaClient.Produce(context.Background(), record, func(r *kgo.Record, err error) {
		if err != nil {
			log.Println("Delivery message error:", err)
			return
		}
	})
}
