# Clients and Integration Guide

## Client components in the repository

The repository includes several client-oriented assets:

- example mobile and desktop clients under the clients directory
- helper scripts and deployment assets
- REST API examples that can be used from terminal tools or custom integrations

## How the system is used from a client perspective

A client can interact with Esp Monitor in several ways:

- browse the web UI
- call the REST API directly
- upload or update device inventory data
- inspect device state and configuration history

## Integration considerations

If you plan to connect an external application to Esp Monitor, keep the following in mind:

- use a valid X-Token for protected endpoints
- treat the API as lightweight and operational rather than fully formalized
- expect device state to be distributed through MQTT, not only through REST calls
- consider Kafka if you need event streaming or audit logging

## Recommended workflow for new integrators

1. start with the server and a single device
2. verify that MQTT state updates are visible
3. try a simple REST API call to list devices
4. add automation around inventory uploads or configuration updates

## Notes

The project is currently more focused on practical device management than on polished SDKs or fully standardized client libraries.
