package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	host := GetEnv("KAFKA_HOST", "localhost")
	port := GetEnv("KAFKA_PORT", "9092")
	address := host + ":" + port

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(address),
		kgo.ConsumeTopics("pins"),
		kgo.ConsumerGroup("pins-worker-v1"),
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
			ssdp := string(record.Key)
			pins := string(record.Value)

			if devices[ssdp] == nil {
				device := Device{}
				devices[ssdp] = device
			}

			device := devices[ssdp]
			count := device[pins]

			device[pins] = count + 1
			devices[ssdp] = device
		}
	}

	output(devices)
}

func output(devices Devices) {
	names := make([]string, 0, len(devices))
	for k := range devices {
		names = append(names, k)
	}
	sort.Strings(names)

	for _, ssdp := range names {
		device := devices[ssdp]
		keys := make([]string, 0, len(device))
		total := 0

		maxKeyLen := 0
		for k, v := range device {
			keys = append(keys, k)
			total += v
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

		fmt.Println()
		fmt.Println()

		fmt.Printf("%*s | Count\n", maxKeyLen, ssdp)
		fmt.Println(separatorLine)

		for _, key := range keys {
			fmt.Printf("%*s | %d\n", maxKeyLen, key, device[key])
		}

		fmt.Println(separatorLine)
		fmt.Printf("%*s | %d\n", maxKeyLen, "Total", total)
		fmt.Println(separatorLine)
	}
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
