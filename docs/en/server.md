# Server Guide

## Role of the server

The server is the central management component of Esp Monitor. It keeps track of registered devices, receives state updates, exposes management endpoints, and sends configuration instructions to devices over MQTT.

## Main responsibilities

- maintain a list of known devices and metadata
- store configuration state for each device
- serve web pages and REST API endpoints
- receive status updates from firmware via MQTT
- push new configuration commands to devices
- optionally emit logs and events to Kafka

## Default runtime model

The Go service initializes:

1. environment configuration
2. logging and event providers
3. the database connection
4. the MQTT listener
5. the REST API server

## Data storage

By default the server uses an embedded SQLite database. The project structure also includes a PostgreSQL-oriented package layer, although it is still minimal.

## API entry points

The server exposes a simple REST API for device management and uploads. Typical actions include:

- listing devices
- filtering devices by identifier prefix
- uploading a CSV inventory file
- updating a device configuration
- deleting a device entry

The endpoints are intended for both browser-based use and automation scripts.

## Web UI

The built-in web interface allows an operator to inspect devices and manage them without writing custom client code. It is intentionally lightweight and focused on practical operational needs.

## Operational notes

- authentication is based on an X-Token value
- the server can be run locally for development or in Docker/Kubernetes for larger deployments
- Kafka integration is optional but useful for observability and event processing
- local development is often easier when MQTT and Kafka are started through the provided deployment assets
