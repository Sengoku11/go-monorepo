#!/bin/bash
set -e

# Function: Walk up the directory tree until we find a directory with "go.work"
find_repo_root() {
    local current_dir
    current_dir="$(pwd)"
    while [ "$current_dir" != "/" ]; do
        if [ -f "$current_dir/go.work" ]; then
            echo "$current_dir"
            return 0
        fi
        current_dir=$(dirname "$current_dir")
    done
    return 1
}

# Find the repository root
REPO_ROOT=$(find_repo_root)
if [ -z "$REPO_ROOT" ]; then
    echo "Error: Repository root (with go.work) not found."
    exit 1
fi

# Load .env file from the repo root and extract FEATURE_TOGGLE_PORT
ENV_FILE="$REPO_ROOT/.env"
if [ -f "$ENV_FILE" ]; then
    FEATURE_TOGGLE_PORT=$(grep '^FEATURE_TOGGLE_PORT=' "$ENV_FILE" | cut -d '=' -f2)
    if [ -z "$FEATURE_TOGGLE_PORT" ]; then
        echo "Error: FEATURE_TOGGLE_PORT is not set in $ENV_FILE"
        exit 1
    fi
else
    echo "Error: .env file not found in $REPO_ROOT"
    exit 1
fi

# Stop and remove any existing container named "flipt"
if [ "$(docker ps -aq -f name=^flipt$)" ]; then
    echo "Stopping and removing existing 'flipt' container..."
    docker stop flipt && docker rm flipt
fi

# Run docker command mapping host port from FEATURE_TOGGLE_PORT to container port 9000, with container name "flipt"
docker run -d \
    --name flipt \
    -p 8073:8080 \
    -p "${FEATURE_TOGGLE_PORT}:9000" \
    -v "$HOME/flipt:/var/opt/flipt" \
    docker.flipt.io/flipt/flipt:latest
