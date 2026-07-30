package mqttservice

import (
	"encoding/json"
	"fmt"
	"remoteesp/internal/domain"
	"remoteesp/internal/model"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	defaultHost               = "localhost"
	defaultPort               = "1883"
	defaultProtocol           = "mqtt://"
	defaultTopic              = "#"
	defaultKeepAlive          = 60 // sec
	defaultPingTimeOut        = 10 // sec
	defaultRetained           = false
	defaultQOS                = 1 // Quality of Service
	defaultAutoReconnect      = true
	defaultShowingMessages    = true
	defaultRemoveAfterReceive = true
)

const (
	topicPlain  = "plain"
	topicState  = "state"
	topicAction = "action"
	topicPhoto  = "photo"
)

type Database interface {
	UpdateDevice(device model.Device) error
	UpdateDeviceAction(ssdp string, action string) error
}

type Kafka interface {
	PinsMessage(device model.Device)
	ActionMessage(ssdp string, action string)
	PhotoMessage(ssdp string, message string)
}

type MqttService struct {
	cfg    Config
	client mqtt.Client
	db     Database
	kafka  Kafka
}

func Start(cfg Config, db Database) *MqttService {
	service := &MqttService{cfg: cfg, db: db}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.broker())
	opts.SetClientID(cfg.ClientID)
	opts.SetKeepAlive(time.Duration(cfg.KeepAlive) * time.Second)
	opts.SetPingTimeout(time.Duration(cfg.PingTimeout) * time.Second)
	opts.AutoReconnect = defaultAutoReconnect

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	opts.SetOnConnectHandler(service.onConnectHandler)
	opts.SetConnectionLostHandler(service.disconnectHandler)
	opts.SetReconnectingHandler(service.reconnectionHandler)

	client := mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Printf("Error MQTT connection: %v\n", token.Error())
	}

	service.client = client
	return service
}

func (s *MqttService) Close() {
	s.client.Disconnect(250)
}

func (s *MqttService) UpdateKafka(kafka Kafka) {
	s.kafka = kafka
}

func (s *MqttService) Subscription() {
	s.subscription(s.cfg.Root + "/#")
}

func (s *MqttService) subscription(topic string) {
	qos := byte(defaultQOS)
	token := s.client.Subscribe(topic, qos, s.messageHandler)
	if token.Wait() && token.Error() != nil {
		fmt.Printf("Subscription error: %v", token.Error())
	}
	fmt.Printf("Subscription on topic %s was succeeded\n", topic)
}

func (s *MqttService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	if s.isRemoved(msg) {
		s.showRemovedMessage(msg)
		return
	}

	if defaultShowingMessages {
		s.showMessage(msg)
	}

	s.handlingProcess(msg)
}

func (s *MqttService) showMessage(msg mqtt.Message) {

	fmt.Printf("--- Received message ---\n")
	fmt.Printf("ID: %d, QoS: %d Topic: %s\n", msg.MessageID(), msg.Qos(), msg.Topic())
	data, err := s.json(msg)
	if err == nil {
		fmt.Printf("Data:    %v\n", data)
	} else {
		fmt.Printf("Data:    %s\n", string(msg.Payload()))
	}
	fmt.Printf("---------------------------\n")
}

func (s *MqttService) showRemovedMessage(msg mqtt.Message) {
	if defaultShowingMessages {
		fmt.Println("Message removed from topic:", msg.Topic())
	}
}

func (s *MqttService) sendMessage(topic string, data Message) {

	msg, err := json.Marshal(data)
	if err != nil {
		domain.Log.Log("Error:", err)
		return
	}
	//text := "Hello world"
	qos := byte(defaultQOS)
	s.client.Publish(topic, qos, defaultRetained, msg)
}

func (s *MqttService) json(msg mqtt.Message) (map[string]any, error) {

	var data map[string]any

	err := json.Unmarshal(msg.Payload(), &data)
	if err != nil {
		fmt.Printf("Parsing JSON error: %v\n", err)
		return nil, err
	}

	return data, nil
}

func (s *MqttService) remove(topic string) {
	msg := map[string]any{}
	s.sendMessage(topic, msg)
}

func (s *MqttService) isRemoved(msg mqtt.Message) bool {
	payload := msg.Payload()
	if len(payload) < 3 {
		return true
	}

	return false
}

func (s *MqttService) handlingProcess(msg mqtt.Message) {

	segment := s.lastTopicSegment(msg.Topic())

	switch segment {

	case topicPlain:
		s.receivedPlainMessage(msg)
		return

	case topicState:
		s.receivedDeviceMessage(msg)
		return

	case topicAction:
		s.receivedActionMessage(msg)
		return

	case topicPhoto:
		s.receivedPhotoMessage(msg)
		return

	}
}

func (s *MqttService) lastTopicSegment(topic string) string {
	lastSlash := strings.LastIndex(topic, "/")

	if lastSlash == -1 {
		return topic
	}

	return topic[lastSlash+1:]
}

func (s *MqttService) onConnectHandler(client mqtt.Client) {
	domain.Log.Log("Connected to MQTT broker: ", s.cfg.broker())
	s.Subscription()
}

func (s *MqttService) disconnectHandler(client mqtt.Client, error error) {
	domain.Log.Log("MQTT disconnected. Error ", error)
}

func (s *MqttService) reconnectionHandler(client mqtt.Client, opts *mqtt.ClientOptions) {
	domain.Log.Log("MQTT reconnection. ", opts.Servers)
}
