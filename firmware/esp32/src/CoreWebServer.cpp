#include "CoreCommon.hpp"
#include "CoreWebServer.hpp"
#include "CoreConfig.hpp"
#include <WiFi.h>
#include <ArduinoJson.h>

WebServer *server = nullptr;

static String getJsonString(const JsonVariant &value) {
  return value.isNull() ? String("") : value.as<String>();
}

void handleNotFound() {
  String message = "File Not Found\n\n";
  message += "URI: ";
  message += server->uri();
  message += "\nMethod: ";
  message += (server->method() == HTTP_GET) ? "GET" : "POST";
  message += "\nArguments: ";
  message += server->args();
  message += "\n";

  for (uint8_t i = 0; i < server->args(); i++) {
    message += " " + server->argName(i) + ": " + server->arg(i) + "\n";
  }

  server->send(404, "text/plain", message);
}

void webServerSetup(int port) {
  if (server) {
    server->stop();
    delete server;
  }

  server = new WebServer(port);
  server->on("/config", HTTP_POST, handleConfigurationRequest);
  server->on("/description.xml", HTTP_GET, handleDescriptionRequest);
  server->on("/", []() {
    server->send(200, "application/json", "{\"status\":\"ok\"}");
  });
  server->onNotFound(handleNotFound);
  server->begin();
  Serial.print("HTTP server started on port ");
  Serial.println(port);
}

void shutdownWebServer(bool stopAp) {
  if (server) {
    server->stop();
    delete server;
    server = nullptr;
  }
  Serial.println("HTTP server stopped");

  if (stopAp) {
    WiFi.softAPdisconnect(true);
    Serial.println("Access point stopped");
  }
}

void loopWebServer(void) {
  if (!server) return;
  server->handleClient();
}

void handleConfigurationRequest() {
  if (server->method() != HTTP_POST) {
    server->send(405, "application/json", "{\"error\":\"Method not allowed\"}");
    return;
  }

  String body = server->arg("plain");
  if (body.length() == 0) {
    server->send(400, "application/json", "{\"error\":\"Empty request body\"}");
    return;
  }

  JsonDocument doc;
  DeserializationError error = deserializeJson(doc, body);
  if (error) {
    server->send(400, "application/json", "{\"error\":\"Invalid JSON\"}");
    return;
  }

  Config newConfig;
  newConfig.wifi.ssid = getJsonString(doc["WIFI_SSID"]);
  newConfig.wifi.pass = getJsonString(doc["WIFI_PASS"]);
  newConfig.mqtt.host = getJsonString(doc["MQTT_HOST"]);
  newConfig.mqtt.port = doc["MQTT_PORT"].isNull() ? 1883 : doc["MQTT_PORT"].as<int>();
  newConfig.mqtt.user = getJsonString(doc["MQTT_USER"]);
  newConfig.mqtt.pass = getJsonString(doc["MQTT_PASS"]);
  newConfig.mqtt.root = getJsonString(doc["MQTT_ROOT"]);
  newConfig.mdns = getJsonString(doc["MDNS"]);
  newConfig.ssdp.name = getJsonString(doc["SSDP"]);
  newConfig.ssdp.http_port = doc["HTTP_PORT"].isNull() ? 80 : doc["HTTP_PORT"].as<int>();

  config = newConfig;
  applyConfigDefaults();

  if (config.wifi.ssid.length() == 0 || config.wifi.pass.length() == 0 || config.mqtt.host.length() == 0) {
    server->send(400, "application/json", "{\"error\":\"Missing required fields\"}");
    return;
  }

  if (!saveConfig()) {
    server->send(500, "application/json", "{\"error\":\"Failed to save config\"}");
    return;
  }

  server->send(200, "application/json", "{\"status\":\"saved\"}");
  delay(200);
  ESP.restart();
}

void handleDescriptionRequest() {
  String xml = "<?xml version=\"1.0\"?>\n";
  xml += "<root xmlns=\"urn:schemas-upnp-org:device-1-0\">\n";
  xml += "  <specVersion>\n";
  xml += "    <major>1</major>\n";
  xml += "    <minor>0</minor>\n";
  xml += "  </specVersion>\n";
  xml += "  <URLBase>http://" + WiFi.localIP().toString() + ":" + String(config.ssdp.http_port) + "</URLBase>\n";
  xml += "  <device>\n";
  xml += "    <deviceType>" + config.ssdp.device_type + "</deviceType>\n";
  xml += "    <friendlyName>" + config.ssdp.name + "</friendlyName>\n";
  xml += "    <manufacturer>" + config.ssdp.manufacturer.name + "</manufacturer>\n";
  xml += "    <manufacturerURL>" + config.ssdp.manufacturer.url + "</manufacturerURL>\n";
  xml += "    <modelName>" + config.ssdp.model.name + "</modelName>\n";
  xml += "    <modelNumber>" + config.ssdp.model.number + "</modelNumber>\n";
  xml += "    <modelURL>" + config.ssdp.model.url + "</modelURL>\n";
  xml += "    <serialNumber>" + config.ssdp.serial + "</serialNumber>\n";
  xml += "    <presentationURL>" + config.ssdp.url + "</presentationURL>\n";
  xml += "  </device>\n";
  xml += "</root>\n";

  server->send(200, "text/xml", xml);
}
