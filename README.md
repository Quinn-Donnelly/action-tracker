# Action Tracker

A habit tracking system built with Go and TimescaleDB for tracking actions over time.

## Features

- Track daily actions and habits
- RESTful API for creating and retrieving actions
- TimescaleDB for efficient time-series data storage
- Docker-compose for easy local development

## Quick Start

1. **Start the services:**
   ```bash
   docker-compose up -d
   ```

2. **Verify services are running:**
   ```bash
   docker-compose ps
   ```

3. **Test the API:**
   
   Create an action:
   ```bash
   curl -X POST http://localhost:8080/api/actions \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Morning Run",
       "description": "5km morning run",
       "category": "Exercise",
       "quantity": 5,
       "unit": "km"
     }'
   ```

   Get all actions:
   ```bash
   curl http://localhost:8080/api/actions
   ```

   Get actions by category:
   ```bash
   curl "http://localhost:8080/api/actions?category=Exercise"
   ```

## API Endpoints

### Actions

- `POST /api/actions` - Create a new action
- `GET /api/actions` - Get all actions (with optional filtering)
- `GET /api/actions/:id` - Get a specific action

### Query Parameters for GET /api/actions

- `name` - Filter by action name
- `category` - Filter by category
- `start_date` - Filter actions from this date (RFC3339 format)
- `end_date` - Filter actions until this date (RFC3339 format)
- `limit` - Maximum number of results (default: 100)
- `offset` - Number of results to skip (default: 0)

## Development

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)

### Local Development

1. **Start TimescaleDB:**
   ```bash
   docker-compose up -d timescaledb
   ```

2. **Run the Go application:**
   ```bash
   go run cmd/main.go
   ```

### Database Schema

The database is automatically initialized when TimescaleDB starts. See `scripts/init-db.sql` for the schema definition.

### Project Structure

```
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── db/                 # Database connection and migrations
│   ├── handlers/           # HTTP handlers
│   └── models/             # Data models
├── scripts/
│   └── init-db.sql         # Database initialization script
├── docker-compose.yml      # Docker services configuration
└── Dockerfile             # Application container definition
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string

Default values are configured for the docker-compose setup.

## Stopping the Services

```bash
docker-compose down
```

To also remove the database volume:
```bash
docker-compose down -v
```