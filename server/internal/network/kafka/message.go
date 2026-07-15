package kafka

import (
	"context"
	"fmt"
	"remoteesp/internal/model"

	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	actionsTopic = "actions"
	logsTopic    = "logs"
	pinsTopic    = "pins"

	showingMessages = true
)

// This method need to sure the Kafka is connected.
func (k *Kafka) HelloMessage() {

	if k.client == nil {
		return
	}

	record := &kgo.Record{
		Key:   []byte("server"),
		Value: []byte("Server started"),
	}

	record.Topic = k.cfg.DefaultTopic

	// Async sending message
	// We give callback (Promise), who will call wnen a brocker answer.
	k.client.Produce(context.Background(), record, func(r *kgo.Record, err error) {
		if err != nil {
			fmt.Println("[Root] Delivery message error:", err)
			return
		}
		if showingMessages {
			fmt.Printf("[Root] Message saved to partition %d with offset %d\n", r.Partition, r.Offset)
		}
	})
}

func (k *Kafka) PinsMessage(device model.Device) {

	if k.client == nil {
		return
	}

	record := &kgo.Record{
		Key:   []byte(device.SSDP),
		Value: []byte(device.PinsString()),
	}

	record.Topic = pinsTopic

	k.client.Produce(context.Background(), record, func(r *kgo.Record, err error) {

		if err != nil {
			fmt.Println("[Pins] Delivery message error:", err)
			return
		}
		if showingMessages {
			fmt.Printf("[Pins] Message saved to partition %d with offset %d\n", r.Partition, r.Offset)
		}
	})
}

func (k *Kafka) ActionMessage(ssdp string, action string) {

	if k.client == nil {
		return
	}

	record := &kgo.Record{
		Key:   []byte(ssdp),
		Value: []byte(action),
	}

	record.Topic = actionsTopic

	k.client.Produce(context.Background(), record, func(r *kgo.Record, err error) {

		if err != nil {
			fmt.Println("[Action] Delivery message error:", err)
			return
		}
		if showingMessages {
			fmt.Printf("[Action] Message saved to partition %d with offset %d\n", r.Partition, r.Offset)
		}
	})
}

func (k *Kafka) LogMessage(str string) {

	if k.client == nil {
		return
	}

	record := &kgo.Record{
		Key:   []byte("event"),
		Value: []byte(str),
	}

	record.Topic = logsTopic

	k.client.Produce(context.Background(), record, func(r *kgo.Record, err error) {
		if err != nil {
			fmt.Println("[Log] Delivery message error:", err)
			return
		}
		if showingMessages {
			fmt.Printf("[Log] Message saved to partition %d with offset %d\n", r.Partition, r.Offset)
		}
	})
}
