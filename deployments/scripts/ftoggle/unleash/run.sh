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

# Load .env file from the repo root
ENV_FILE="$REPO_ROOT/.env"
if [ ! -f "$ENV_FILE" ]; then
    echo "Error: .env file not found in $REPO_ROOT"
    exit 1
fi

# Define required variables based on your mapping
required_vars=(
  FEATURE_TOGGLE_INIT_FRONTEND_API_TOKENS
  FEATURE_TOGGLE_INIT_CLIENT_API_TOKENS
  FEATURE_TOGGLE_DATABASE_USERNAME
  FEATURE_TOGGLE_DATABASE_HOST
  FEATURE_TOGGLE_DATABASE_PASSWORD
  FEATURE_TOGGLE_DATABASE_NAME
  FEATURE_TOGGLE_PORT
)

# Parse .env file and export variables (ignoring comments and empty lines)
while IFS='=' read -r key value; do
    # Trim whitespace from key and value
    key=$(echo "$key" | xargs)
    value=$(echo "$value" | xargs)
    # Skip empty lines or comments
    if [[ -z "$key" || "$key" =~ ^# ]]; then
        continue
    fi
    export "$key"="$value"
done < "$ENV_FILE"

# Check that all required variables are set
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: $var is not set in $ENV_FILE"
        exit 1
    fi
done

# Change to the repository root (where docker-compose.yml is expected)
cd "$(dirname "$0")"

# Tear down any existing compose stack (if running)
docker compose down || true

# Start the stack in detached mode
docker compose up -d
