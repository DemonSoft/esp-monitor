# ESP8266 / ESP32 Metrics Visualization and Monitoring Guide

This document provides practical recommendations, PromQL queries, and Grafana templates for visualizing and monitoring ESP8266/ESP32 GPIO metrics collected through the Vector → VictoriaMetrics pipeline.

---

## 1. Unlocking the Full Power of Grafana & PromQL

Thanks to the label-based metric model—where all GPIO values are stored under a single metric (`esp_pin_value`), while board identifiers (`device_id`) and pin names (`pin`) are represented as labels—Grafana can build dynamic, reusable dashboards instead of requiring hundreds of nearly identical charts.

### Basic PromQL Queries

**Display the value of a specific pin on a specific device:**

```promql
esp_pin_value{device_id="esp32_kitchen", pin="GPIO12"}
```

**Aggregate data (for example, the average value of an analog pin across all devices):**

```promql
avg by (pin) (esp_pin_value{pin="A0"})
```

### Dynamic Dashboards Using Variables

Instead of creating a separate panel for every device or GPIO pin, use Grafana variables. When a device is selected from a drop-down list, the available pin list will automatically adjust.

#### 1. Create a device variable (`$device`)

- **Type:** Query
- **Name:** `device`
- **Data Source:** Your VictoriaMetrics/Prometheus instance
- **Query:**

```promql
label_values(esp_pin_value, device_id)
```

- **Multi-value / Include All:** Optional.

#### 2. Create a dependent pin variable (`$pin`)

- **Type:** Query
- **Name:** `pin`
- **Data Source:** Your VictoriaMetrics/Prometheus instance
- **Query:**

```promql
label_values(esp_pin_value{device_id="$device"}, pin)
```

#### 3. Universal panel query

Use this query in graphs, gauges, or State Timeline panels:

```promql
esp_pin_value{device_id="$device", pin="$pin"}
```

### Event Analytics (Rate / Changes)

For digital GPIO pins (values `0` or `1`) connected to relays, buttons, or PIR motion sensors, it is often more useful to measure activity than to display only the current state.

**Count the number of state changes during the last hour:**

```promql
changes(esp_pin_value{device_id="$device", pin="$pin"}[1h])
```

This query is particularly useful for Heatmaps or Bar Charts that visualize sensor activity throughout the day.

---

## 2. Device Health Monitoring (Heartbeat / Watchdog)

In a streaming IoT architecture, if a microcontroller crashes, loses Wi-Fi connectivity, or loses power, it simply stops publishing messages to Kafka. Consequently, VictoriaMetrics receives no new samples, while Grafana continues displaying the last known value (or an empty panel, depending on visualization settings).

The following PromQL techniques can be used to detect offline devices.

### Method 1. Using `absent()`

The `absent()` function returns `1` when a metric has completely disappeared.

**Check whether a specific device is still reporting:**

```promql
absent(esp_pin_value{device_id="esp32_kitchen"}) == 1
```

**Limitation:** This approach does not scale well when monitoring hundreds of dynamically created devices.

---

### Method 2. Comparing Timestamps (Recommended)

A more scalable approach is to compare the current server time with the timestamp of the most recent metric received from each device.

**Detect devices that have not reported data for more than two minutes (120 seconds):**

```promql
time() - max by (device_id) (timestamp(esp_pin_value)) > 120
```

---

## Configuring Grafana Alerts

The previous query can be used directly as a Grafana Alert Rule.

1. **Expression**

```promql
time() - max by (device_id) (timestamp(esp_pin_value)) > 120
```

2. **Condition**

If the expression evaluates to **true** (returns one or more series), one or more devices have stopped sending metrics.

3. **Annotations**

**Summary**

```
Device {{ $labels.device_id }} is offline.
```

**Description**

```
The board {{ $labels.device_id }} has not reported GPIO metrics to Kafka for more than two minutes. Check its power supply and Wi-Fi connectivity.
```

4. **Notification Channels**

Configure notifications through Grafana Alertmanager using any supported integration, such as:

- Telegram
- Slack
- Email
- Discord