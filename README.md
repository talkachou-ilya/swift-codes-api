# Swift Codes API

A REST API service built with Go and MongoDB for managing and retrieving SWIFT codes.

## Prerequisites

- Docker and Docker Compose
- Make sure ports 8080 (API service) and 27017 (MongoDB) are available on your machine (or configure other ports for the application)

## Environment Variables

The project uses environment variables for configuration. Create a `.env` file in the root directory based on the provided `.env.example`:
```
PORT=8080 DB_USERNAME=root DB_PASSWORD=root DB_NAME=swiftdb DB_PORT=27017 TEST_PORT=8081 TEST_DB_PORT=27018
``` 

You can modify these values according to your preferences:
- `PORT`: The port on which the API service will run
- `DB_USERNAME` and `DB_PASSWORD`: MongoDB credentials
- `DB_NAME`: MongoDB database name
- `DB_PORT`: Port mapping for the MongoDB service
- `TEST_PORT`: Port for the test API service
- `TEST_DB_PORT`: Port mapping for the test MongoDB service

## Running the Application

To start the application and all required services:
```bash 
  docker compose up api
``` 

To run in detached mode:
```bash 
  docker compose up api -d
``` 

The project uses Docker Compose watch mode for development, which automatically rebuilds the application when changes are detected in the project files.

## Services

The application consists of the following services:

- **api**: Main API service running on port 8080 (or your configured PORT)
- **mongo**: MongoDB instance with persistent storage, accessible on port 27017 (or your configured DB_PORT)
- **mongo-seed**: Service that seeds the database with initial data
- **api-test**: Test service for running integration tests
- **mongo-test**: Dedicated MongoDB instance for testing, accessible on port 27018 (or your configured TEST_DB_PORT)

## API Endpoints

The API server runs on `http://localhost:8080` (by default) with the following endpoints:

- **GET /v1/swift-codes/:swift-code** - Retrieve a specific SWIFT code by its identifier
- **GET /v1/swift-codes/country/:countryISO2code** - Get all SWIFT codes for a specific country
- **POST /v1/swift-codes** - Add a new SWIFT code
- **DELETE /v1/swift-codes/:swift-code** - Delete a SWIFT code by its identifier

All endpoints return JSON responses.

## Testing

### Unit Tests

Unit tests can be run locally without any external dependencies:
```bash
  go test ./tests/unit/... -v
``` 

These tests use mocks to simulate the repository layer and don't require a MongoDB connection.

### Integration Tests

Integration tests are configured to run with a dedicated test service (`api-test`) and MongoDB instance:
```bash 
  docker compose up api-test
``` 

You can add `--build` flag to ensure that the latest code changes are applied before testing.

This will start the test service that connects to a separate MongoDB instance and runs all integration tests.

## Development

The application is structured using a clean architecture pattern:

- **cmd**: Application entry point
- **handlers**: HTTP request handlers
- **repositories**: Data access layer
- **models**: Data models
- **routes**: API route definitions

## Stopping the Application

To stop the services:
```bash 
  docker compose down
``` 

To stop and remove volumes (clears database data):
```bash 
  docker compose down -v
```
