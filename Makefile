.PHONY: up down logs seed migrate-up migrate-down build test help \
        local-up local-down local-migrate-up local-migrate-down local-test

# ─── Config ────────────────────────────────────────────────────────────────────
include .env
export

DB_URL=postgresql://$(DATABASE_USER):$(DATABASE_PASSWORD)@localhost:5432/$(DATABASE_NAME)?sslmode=disable
DB_URL_DOCKER=postgresql://$(DATABASE_USER):$(DATABASE_PASSWORD)@db:5432/$(DATABASE_NAME)?sslmode=disable

# ─── Help ──────────────────────────────────────────────────────────────────────
help:
	@echo ""
	@echo "DOCKER targets:"
	@echo "  make up                - Start all services (with build)"
	@echo "  make down              - Stop all services"
	@echo "  make logs              - Tail app logs"
	@echo "  make seed              - Run seed inside container"
	@echo "  make build             - Build Docker image"
	@echo "  make migrate-up        - Run pending migrations (in container)"
	@echo "  make migrate-down      - Rollback last migration (in container)"
	@echo "  make migrate-create    - Create new migration (in container)"
	@echo "  make clean             - Stop and remove volumes"
	@echo "  make ps                - Show running containers"
	@echo ""
	@echo "LOCAL targets:"
	@echo "  make local-up          - Run app locally (go run)"
	@echo "  make local-test        - Run tests locally"
	@echo "  make local-migrate-up  - Run pending migrations locally"
	@echo "  make local-migrate-down- Rollback last migration locally"
	@echo "  make local-migrate-create - Create new migration locally"
	@echo ""

# ─── Docker ────────────────────────────────────────────────────────────────────
up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f app

seed:
	docker compose exec app ./seed

build:
	docker compose build

migrate-up:
	docker compose exec app migrate -path /app/migrations -database "$(DB_URL_DOCKER)" up

migrate-down:
	docker compose exec app migrate -path /app/migrations -database "$(DB_URL_DOCKER)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	docker compose exec app migrate create -ext sql -dir /app/migrations $$name

clean:
	docker compose down -v

ps:
	docker compose ps

# ─── Local ─────────────────────────────────────────────────────────────────────
local-up:
	go run .\cmd\api\main.go

local-test:
	go test ./...

local-migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

local-migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down

local-migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations $$name