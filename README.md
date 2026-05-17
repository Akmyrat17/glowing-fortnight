# Go Echo Boilerplate

A Go web API boilerplate using Echo framework, PostgreSQL, and Redis.

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose (for containerized setup)
- PostgreSQL and Redis (for local development)

## Quick Start with Docker Compose

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd boilerplate
   ```

2. Create a `.env` file:
   ```bash
   cp .env.example .env
   ```
   Update with your database password and other secrets.

3. Start the services:
   ```bash
   make up
   ```

   This will:
   - Build the Docker image
   - Start PostgreSQL, Redis, and the App
   - Automatically run all database migrations
   - Start the server

4. Seed the database (optional):
   ```bash
   make seed
   ```

5. The API will be available at `http://localhost:8083`

6. Check health endpoint:
   ```bash
   curl http://localhost:8083/health
   ```

## Manual Docker Setup (without Docker Compose)

If you prefer to run without Docker Compose:

1. Start PostgreSQL and Redis separately (e.g., using Docker or local installation)

2. Build the Docker image:
   ```bash
   docker build -t go-echo-boilerplate .
   ```

3. Create a root `.env` file based on `.env.example`:
   ```bash
   cp .env.example .env
   ```
   Update the database password, JWT secret, and any private keys.

4. Run the container:
   ```bash
   docker run -p 8080:8080 -v $(pwd)/configs:/app/configs --env-file .env go-echo-boilerplate
   ```

## Local Development Setup

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Set up PostgreSQL and Redis locally.

3. Create database `boilerplate` in PostgreSQL.

4. Update `configs/config.yaml` to use `localhost` for database and Redis hosts.

5. Run migrations manually:
   - Use your PostgreSQL client to run the `.up.sql` files in `migrations/` in order.

6. Seed the database:
   ```bash
   go run ./cmd/seed/main.go
   ```

7. Run the application:
   ```bash
   go run ./cmd/api/main.go
   ```

## Available Make Commands

- `make up` - Start all services with Docker Compose (automatically runs migrations)
- `make down` - Stop all services
- `make logs` - View application logs
- `make seed` - Run database seeding
- `make build` - Build Docker image
- `make test` - Run tests
- `make migrate-up` - Apply pending migrations manually
- `make migrate-down` - Rollback last migration manually

## API Endpoints

- `GET /health` - Health check
- User management endpoints (see routes in `internal/modules/user/infra/http/`)
- Permission management endpoints (see routes in `internal/modules/permission/infra/http/`)

## Configuration

Configuration is loaded from `configs/config.yaml` first. If a root `.env` file exists, it is merged on top of the YAML settings.

Environment variables from the operating system also override YAML and `.env` values.

For Docker, the compose stack can load `./.env` and pass those settings into the container.

For local development, use `configs/config.yaml` for general settings and `.env` for secrets or runtime overrides.

## Database Migrations

**Note**: Database migration commands are planned but not yet implemented. For now:

1. Run the `.up.sql` files in `migrations/` manually in order.
2. For rollbacks, run the corresponding `.down.sql` files.

## Project Structure

See `RULES.MD` for detailed architecture and coding conventions.

## What's Missing

- Database migration commands (`migrate up`, `migrate down`, `migrate create`) - currently TODO
- Environment variable loading from `.env` file (config currently loads from YAML)
- Comprehensive tests
- API documentation</content>
<parameter name="filePath">c:\Users\user\Documents\boilerplate\README.md