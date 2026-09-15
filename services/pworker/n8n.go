package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type N8nPayload struct {
	DeviceID  string `json:"device_id"`
	PhotoURL  string `json:"photo_url"`
	Timestamp int64  `json:"timestamp"`
}

func send(deviceID, presignedURL string) error {
	payload := N8nPayload{
		DeviceID:  deviceID,
		PhotoURL:  presignedURL,
		Timestamp: time.Now().Unix(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	webHookUrl := os.Getenv("N8N_WEBHOOK_URL")
	seriliazation := "application/json"
	// fmt.Println(">>", webHookUrl)
	// fmt.Printf(">>\n%v\n", payload)
	_, err = http.Post(webHookUrl, seriliazation, bytes.NewBuffer(jsonData))
	return err
}
