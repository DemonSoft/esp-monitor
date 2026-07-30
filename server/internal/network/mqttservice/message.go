package mqttservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"remoteesp/internal/domain"
	"remoteesp/internal/model"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Message any

type PlainMessage struct {
	Message string
}

func createPlainMessage(text string) Message {
	return &PlainMessage{Message: text}
}

type DeviceMessage struct {
	SSDP    string
	MDNS    string
	Started int64
}

type ConfigMessage struct {
	SSDP   string
	Config any
}

func (s *MqttService) receivedPlainMessage(msg mqtt.Message) {
	if defaultRemoveAfterReceive {
		s.remove(msg.Topic())
	}
}

func (s *MqttService) receivedDeviceMessage(msg mqtt.Message) {
	var dm model.Device
	err := json.Unmarshal(msg.Payload(), &dm)
	if err != nil {
		domain.Log.Log("Update device was wrong:", err)
		return
	}

	if s.checkDeviceInTopic(dm.SSDP, msg) {
		err = s.db.UpdateDevice(dm)
		if err != nil {
			domain.Log.Log("DB UPDATE ERROR: ", err)
		}

		if s.kafka != nil {
			s.kafka.PinsMessage(dm)
		}
	}
	if defaultRemoveAfterReceive {
		s.remove(msg.Topic())
	}
}

func (s *MqttService) receivedActionMessage(msg mqtt.Message) {
	// DON'T REMOVE CFG MESSAGE HERE!
	// CFG message must remove only microcontroller!
	ssdp, err := s.extractDeviceFromMsg(msg)
	if err != nil {
		domain.Log.Log(err)
		return
	}
	action := string(msg.Payload())
	err = s.db.UpdateDeviceAction(ssdp, action)
	if err != nil {
		domain.Log.Log(err)
	}

	if s.kafka != nil {
		s.kafka.ActionMessage(ssdp, action)
	}

}

func (s *MqttService) receivedPhotoMessage(msg mqtt.Message) {
	ssdp, err := s.extractDeviceFromMsg(msg)
	if err != nil {
		domain.Log.Log(err)
		return
	}
	data := msg.Payload()

	str := fmt.Sprintf("%s took photo, %d bytes.", ssdp, len(data))
	fmt.Println(str)
	if s.kafka != nil {
		s.kafka.PhotoMessage(ssdp, str)
	}

	if defaultRemoveAfterReceive {
		s.remove(msg.Topic())
	}
}

// Pre-Last segment of topic MUST BE the same name as the field SSDP of the device object.
// For example,
// if ssdp == esp-001 the topic token/esp-001/state is valid,
// but token/esp-011/state is not valid
func (s *MqttService) checkDeviceInTopic(ssdp string, msg mqtt.Message) bool {

	topic := msg.Topic()
	segments := strings.Split(topic, "/")

	if len(segments)-2 < 0 {
		return false
	}

	prelast, err := s.extractDeviceFromMsg(msg)
	if err != nil || prelast != ssdp {
		return false
	}

	return true
}

func (s *MqttService) extractDeviceFromMsg(msg mqtt.Message) (string, error) {

	topic := msg.Topic()
	segments := strings.Split(topic, "/")

	if len(segments)-2 < 0 {
		return "", errors.New("Wrong lenght of topic")
	}

	return segments[len(segments)-2], nil
}

func (s *MqttService) SendDeviceConfig(ssdp string, cfg model.ConfigPayload) {

	fmt.Printf("Получена конфигурация для устройства: %s (mDNS: %s)\n", ssdp, cfg.Mdns)
	fmt.Printf("Wi-Fi SSID: %s\n", cfg.Wifi)
	fmt.Printf("SSDP name: %s (mDNS: %s)\n", cfg.Ssdp.Name, cfg.Mdns)
	fmt.Printf("MQTT Host: %s\n", cfg.Mqtt.Host)

	actionTopic := fmt.Sprintf("%s/%s/action", s.cfg.Root, ssdp)
	msg := map[string]any{}
	msg["config"] = cfg
	s.sendMessage(actionTopic, msg)
}
