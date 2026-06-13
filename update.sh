#!/bin/bash
set -e

IMAGE="ghcr.io/mosmo1212312121/huawei-solar-to-influx"
CONTAINER="huawei-solar-to-influx"
ENV_FILE="/home/mosmo/environment/huawei-to-influx.env"
NETWORK="monitoring-network"

echo "Pulling latest image..."
docker pull "$IMAGE:latest"

NEW_ID=$(docker inspect --format='{{.Id}}' "$IMAGE:latest")

OLD_ID=""
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER}$"; then
    OLD_ID=$(docker inspect --format='{{.Image}}' "$CONTAINER" 2>/dev/null || true)
    echo "Stopping and removing container: $CONTAINER"
    docker stop "$CONTAINER"
    docker rm "$CONTAINER"
fi

echo "Starting new container..."
docker run -d \
    --name "$CONTAINER" \
    --env-file "$ENV_FILE" \
    --network "$NETWORK" \
    "$IMAGE:latest"

echo "Container started: $(docker ps --filter name=$CONTAINER --format '{{.ID}}')"

if [ -n "$OLD_ID" ] && [ "$OLD_ID" != "$NEW_ID" ]; then
    echo "Removing old image: $OLD_ID"
    docker rmi "$OLD_ID" 2>/dev/null || echo "Old image already removed or in use, skipping."
fi

echo "Done. Running containers:"
docker ps --filter name="$CONTAINER"
