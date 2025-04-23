# Swift Codes API

A REST API service built with Go and MongoDB.

## Prerequisites

- Docker and Docker Compose
- Make sure ports 8080 (or your configured PORT) are available on your machine

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
env PORT=8080 DB_USERNAME=your_username DB_PASSWORD=your_password
``` 

## Running the Application

To start the application and MongoDB:

```bash
  docker compose up
```

To run in detached mode:

```bash
  docker compose up -d
```

The project uses Docker Compose watch mode for development. Any changes to `go.mod` will trigger automatic rebuilds.

## Services

The application consists of two main services:

- API service running on port 8080 (or your configured PORT)
- MongoDB instance with persistent storage

## API Endpoints

The API server runs on `http://localhost:8080`

## Testing

### Unit Tests

Unit tests can be run locally without any external dependencies:

```bash 
  go test ./tests/unit/...
``` 

These tests use mocks to simulate the repository layer and don't require a MongoDB connection.

### Integration Tests

Integration tests are configured to run with a dedicated test service (`api-test`) and MongoDB instance in Docker
Compose (add `--build` flag to apply code changes):

```bash 
  docker compose up api-test
``` 

This will start the test service that connects to a separate MongoDB instance and runs all integration tests.

## Stopping the Application

To stop the services:

```bash 
  docker compose down
``` 

To stop and remove volumes:

```bash 
  docker compose down -v
```