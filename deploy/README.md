## **Tiny Configuration (Only server, uses external MQTT broker):**

**You should necessary change MQTT_PROTOCOL=mqtt:// and MQTT_HOST=test.mosquitto.org in .env**  
docker compose -f deploy/docker/tiny.yaml up -d

## **Small Configuration (Server + local broker):**

**.env must be: MQTT_PROTOCOL=mqtt:// and MQTT_HOST=mosquitto**  
docker compose -f deploy/docker/small.yaml up -d

## **Medium Configuration (Server + MQTT Broker + Kafka + Kafdrop + Workers):**

docker compose -f deploy/docker/medium.yaml up -d

## **Large Configuration (Whole Stack in Docker Compose):**

docker compose -f deploy/docker/large.yaml up -d

## **Custom Configuration (Advanced Playground):**

If you want to experiment with advanced or customized configurations, you can use the standalone compose file. It fully relies on your .env file variables:  
docker compose -f docker/docker.yaml up -d

## **Huge Configuration (Production-ready local deployment in Kubernetes / Minikube):**

This configuration migrates the entire stack to Kubernetes, introduces auto-initialization of Kafka topics via a separate Job resource, handles persistent volume storage (PVC) for all data, and securely mounts your OpenRouter API key.  
Before launching, make sure you have your OpenRouter API token inside the .env file:  
OPENROUTER_API_KEY="your_actual_token_here"

To run the full K8s stack automatically, make the automation scripts executable and run one of them:

* **For English interface:**  
  chmod +x en-start-cluster  
  ./en-start-cluster

* **For Russian interface:**  
  chmod +x ru-start-cluster  
  ./ru-start-cluster

The script will automatically spin up Minikube, create the necessary secrets, apply all manifests in the correct order (pvc-stack.yaml, mosquitto.yaml, kafka-stack.yaml, monitoring-stack.yaml, apps-stack.yaml), set up port-forwarding, and launch the monitoring dashboards.

# **User Interface Available**

### **Local Deployment (Tiny, Small, Medium, Large via Docker Compose):**

* **Server (Backend API):** http://localhost  
* **Kafdrop (Kafka UI):** http://localhost:9000  
* **Grafana Dashboards:** http://localhost:3000

### **Kubernetes Deployment (Huge via Minikube):**

* **Server (Backend API):** http://localhost:8080 *(Port-forwarded from internal port 80)*  
* **Kafdrop (Kafka UI):** Available via automatic Minikube routing window  
* **Grafana Dashboards:** http://localhost:3000 *(Port-forwarded from K8s service)*  
* **Kubernetes Dashboard:** Opened automatically via script to monitor pods, PVCs, and logs


### **Grafana:**
1. **Login:** admin /admin
2. **Connections/Add new connection** Prometheus
3. **Data source** Add new data source
4. **Connection/Prometheus server URL ***http://victoria-metrics:8428 (for k8s: http://victoria-metrics-service:8428)
5. **Scroll at bottom** Save and test
6. **Dashboards** Create dashboard
7. **Dashboards** Add visualization
8. **Select data source** prometheus
- pin_values{device_id="esp32-001", pin="GPIO12"}
- pin_values{pin="A0"}
- pin_values{device_id="esp32-001"}
9. **On top right** save
