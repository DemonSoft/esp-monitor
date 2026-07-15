#!/usr/bin/env bash

# Function to clean up background processes on exit (Ctrl+C)
cleanup() {
    echo -e "\n🛑 Stopping all background processes..."
    kill $TUNNEL_PID $FORWARD_MOSQUITTO_PID $FORWARD_KAFKA_PID $FORWARD_GRAFANA_PID $FORWARD_SERVER_PID 2>/dev/null
    echo "👋 Environment stopped cleanly."
    exit 0
}
trap cleanup SIGINT SIGTERM
# Navigate to the directory where this script is physically located
cd "$(dirname "$0")" || exit 1

# 0. Start Docker Desktop if it is not already running
if ! docker info > /dev/null 2>&1; then
    echo "🐳 Docker Desktop is not running. Starting..."
    open -a Docker

    # Wait until Docker is fully initialized
    echo "⏳ Waiting for the Docker daemon to start..."
    until docker info > /dev/null 2>&1; do
        sleep 2
    done
    echo "✅ Docker Desktop started successfully!"
else
    echo "✅ Docker Desktop is already running."
fi


# Exit immediately if a command exits with a non-zero status
set -e

echo "===================================================="
# 
echo "Starting Kubernetes & App Stack Infrastructure Setup"
echo "===================================================="

# --- STEP 1: Dependencies Installation ---
echo "--> Checking and installing host dependencies via Homebrew..."
if ! command -v brew &> /dev/null; then
    echo "Homebrew is missing! Please install it first from https://brew.sh/"
    exit 1
fi

# Install Minikube if not present
if ! command -v minikube &> /dev/null; then
    echo "Installing Minikube..."
    brew install minikube
else
    echo "Minikube is already installed."
fi

# Install kubectl if not present
if ! command -v kubectl &> /dev/null; then
    echo "Installing kubectl..."
    brew install kubectl
else
    echo "kubectl is already installed."
fi


# --- STEP 2: Pure Cluster Initialization ---
echo "--> Cleaning up any corrupted cluster state..."
minikube delete || true

echo "--> Starting Minikube cluster tailored for Mac Intel (2017)..."
# Allocated 4 CPUs and 6GB RAM to ensure Kafka, VictoriaMetrics, and Workers run smoothly
minikube start --cpus=4 --memory=6144 --driver=docker
sleep 20

# --- STEP : Run Cluster Start Script ---
echo "--> Handing over execution to en-start-cluster for deployment..."
if [ -f "en-start-cluster" ]; then
    # Make the script executable in case permissions are not set
    chmod +x en-start-cluster
    # Run the startup script. It will deploy infra.yaml and show the status.
    exec ./en-start-cluster
else
    echo "ERROR: en-start-cluster script is missing! Cannot start the environment."
    exit 1
fi