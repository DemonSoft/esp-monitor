# Firmware Guide

## Supported targets

Esp Monitor currently supports firmware for:

- ESP8266
- ESP32

The firmware code lives in the firmware/esp8266 and firmware/esp32 directories.

## Firmware responsibilities

The embedded software is responsible for:

- entering provisioning mode when no valid configuration exists
- accepting initial Wi-Fi and MQTT settings
- storing configuration locally in a JSON file
- connecting to Wi-Fi and synchronizing time
- advertising device services and exposing a local web endpoint
- publishing state updates over MQTT
- receiving action messages from the server and applying them
- resetting configuration when requested

## Provisioning flow

When the firmware starts without a saved configuration file:

1. it starts a local access point using a default setup SSID
2. it exposes a simple web service for configuration requests
3. it accepts Wi-Fi credentials, MQTT parameters, and device naming values
4. it saves the resulting configuration and reboots

After reboot, the device uses the saved configuration to connect to the network and register with the backend.

## MQTT behavior

The firmware uses MQTT for two main purposes:

- publishing device state updates
- receiving remote action messages

State updates typically include a device identifier, runtime information, and current pin values. The server can then use this information to decide whether a device needs a new configuration or other maintenance action.

## Remote actions

The firmware can process action messages sent to topics such as:

- device state topic
- action topic

If an action contains a configuration object, the device may update its saved configuration, persist it, and reboot. The exact set of actions is still being expanded depending on project needs.

## Reset behavior

If the user triggers a long button press, the device may clear its stored configuration and reboot. This is useful when the device needs to be re-provisioned from scratch.

## Operational tips

- use a unique device name and MQTT root topic for each deployment
- keep MQTT credentials private and avoid public brokers for production use
- check the serial logs when provisioning fails
- use the server logs and MQTT topics to verify that device messages arrive as expected
