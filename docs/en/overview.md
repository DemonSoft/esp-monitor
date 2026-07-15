# Product Overview

## Purpose

Esp Monitor is a client-server platform designed for provisioning, discovering, monitoring, and remotely managing ESP-based devices. It is especially useful for installers, integrators, and service engineers who deploy ESP devices into customer networks and need a simple way to configure them without physical access to each board.

## Main capabilities

The platform provides the following capabilities:

- automatic or manual device registration
- initial Wi-Fi provisioning for new boards
- remote configuration delivery through MQTT
- runtime status monitoring and telemetry reporting
- basic web-based administration and REST API access
- optional event logging through Kafka
- support for both ESP8266 and ESP32 firmware variants

## Typical use cases

- deploying a batch of ESP devices in a building or facility
- configuring Wi-Fi credentials on first boot
- updating device parameters remotely after installation
- tracking whether devices remain online and reachable
- collecting operational events for troubleshooting and auditing

## Product architecture

Esp Monitor consists of four main layers:

1. Device firmware
   - runs on ESP8266 or ESP32
   - handles local setup, Wi-Fi connection, MQTT communication, and device actions

2. Backend server
   - stores device metadata and configuration state
   - accepts updates from the web UI or REST API
   - pushes configuration changes to devices over MQTT

3. Management interfaces
   - web UI for device inspection and configuration
   - REST API for automation and integrations
   - optional command-line examples using curl

4. Infrastructure services
   - MQTT broker for device communication
   - Kafka for event logging and monitoring
   - optional Docker or Kubernetes-based deployment helpers

## Component map

- firmware/esp8266 and firmware/esp32 contain the device-side code
- server contains the Go-based management service
- clients contains example clients and helper tools
- deploy contains Docker Compose and Kubernetes deployment assets
- docs contains project documentation and deployment notes

## Operational workflow

1. A device boots and, when no saved configuration exists, starts a temporary setup mode.
2. The device exposes a provisioning endpoint or local access point for initial configuration.
3. A user provides Wi-Fi credentials and MQTT settings.
4. The device stores the configuration locally and reconnects to the network.
5. The server receives device state updates and manages the device inventory.
6. The server can push new configuration instructions and monitor device activity.

## Design principles

- keep device onboarding simple
- reduce the amount of manual intervention after deployment
- use lightweight messaging for device communication
- support both local development and a more production-like deployment model
- keep the system extensible for future integrations

## Notes for new users

The repository is a practical prototype rather than a fully polished commercial product. Some parts are intentionally simple, and some capabilities are still evolving. For most first-time users, the main goal is to get a single device provisioned and connected to the server, then expand to a small fleet once the basic flow works.
