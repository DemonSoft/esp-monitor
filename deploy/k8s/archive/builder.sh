#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

# Navigate to the directory where this builder script is physically located
cd "$(dirname "$0")" || exit 1

echo "===================================================="
echo "🛠️ Custom Applications Builder (Minikube Context)"
echo "===================================================="

# Check if Minikube is actually running before trying to use its docker engine
if ! minikube status > /dev/null 2>&1; then
    echo "ERROR: Minikube is not running! Please start the cluster first."
    echo "Run './en-start-cluster' from the root folder, then try building again."
    exit 1
fi

echo "--> Pointing local Docker CLI inside Minikube's Docker Engine..."
eval $(minikube docker-env)

echo "--> Compiling local applications directly inside the cluster..."

# We are inside 'archive/', so we need to go 3 levels up to reach root: ../../../
PROJECTS="esp-server kworker ai-worker hworker"

for PROJ in $PROJECTS; do
    case "$PROJ" in
        "esp-server") TARGET_DIR="../../../server" ;;
        "kworker")    TARGET_DIR="../../../services/kworker" ;;
        "ai-worker")  TARGET_DIR="../../../services/ai-worker" ;;
        "hworker")    TARGET_DIR="../../../services/hworker" ;;
        *)            echo "Unknown project: $PROJ"; exit 1 ;;
    esac
    
    if [ -d "$TARGET_DIR" ]; then
        echo "----------------------------------------------------"
        echo "🚀 Navigating to $TARGET_DIR for building $PROJ..."
        cd "$TARGET_DIR"
        
        echo "📦 Building $PROJ:latest inside Minikube..."
        docker build -t "$PROJ:latest" .
        
        echo "🏷️ Tagging $PROJ with organization prefix 'demonsoft'..."
        docker tag "$PROJ:latest" "demonsoft/$PROJ:latest"
        
        # Return to archive/ directory
        cd - > /dev/null
    else
        echo "⚠️ WARNING: Project directory for $PROJ not found at $TARGET_DIR! Skipping."
    fi
done

echo "===================================================="
echo "--> Restoring host Docker context..."
# Возвращаем контекст Docker хостовой машине на случай, 
# если скрипт был запущен через 'source'
eval $(minikube docker-env -u)

echo "===================================================="
echo "✅ All local builds completed successfully inside Minikube!"
echo "===================================================="