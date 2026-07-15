# Deployment Documentation

The **./../../deploy** directory contains the configurations and scripts required to deploy and run the IoT backend system. The system supports two independent deployment strategies depending on your current needs: Lightweight Local Development (Docker Compose) and Full Local Cluster Orchestration (Kubernetes via Minikube).

## Strategy 1: Lightweight Local Development (Docker Compose)
Use this strategy when you are actively writing and debugging code for esp-server or esp-worker locally in your IDE (for example, running via "go run") and only need the infrastructure dependencies (Kafka, Kafdrop) running in the background.

Files Involved:
-	docker-compose.yml: Launches external infrastructure dependencies. It spins up a standalone Kafka broker (KRaft mode) and a Kafdrop web interface.
-	restart_containers: A helper script that restarts the Docker Compose stack, wipes out old state or messages, and manually creates the required Kafka topics: root, logs, pins, actions.

How to Use:
1.	Start the infrastructure using the command: docker compose up -d
2.	Initialize or reset topics using the script: ./restart_containers
3.	Run your Go applications locally targeting localhost:9092 for Kafka.

## Strategy 2: Full Local Cluster Orchestration (Kubernetes / Minikube)
Use this strategy to simulate a production-like environment where all components, both infrastructure (Kafka, Mosquitto) and your actual backend applications (esp-server, esp-worker), run inside a unified Kubernetes network.
1.	Infrastructure Layer (Deploy First)
-	mosquitto.yaml: Deploys the Eclipse Mosquitto MQTT broker. Configures a LoadBalancer service and a dedicated NodePort 31883 allowing external physical hardware (such as ESP8266 or Wemos D1 mini) from your local Wi-Fi network to connect.
-	kafka-stack.yaml: Deploys a Kafka broker inside the cluster with automatic topic creation enabled. It also sets up Kafdrop, exposing its web UI on NodePort 31000.

2.	Application Layer (Deploy Second)
-	apps-stack.yaml: Deploys your business logic. It contains an initContainers block inside the esp-server deployment to poll Kafka and ensure topics are ready before the server starts. It deploys esp-server (image version v4) and esp-worker (image version v2). It uses imagePullPolicy Never to utilize local Docker images built on your machine. It exposes the esp-server REST API via a LoadBalancer service on NodePort 31080.
3.	Automation Script
-	start-cluster.sh: The main control script designed for macOS. It automates the entire cluster lifecycle: Ensures Docker Desktop is running, starts Minikube using the Docker driver, cleans up any zombie processes occupying port 1883, establishes a background minikube tunnel, sets up kubectl port-forward to map Mosquitto directly to your host machine, and opens the Kafdrop dashboard.
Quick Reference Cheatsheet
-	Developing code in IDE (need Kafka only): Run "docker compose up -d && ./restart_containers"
-	Testing full system in K8s: Run "./start-cluster.sh"
-	Updated backend code: Run "kubectl apply -f apps-stack.yaml"
-	Wipe and reset K8s infrastructure: Run "kubectl delete -f kafka-stack.yaml -f mosquitto.yaml"
Network and Port Mapping Summary
-	REST API (esp-server): Accessible at http://127.0.0.1:80 (via tunnel) or NodePort 31080.
-	MQTT Broker (mosquitto): Exposed on your machine at 0.0.0.0:1883 (via port-forward) or NodePort 31883.
-	Kafka Dashboard (kafdrop): Available at http://127.0.0.1:9000 (Docker Compose) or NodePort 31000 (Minikube).
