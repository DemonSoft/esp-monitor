package mqttservice

import (
	"fmt"
	"remoteesp/internal/domain/utils"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string
	Port        string
	Protocol    string
	Username    string
	Password    string
	Root        string
	ClientID    string
	KeepAlive   int
	PingTimeout int
}

func (c *Config) broker() string {
	return c.Protocol + c.Host + ":" + c.Port
}

func CreateCfg() Config {
	_ = godotenv.Load() // Safety load for manifest using (without error generation)

	host := utils.GetEnv("MQTT_HOST", defaultHost)
	port := utils.GetEnv("MQTT_PORT", defaultPort)
	protocol := utils.GetEnv("MQTT_PROTOCOL", defaultProtocol)
	username := utils.GetEnv("MQTT_USERNAME", "")
	password := utils.GetEnv("MQTT_PASSWORD", "")
	root := utils.GetEnv("MQTT_ROOT", defaultTopic)
	clientId := utils.GetEnv("MQTT_CLIENT_ID", fmt.Sprintf("clientId_%d", time.Now().UTC().Unix()))
	keepAliveStr := utils.GetEnv("MQTT_KEEP_ALIVE", fmt.Sprintf("%d", defaultKeepAlive))
	pingTimeOutStr := utils.GetEnv("MQTT_PING_TIMEOUT", fmt.Sprintf("%d", defaultPingTimeOut))

	keepAlive, err := strconv.Atoi(keepAliveStr) // default 60
	if err != nil {
		keepAlive = defaultKeepAlive
	}

	timeOut, err := strconv.Atoi(pingTimeOutStr) // default 10
	if err != nil {
		timeOut = defaultPingTimeOut
	}

	return Config{
		Host:        host,
		Port:        port,
		Protocol:    protocol,
		Username:    username,
		Password:    password,
		Root:        root,
		ClientID:    clientId,
		KeepAlive:   keepAlive,
		PingTimeout: timeOut,
	}
}
