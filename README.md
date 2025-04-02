# Transaction Processing Service

A Go-based service that handles financial transactions and user balances with PostgreSQL backend.

## Features

- Process win/lose transactions
- Track user balances
- Support multiple transaction sources (game, server, payment)
- Concurrent transaction processing
- Load tested for performance

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)
- PostgreSQL client (for running migrations)

## Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd awesomeProject
```

2. Choose how to start the services:

   a. Full setup with migrations (recommended):
   ```bash
   make up
   ```

   b. Start services only (no migrations):
   ```bash
   make up-only
   ```

   c. Run migrations separately if needed:
   ```bash
   make migrate
   ```

The service will be available at `http://localhost:8080`

To stop all services:
```bash
make down
```

## Development Setup

For local development without Docker:

1. Start only the PostgreSQL database:
```bash
docker-compose up -d db
```

2. Initialize the database:
```bash
make init-db
```

3. Run the application locally:
```bash
go run main.go
```

## API Endpoints

### Get User Balance
```
GET /user/{userId}
```

Response:
```json
{
    "userId": 1,
    "balance": "100.00"
}
```

### Create Transaction
```
POST /user/{userId}/transaction
Header: Source-Type: game|server|payment
```

Request body:
```json
{
    "state": "win|lose",
    "amount": "10.00",
    "transactionId": "unique-transaction-id"
}
```

## Running Tests

Run the load tests:
```bash
go test ./test -v
```

The load test simulates 100 concurrent transactions at 25 TPS.

## Database Schema

The service uses two main tables:
- `users`: Stores user IDs and their current balance
- `transactions`: Records all processed transactions

## Makefile Commands

- `make build`: Build the Docker images only
- `make up`: Complete setup - build, start services, and run migrations
- `make up-only`: Start services without migrations (includes db health check)
- `make down`: Stop and remove all services
- `make restart`: Restart all services (includes migrations)
- `make migrate`: Run database migrations manually
- `make init-db`: Full initialization (same as `make up`)
- `make check-db`: Check database connectivity (used internally)

## Docker Commands

- `docker-compose up -d`: Start all services
- `docker-compose down`: Stop all services
- `docker-compose logs -f app`: Follow application logs
- `docker-compose logs -f db`: Follow database logs
- `docker-compose restart app`: Restart the application

## Configuration

The service uses the following default configuration:
- Database: PostgreSQL on localhost:5432
- Application port: 8080
- Database credentials: user/password
- Database name: core
