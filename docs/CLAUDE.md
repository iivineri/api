# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Iivineri API is a Go-based REST API built with the Fiber framework. It provides authentication services with JWT and 2FA support, using PostgreSQL for data storage and Redis for caching. The project includes Swagger documentation, Prometheus metrics, and Docker containerization.

## Architecture

The project follows a modular architecture with clear separation of concerns:

- **`main.go`**: Entry point using Cobra CLI framework
- **`cmd/`**: CLI commands (`serve` for running the API, `migration` for database operations)
- **`internal/`**: Application code organized by domain:
  - `config/`: Configuration management (app, database, cache, thumbor)
  - `fiber/`: HTTP server with modules (auth module with MVC pattern)
  - `database/`: Database connection and utilities
  - `logger/`: Structured logging with logrus
  - `metrics/`: Prometheus metrics collection
  - `migration/`: Database migration utilities
  - `wire/`: Dependency injection using Google Wire
- **`pkg/swagger/`**: Generated Swagger documentation
- **`migrations/`**: SQL migration files

### Key Components

- **Dependency Injection**: Uses Google Wire for compile-time DI
- **Authentication Module**: Located in `internal/fiber/modules/auth/` with full MVC structure
- **Middleware Stack**: Request ID, metrics, logging, recovery, security, CORS, compression, rate limiting
- **Health Checks**: Available at `/health` endpoint
- **Metrics**: Prometheus metrics at `/metrics` endpoint

## Development Commands

### Building and Running
```bash
# Build the application
go build -o iivineri .

# Run locally (development mode)
go run main.go serve

# Generate Swagger documentation
swag init -g main.go -o pkg/swagger
```

### Docker Development
```bash
# Setup Docker environment (checks ports, creates .env)
./tool setup

# Start basic services (API + PostgreSQL + Redis)
./tool start

# Start with development tools (+ Adminer + Redis Commander)
./tool dev

# Start with monitoring (+ Prometheus)
./tool monitoring

# View logs
./tool logs api -f

# Stop all services
./tool stop

# Complete cleanup (containers, images, volumes)
./tool cleanup
```

### Database Operations
```bash
# Run migrations
go run main.go migration up

# Rollback migrations
go run main.go migration down

# Create new migration
go run main.go migration create <name>

# Docker database operations
./tool backup-db
./tool restore-db backup.sql
./tool shell postgres
```

### Testing
```bash
# Run tests (no test files currently exist)
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./internal/fiber/modules/auth/...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run vet for static analysis
go vet ./...

# Run with race detection
go run -race main.go serve
```

## Environment Configuration

Key environment variables (see `.env.example`):

- `ENV`: Environment (development/production)
- `PORT`: API server port (default: 8080)
- `LOG_LEVEL`: Logging level (debug/info/warn/error)
- Database: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USERNAME`, `DB_PASSWORD`
- Cache: `CACHE_HOST`, `CACHE_PORT` (Redis)
- Docker ports: `API_EXTERNAL_PORT`, `DB_EXTERNAL_PORT`, `CACHE_EXTERNAL_PORT`

## Service Access

When running with `./tool dev`:
- **API**: http://localhost:8080
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health
- **Metrics**: http://localhost:8080/metrics
- **Adminer (DB UI)**: http://localhost:8081
- **Redis Commander**: http://localhost:8082
- **Prometheus**: http://localhost:9090 (with monitoring profile)

## Authentication

The API uses JWT authentication with 2FA support. Key endpoints:
- `POST /api/v1/auth/register`: User registration
- `POST /api/v1/auth/login`: User login
- Authentication header: `Authorization: Bearer <token>`

## Development Notes

- The project uses Go modules for dependency management
- Wire is used for dependency injection (run `wire` to regenerate)
- Swagger documentation is generated from code annotations
- The Docker tool automatically handles port conflicts
- Database migrations are managed through the CLI
- Metrics are automatically collected for all HTTP requests
- CORS is configured for all origins in development mode