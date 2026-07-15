package kafka

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	defaultHost    = "localhost"
	defaultPort    = "9092"
	defaultTopic   = "root"
	defaultTimeout = 10 // sec
)

type Rest any

type Kafka struct {
	cfg    Config
	rest   Rest
	client *kgo.Client
}

func Create(cfg Config) *Kafka {
	kafka := &Kafka{cfg: cfg, client: nil}
	cl, err := createClient(cfg, kafka)
	log.Printf("Kafka initialization. Run %v, Host: %s, Port: %s, Topic: %s", cfg.Run, cfg.Host, cfg.Port, cfg.DefaultTopic)
	if err != nil {
		fmt.Println("Cannot create Kafka client:", err)
		return kafka
	}

	kafka.client = cl
	return kafka
}

func (k *Kafka) Close() {
	if k.client == nil {
		return
	}
	defer k.client.Close()
}

func (k *Kafka) UpdateRest(rest Rest) {
	k.rest = rest
}

func createClient(cfg Config, kafka *Kafka) (*kgo.Client, error) {
	if !cfg.Run {
		return nil, nil
	}

	if cfg.DefaultTopic == "" {
		return kgo.NewClient(
			kgo.SeedBrokers(cfg.brocker()),
			kgo.RequiredAcks(cfg.Acks),
			kgo.WithHooks(kafka),
		)
	}

	return kgo.NewClient(
		kgo.SeedBrokers(cfg.brocker()),
		kgo.DefaultProduceTopic(cfg.DefaultTopic),
		kgo.RequiredAcks(cfg.Acks),
		kgo.WithHooks(kafka),
		kgo.RecordDeliveryTimeout(time.Duration(cfg.Timeout)*time.Second),
	)
}

func (k *Kafka) OnBrokerConnect(meta kgo.BrokerMetadata, time time.Duration, conn net.Conn, err error) {
	if err != nil {
		log.Printf("[KAFKA WARN] Cannot connect to Kafka broker %s: %v", meta.Host, err)
		return
	}

	message := fmt.Sprintf("[KAFKA INFO] Connected to broker Kafka successfully: %s (connected for %v)", meta.Host, time)
	k.LogMessage(message)
}

func (k *Kafka) OnBrokerDisconnect(meta kgo.BrokerMetadata, conn net.Conn) {
	log.Printf("[KAFKA WARN] Disconnected Kafka %s !", meta.Host)
}
