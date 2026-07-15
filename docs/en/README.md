# Esp Monitor Documentation

This directory contains the English documentation for Esp Monitor, a platform for provisioning, monitoring, and remotely managing ESP-based devices.

## Documentation map

- [Overview](overview.md) — purpose, capabilities, architecture, and onboarding context
- [Getting Started](getting-started.md) — prerequisites, installation, first run, and basic usage
- [Firmware Guide](firmware.md) — ESP firmware behavior, provisioning flow, and MQTT actions
- [Deployment Guide](deployment.md) — local development and Kubernetes-based deployment options

## What this project provides

Esp Monitor helps installers, integrators, and service engineers work with fleets of ESP devices in a more structured way:

- initial Wi-Fi setup for new devices
- centralized device inventory and status tracking
- remote configuration updates over MQTT
- basic web and REST-based management interfaces
- optional logging and event streaming through Kafka

## Recommended reading order

1. Start with [Overview](overview.md)
2. Follow [Getting Started](getting-started.md)
3. Review [Firmware Guide](firmware.md) if you plan to work with device firmware
4. Use [Deployment Guide](deployment.md) for infrastructure setup

## Notes for new users

- The repository contains both firmware and backend services.
- The server uses an embedded SQLite database by default and can be extended to PostgreSQL later.
- MQTT and Kafka are optional infrastructure components, but they are central to remote management and telemetry.
- The current project is actively evolving; some components may still be experimental.
