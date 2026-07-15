package kafka

import (
	"fmt"
	"remoteesp/internal/domain/utils"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Config struct {
	Host         string
	Port         string
	DefaultTopic string
	Timeout      int
	Acks         kgo.Acks
	Run          bool
}

func CreateCfg() Config {
	_ = godotenv.Load() // Safety load for manifest using (without error generation)

	host := utils.GetEnv("KAFKA_HOST", defaultHost)
	port := utils.GetEnv("KAFKA_PORT", defaultPort)
	root := utils.GetEnv("KAFKA_DEFAULT_TOPIC", defaultTopic)
	asksEnv := utils.GetEnv("KAFKA_ACKS", "")
	acks := acks(asksEnv)
	timeoutStr := utils.GetEnv("KAFKA_TIMEOUT", fmt.Sprintf("%d", defaultTimeout))
	timeout, err := strconv.Atoi(timeoutStr) // default 10
	runStr := utils.GetEnv("KAFKA_RUN", "false")
	run, err := strconv.ParseBool(runStr)
	if err != nil {
		run = false
	}

	if err != nil {
		timeout = 10 // default 10
	}

	return Config{
		Host:         host,
		Port:         port,
		DefaultTopic: root,
		Acks:         acks,
		Timeout:      timeout,
		Run:          run,
	}
}

func (c *Config) brocker() string {
	return c.Host + ":" + c.Port
}

func acks(ask string) kgo.Acks {
	switch ask {
	case "NO":
		return kgo.NoAck()
	case "LEADER":
		return kgo.LeaderAck()
	default:
		return kgo.AllISRAcks()
	}
}
