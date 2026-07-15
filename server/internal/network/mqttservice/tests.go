package mqttservice

import (
	"fmt"
	"remoteesp/internal/model"
)

func (c *Config) testTopic() string {
	return c.Root + "/" + c.ClientID
}

// This method must be use only on development mode for testing
func (s *MqttService) TestSendMessage() {
	// TEST plain message
	//msg := createPlainMessage("Hello world")
	//s.sendMessage(s.cfg.testTopic()+"/plain", msg)

	// TEST state message
	// device := model.Device{SSDP: "esp-01-01", MDNS: "glob.esp", Started: time.Now().UTC().Unix()}
	// deviceTopic := fmt.Sprintf("%s/%s/state", s.cfg.Root, device.SSDP)
	// s.sendMessage(deviceTopic, device)

	// TEST action message
	// YOU CANNOT SEND ANY ACTIONS BEFORE TEST state message!
	// because ONLY STATE topic adds a new device by MQTT.
	// If it is inpossible, use .csv file BEFORE this test.
	action := model.Device{MDNS: "unknown"}
	actionTopic := fmt.Sprintf("%s/%s/action", s.cfg.Root, "esp-01-01")
	s.sendMessage(actionTopic, action)
}
