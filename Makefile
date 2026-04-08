# ─── Crowdfunding Platform Makefile ──────────────────────────────────────────

.PHONY: up down build restart logs \
        run tidy lint test \
        db-shell redis-shell \
        help

# ── Docker Compose ───────────────────────────────────────────────────────────

## Start all services in background
up:
	docker compose up -d

## Start all services and show logs
up-logs:
	docker compose up

## Stop and remove containers
down:
	docker compose down

## Rebuild images and start
build:
	docker compose up --build -d

## Restart all services
restart:
	docker compose restart

## Follow logs (all services)
logs:
	docker compose logs -f

## Follow backend logs only
logs-back:
	docker compose logs -f back

# ── Local development (without Docker) ───────────────────────────────────────

## Run backend locally (requires local Postgres and Redis)
run:
	cd back && go run ./cmd/main.go

## Download and tidy Go dependencies
tidy:
	cd back && go mod tidy

## Build backend binary
build-bin:
	cd back && go build -o bin/server ./cmd/main.go

## Run tests
test:
	cd back && go test ./...

## Run linter (requires golangci-lint)
lint:
	cd back && golangci-lint run ./...

# ── Database ─────────────────────────────────────────────────────────────────

## Open psql shell in running Postgres container
db-shell:
	docker compose exec postgres psql -U postgres -d crowdfunding

## Open redis-cli in running Redis container
redis-shell:
	docker compose exec redis redis-cli

# ── Help ─────────────────────────────────────────────────────────────────────

## Show this help
help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@grep -E '^##' Makefile | sed 's/## /  /'
	@echo ""
