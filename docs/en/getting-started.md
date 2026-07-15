# Getting Started

## Prerequisites

Before using the project, make sure you have the following tools available:

- Go
- Docker and Docker Compose
- optional: PlatformIO for firmware builds
- optional: Minikube and kubectl for Kubernetes-based deployment

## Repository structure relevant to onboarding

- server — Go backend service
- firmware/esp8266 and firmware/esp32 — device firmware
- deploy — Docker Compose and Kubernetes assets
- docs — documentation in English and Russian

## 1. Prepare the environment

Create a local environment file for the server if needed:

```bash
cp server/.env.example server/.env
```

If no example file exists, create a minimal `.env` file with values for the database, REST port, MQTT, and optional Kafka settings.

## 2. Start infrastructure dependencies

For local development, start the supporting services first:

```bash
docker compose -f deploy/docker-compose.yml up -d
```

This starts Kafka and Kafdrop locally.

## 3. Run the server

From the server directory:

```bash
cd server
go run ./server
```

The server will initialize the database, connect to MQTT if configured, and start the REST endpoints.

## 4. Provision a device

For a new ESP device:

1. Flash the firmware for the target platform.
2. Power the device and let it enter setup mode when no configuration file is present.
3. Provide Wi-Fi credentials and MQTT settings through the provisioning flow.
4. Reboot the device so it joins the configured network.

## 5. Register and manage devices

After the device connects:

- verify that it appears in the server inventory
- inspect its status and recent actions
- send a new configuration if needed
- monitor logs and state reports from the MQTT channel

## 6. Useful commands

```bash
# list available devices through the REST API
curl -H "xToken: <TOKEN>" http://127.0.0.1/api/v1/devices

# upload a CSV inventory file
curl -H "xToken: <TOKEN>" -X POST http://127.0.0.1/api/v1/upload -F "body=@./tests/csv/devices.csv"
```

## Common issues

- authentication fails: verify the X-Token value used by the server and client
- MQTT messages are not received: check broker host, port, and topic root
- Kafka is not required for basic use: it can be disabled or commented out during startup
- port conflicts can happen on local machines; inspect active listeners before starting services

## Recommended next steps

- test with a single device first
- add a small CSV inventory list
- verify state updates through MQTT
- then expand to a larger fleet or a Docker/Kubernetes deployment
