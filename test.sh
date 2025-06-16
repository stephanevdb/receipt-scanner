#!/bin/bash

# Print commands and their arguments as they are executed
set -x

# Function to handle cleanup
cleanup() {
    echo -e "\nStopping the application..."
    docker-compose down
    exit 0
}

# Set up trap for SIGINT (Ctrl+C) and SIGTERM
trap cleanup SIGINT SIGTERM

# Stop any running containers
echo "Stopping existing containers..."
docker-compose down

# Remove old images to ensure a clean build
echo "Removing old images..."
docker-compose rm -f
docker-compose pull

# Rebuild the images
echo "Rebuilding images..."
docker-compose build --no-cache

# Start the application
echo "Starting the application..."
docker-compose up -d

# Show the logs and handle ESC key
echo "Container logs (Press ESC to exit):"
(
    # Start docker-compose logs in the background
    docker-compose logs -f &
    LOGS_PID=$!

    # Wait for ESC key
    while true; do
        # Read a single character
        read -rsn1 key
        # Check if it's ESC (ASCII 27)
        if [[ $key == $'\e' ]]; then
            kill $LOGS_PID
            cleanup
        fi
    done
) 