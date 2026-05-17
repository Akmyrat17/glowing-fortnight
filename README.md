# Go Echo Boilerplate

A production-ready Go web API boilerplate using Echo v4, PostgreSQL (pgx), Redis, Squirrel, Viper and Docker.

This repository follows the project's canonical architecture and conventions outlined in [RULES.MD](RULES.MD).

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose (recommended for development and CI)
- Make (optional but recommended)

## Quick Start (Docker Compose)

1. Clone the repository and change into it:

```bash
git clone <repository-url> portfolio
cd portfolio
```

2. Copy the example environment file and edit values:

```bash
cp .env.example .env
# Edit .env to set DB password, JWT secret, and other secrets
```

3. Start the stack (builds images, starts DB + Redis, runs the app):

```bash
make up
```

4. Seed the database (optional):

```bash
make seed
```

5. Check the API health endpoint:

```bash
curl http://localhost:8083/health
```

Default app port and other runtime settings come from `configs/config.yaml` and are overridden by `.env` and environment variables.

## Development (local without Docker)

1. Install Go modules:

```bash
go mod download
```

2. Ensure PostgreSQL and Redis are running locally and update `configs/config.yaml` accordingly.

3. Create the application database (name and user per `configs/config.yaml`).

4. Apply migrations (via Docker container or DB client). Using the container:

```bash
docker compose exec app ./server migrate up
```

Or run SQL files from `platform/database/migrations/` in order.

5. Seed the database:

```bash
make seed
```

6. Run the server locally for development:

```bash
go run ./cmd/api/main.go
```

## Make Targets

- `make up` — Start services with Docker Compose
- `make down` — Stop services
- `make logs` — Tail app logs
- `make seed` — Run database seeds
- `make build` — Build Docker image
- `make test` — Run unit tests
- `make migrate-up` — Apply pending migrations (containerized)
- `make migrate-down` — Roll back last migration (containerized)

## API Endpoints

- `GET /health` — Health check
- See module routes under `internal/modules/*/infra/http` for available endpoints (users, permissions, auth)

## Configuration

Configuration is loaded from `configs/config.yaml` first. If a root `.env` file exists it is merged on top, and environment variables override both.

Sensitive configuration (passwords, secrets) should live only in `.env` or environment variables — `configs/.env` is gitignored.

## Database Migrations

This project uses SQL migration files in `platform/database/migrations/`. Use the containerized migrate commands (via `make` or `docker compose exec`) to apply or roll back migrations.

## Project Structure

Follow the conventions in `RULES.MD`. Key folders:

- `cmd/` — entry points (`api`, `seed`)
- `internal/` — application code grouped by module
- `pkg/` — reusable packages (logger, req_ctx)
- `platform/` — infra: database, cache, migrations

## Finalization Checklist

Before marking this repository ready for release, ensure:

- `.env` has been created from `.env.example` locally and contains no secrets in source control
- Database migrations have been applied to production and staging databases
- Seeds have been run where appropriate
- Tests pass: `make test` (or `go test ./...`)
- CI is configured to run `make test` and build the Docker image

## Contributing

Follow `RULES.MD` for coding patterns and architecture. Use branch names per the Git Workflow section in that file.

---

If you'd like, I can also:

- Run `go test ./...` and fix any failing tests
- Add CI configuration (GitHub Actions) to run tests and build images
- Prepare a release checklist and changelog

Let me know which of those you'd like next.
<parameter name="filePath">c:\Users\user\Documents\boilerplate\README.md