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

The application consists of two services:
- API service running on port 8080 (or your configured PORT)
- MongoDB instance with persistent storage

## API Endpoints

The API server runs on `http://localhost:8080`

## Stopping the Application

To stop the services:
```bash 
docker compose down
``` 

To stop and remove volumes:
```bash 
docker compose down -v
```