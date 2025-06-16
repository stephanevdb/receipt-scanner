# Receipt Scanner

A receipt scanning application built with PocketBase and Docker.

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. Clone this repository
2. Run the application:
   ```bash
   docker-compose up -d
   ```
3. Access the PocketBase admin UI at: http://localhost:8080/_/

## Features

- Runs PocketBase in a Docker container
- Persists data using Docker volumes
- Automatically restarts on failure
- Secure non-root user configuration

## Stopping the Application

To stop the application:
```bash
docker-compose down
```

## Data Persistence

The application data is stored in the `pb_data` directory, which is mounted as a volume in the Docker container.