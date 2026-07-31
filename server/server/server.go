package main

import (
	"remoteesp/internal/database/sqlite"
	"remoteesp/internal/domain"
	"remoteesp/internal/domain/env"
	"remoteesp/internal/domain/minioStorage"
	"remoteesp/internal/network/kafka"
	"remoteesp/internal/network/mqttservice"
	"remoteesp/internal/network/rest"
)

func main() {

	// Use environtment
	env.Apply()

	// Create minIO photo storage, if you need
	minio, _ := minioStorage.NewMinIOClient()

	//Create log provider
	log := domain.Create()

	//Start Kafka
	kafkaCfg := kafka.CreateCfg()
	kafka := kafka.Create(kafkaCfg)
	log.UpdateKafka(kafka)
	kafka.HelloMessage()
	defer kafka.Close()

	// Connection to Database
	db := sqlite.Open()
	defer db.Close()

	// Start MQTT listener
	mqttCfg := mqttservice.CreateCfg()
	service := mqttservice.Start(mqttCfg, db, minio)
	service.UpdateKafka(kafka) // add Kafka to MQTT
	//service.TestSendMessage() // ONLY FOR TESTING!
	defer service.Close()

	// Start HTTP REST server
	cfg := rest.CreateCfg()
	server := rest.Create(cfg, db, service)
	kafka.UpdateRest(server) // add mqtt to Kafka.
	server.Start()
}
