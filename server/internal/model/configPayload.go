package model

type ConfigPayload struct {
	Wifi WifiConfig `json:"wifi"`
	Ssdp SsdpConfig `json:"ssdp"`
	Mdns string     `json:"mdns"`
	Mqtt MqttConfig `json:"mqtt"`
}

type WifiConfig struct {
	Ssid string `json:"ssid"`
	Pass string `json:"pass"`
}

type SsdpConfig struct {
	Name         string             `json:"name"`
	Url          string             `json:"url"`
	Serial       string             `json:"serial"`
	Model        ModelConfig        `json:"model"`
	Manufacturer ManufacturerConfig `json:"manufacturer"`
	DeviceType   string             `json:"device_type"`
	SchemaUrl    string             `json:"schema_url"`
	HttpPort     string             `json:"http_port"`
}

type ModelConfig struct {
	Name   string `json:"name"`
	Number string `json:"number"`
	Url    string `json:"url"`
}

type ManufacturerConfig struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type MqttConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Root string `json:"root"`
}
